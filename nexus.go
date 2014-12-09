// Package nexus is a client for Go code to pull data from a Sonatype Nexus instance. Nexus provides a REST API,
// although some information may require several calls to collate all the data. So this client provides some methods
// to abstract away the necessary plumbing.
package nexus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/search"
)

// Client accesses a Nexus instance. The default Client should work for the newest Nexus versions. Older Nexus
// versions may need or benefit from a specific client.
type Client interface {
	// Returns all artifacts in this Nexus which satisfy the given criteria.
	Artifacts(criteria search.Criteria) ([]*Artifact, error)

	// Returns all repositories in this Nexus.
	Repositories() ([]*Repository, error)

	// Returns a Client which uses the given credentials to access this Nexus, without modifying the original client.
	// nil is understood to mean credentials.None.
	As(credentials credentials.Credentials) Client
}

// Nexus2x represents a Nexus v2.x instance. It's the default Client implementation.
type Nexus2x struct {
	Url string // e.g. http://nexus.somewhere.com:8080/nexus

	credentials.Credentials // e.g. BasicAuth{"username", "password"}
}

// New creates a new Nexus client, using the default Client implementation.
func New(url string, c credentials.Credentials) Client {
	return &Nexus2x{Url: url, Credentials: credentials.OrZero(c)}
}

// builds the proper URL with parameters for GET-ing.
func (nexus Nexus2x) fullUrlFor(query string, filter map[string]string) string {
	params := []string{}

	for k, v := range filter {
		params = append(params, k+"="+v)
	}

	if len(params) == 0 {
		return nexus.Url + "/" + query
	} else {
		return nexus.Url + "/" + query + "?" + strings.Join(params, "&")
	}
}

// does the actual legwork, going to Nexus and validating the response.
func (nexus Nexus2x) fetch(url string, params map[string]string) (*http.Response, error) {
	get, err := http.NewRequest("GET", nexus.fullUrlFor(url, params), nil)
	if err != nil {
		return nil, err
	}

	nexus.Credentials.Sign(get)
	get.Header.Add("Accept", "application/json")

	// go for it!
	response, err := http.DefaultClient.Do(get)
	if err != nil {
		return nil, err
	}

	// for us, 4xx are 5xx are errors, so we need to validate the response
	if 400 <= response.StatusCode && response.StatusCode < 600 {
		return response, &BadResponseError{Url: nexus.Url, StatusCode: response.StatusCode, Status: response.Status}
	}

	// everything alright, carry on
	return response, err
}

func bodyToBytes(body io.ReadCloser) ([]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, err
	}
	defer body.Close() // don't forget to Close() body at the end!

	return buf.Bytes(), nil
}

// As returns a copy of this client which uses the given credentials to access the same Nexus instance. Nil is
// understood to mean credentials.None.
func (nexus Nexus2x) As(c credentials.Credentials) Client {
	return &Nexus2x{Url: nexus.Url, Credentials: credentials.OrZero(c)}
}

// Artifacts returns all artifacts in this Nexus which satisfy the given criteria. This implementation errors out on a
// full search (n.Artifacts(search.None)).
func (nexus Nexus2x) Artifacts(criteria search.Criteria) ([]*Artifact, error) {
	params := criteria.Parameters()

	if len(params) == 0 {
		return nil, fmt.Errorf("Full search isn't supported!")
	}

	if len(params) == 1 {
		if repoId, ok := params["repositoryId"]; ok {
			return nexus.readArtifactsFrom(repoId)
		}
	}

	return nexus.readArtifactsWhere(params)
}

type artifactSearchResponse struct {
	TotalCount int
	Data       []struct {
		GroupId      string
		ArtifactId   string
		Version      string
		ArtifactHits []struct {
			ArtifactLinks []struct {
				Extension  string
				Classifier string
			}
		}
	}
}

func extractArtifactPayloadFrom(body []byte) (*artifactSearchResponse, error) {
	var payload *artifactSearchResponse

	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func extractArtifactsFrom(payload *artifactSearchResponse) []*Artifact {
	var artifacts = []*Artifact{}

	for _, artifact := range payload.Data {
		g := artifact.GroupId
		a := artifact.ArtifactId
		v := artifact.Version

		for _, hit := range artifact.ArtifactHits {
			for _, link := range hit.ArtifactLinks {
				e := link.Extension
				c := link.Classifier

				artifacts = append(artifacts, &Artifact{g, a, v, c, e})
			}
		}
	}

	return artifacts
}

// returns all artifacts in this Nexus which pass the given filter. The expected keys in filter are the flags Nexus'
// REST API accepts, with the same semantics.
func (nexus Nexus2x) readArtifactsWhere(filter map[string]string) ([]*Artifact, error) {
	// This implementation is slightly tricky. As artifactSearchResponse shows, Nexus always wraps the artifacts in a
	// GAV structure. This structure doesn't mean that within the wrapper are *all* the artifacts within that GAV, or
	// that the next page won't repeat artifacts if an incomplete GAV was returned earlier.
	//
	// On top of that, I haven't quite figured out how Nexus is counting artifacts for paging purposes. POMs don't
	// seem to count as artifacts, except when the project has a 'pom' packaging (which I can't know for sure without
	// GET-ing every POM), but the math still doesn't quite come together. So I took a conservative approach, which
	// forces a sequential algorithm. This search can be parallelized if the paging problem is solved.

	from := 0
	offset := 0
	started := false              // do-while can sometimes be useful
	artifacts := newArtifactSet() // acumulates the artifacts

	for offset != 0 || !started {
		started = true // do-while can sometimes be useful

		from = from + offset
		filter["from"] = strconv.Itoa(from)

		resp, err := nexus.fetch("service/local/lucene/search", filter)
		if err != nil {
			return nil, err
		}

		body, err := bodyToBytes(resp.Body)
		if err != nil {
			return nil, err
		}

		payload, err := extractArtifactPayloadFrom(body)
		if err != nil {
			return nil, err
		}

		// extract and store the artifacts. The set ensures we ignore repeated artifacts.
		artifacts.add(extractArtifactsFrom(payload))

		// a lower bound for the number of artifacts returned, since every GAV holds at least one artifact.
		// There will be some repetitions, but artifacts takes care of that.
		offset = len(payload.Data)
	}

	return artifacts.data, nil
}

// returns the first-level directories in the given repository.
func (nexus Nexus2x) firstLevelDirsOf(repositoryId string) ([]string, error) {
	// XXX Don't forget the ending /, or the response is always XML!
	resp, err := nexus.fetch("service/local/repositories/"+repositoryId+"/content/", nil)
	if err != nil {
		return nil, err
	}

	// fill payload with the given response
	body, err := bodyToBytes(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload *struct {
		Data []struct {
			Leaf bool
			Text string
		}
	}

	err = json.Unmarshal([]byte(body), &payload)
	if err != nil {
		return nil, err
	}

	// extract the directories from payload
	result := []string{}
	for _, dir := range payload.Data {
		if !dir.Leaf {
			result = append(result, dir.Text)
		}
	}

	return result, nil

}

// returns all artifacts in the given repository.
func (nexus Nexus2x) readArtifactsFrom(repositoryId string) ([]*Artifact, error) {
	// This function also has some tricky details. In the olden days (around version 1.8 or so), one could get all the
	// artifacts in a given repository by searching for *. This has been disabled in the newer versions, without any
	// official alternative for "give me everything you have". So, the solution adopted here is:
	// 1) get the first level directories in repositoryId
	// 2) for every directory 'dir', do a search filtering for the groupId 'dir*' and the repository ID
	// 3) accumulate the results in an artifactSet to avoid duplicates (e.g. the results in common* appear also in com*)

	result := newArtifactSet()

	// 1)
	dirs, err := nexus.firstLevelDirsOf(repositoryId)
	if err != nil {
		return nil, err
	}

	// 2) these searches can be done concurrently :)
	artifacts := make(chan []*Artifact)
	errors := make(chan error)
	for _, dir := range dirs {
		go func(dir string) {
			a, err := nexus.readArtifactsWhere(map[string]string{"g": dir + "*", "repositoryId": repositoryId})
			if err != nil {
				errors <- err
				return
			}

			artifacts <- a
		}(dir)
	}

	// 3)
	for i := 0; i < len(dirs); i++ {
		select {
		case a := <-artifacts:
			result.add(a)
		case err := <-errors:
			return nil, err
		}
	}

	return result.data, nil
}

// Repositories

// Repositories returns all repositories in this Nexus.
func (nexus Nexus2x) Repositories() ([]*Repository, error) {
	resp, err := nexus.fetch("service/local/repositories", nil)
	if err != nil {
		return nil, err
	}

	body, err := bodyToBytes(resp.Body)
	if err != nil {
		return nil, err
	}

	payload, err := extractRepoPayloadFrom(body)
	if err != nil {
		return nil, err
	}

	return extractReposFrom(payload), nil
}

type repoSearchResponse struct {
	Data []struct {
		Id         string
		Name       string
		RepoType   string
		RepoPolicy string
		Format     string
		RemoteUri  string
	}
}

func extractRepoPayloadFrom(body []byte) (*repoSearchResponse, error) {
	var payload *repoSearchResponse

	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func extractReposFrom(payload *repoSearchResponse) []*Repository {
	result := []*Repository{}

	for _, repo := range payload.Data {
		newRepo := &Repository{
			Id:        repo.Id,
			Name:      repo.Name,
			Type:      repo.RepoType,
			Format:    repo.Format,
			Policy:    repo.RepoPolicy,
			RemoteURI: repo.RemoteUri,
		}

		result = append(result, newRepo)
	}

	return result
}
