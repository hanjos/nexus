// Package nexus is a client for Go code to pull data from a Sonatype Nexus instance. Nexus provides a REST API,
// although some information may require several calls to collate all the data. So this client provides some methods
// to abstract away the necessary plumbing.
package nexus

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/errors"
	"github.com/hanjos/nexus/search"
	"github.com/hanjos/nexus/util"
	"time"
)

// Client accesses a Nexus instance. The default Client should work for the newest Nexus versions. Older Nexus
// versions may need or benefit from a specific client.
type Client interface {
	// Returns all artifacts in this Nexus which satisfy the given criteria.
	// Nil is the same as search.None. If no criteria are given
	// (e.g. search.None), it does a full search in all repositories.
	Artifacts(criteria search.Criteria) ([]*Artifact, error)

	// Returns all repositories in this Nexus.
	Repositories() ([]*Repository, error)

	// Returns extra information about the given artifact.
	InfoOf(artifact *Artifact) (*ArtifactInfo, error)
}

// Nexus2x represents a Nexus v2.x instance. It's the default Client implementation.
type Nexus2x struct {
	Url         string                  // e.g. http://nexus.somewhere.com:8080/nexus
	Credentials credentials.Credentials // e.g. credentials.BasicAuth{"username", "password"}
	HttpClient  *http.Client            // the network client
}

// New creates a new Nexus client, using the default Client implementation.
func New(url string, c credentials.Credentials) Client {
	return &Nexus2x{Url: url, Credentials: credentials.OrZero(c), HttpClient: &http.Client{}}
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
	fullUrl, err := util.CleanSlashes(nexus.fullUrlFor(url, params))
	if err != nil {
		return nil, err
	}

	get, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	nexus.Credentials.Sign(get)
	get.Header.Add("Accept", "application/json")

	// go for it!
	response, err := nexus.HttpClient.Do(get)
	if err != nil {
		return nil, err
	}

	// lets see if everything is alright
	status := response.StatusCode
	switch true {
	case status == http.StatusUnauthorized:
		// the credentials don't check out
		return nil, &errors.UnauthorizedError{fullUrl, nexus.Credentials}
	case 400 <= status && status < 600:
		// Nexus complained, so error out
		return nil, &errors.BadResponseError{nexus.Url, status, response.Status}
	}

	// all is good, carry on
	return response, nil
}

func bodyToBytes(body io.ReadCloser) ([]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, err
	}
	defer body.Close() // don't forget to Close() body at the end!

	return buf.Bytes(), nil
}

// Artifacts implements the Client interface, returning all artifacts in this Nexus which satisfy the given criteria.
// Nil is the same as search.None. If no criteria are given (e.g. search.None), it does a full search in all
// repositories.
//
// Generally you don't want that, especially if you have proxy repositories; Maven Central (which many people will
// proxy) has, at the time of this comment, over 800,000 artifacts (!), which in this implementation will be all loaded
// into memory (!!). But, if you insist...
func (nexus Nexus2x) Artifacts(criteria search.Criteria) ([]*Artifact, error) {
	params := search.OrZero(criteria).Parameters()

	if len(params) == 0 { // full search
		return nexus.readAllArtifacts()
	}

	if len(params) == 1 {
		if repoId, ok := params["repositoryId"]; ok { // all in repo search
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
			RepositoryId  string
			ArtifactLinks []struct {
				Extension  string
				Classifier string
			}
		}
	}
}

func extractArtifactsFrom(payload *artifactSearchResponse) []*Artifact {
	var artifacts = []*Artifact{}

	for _, artifact := range payload.Data {
		g := artifact.GroupId
		a := artifact.ArtifactId
		v := artifact.Version

		for _, hit := range artifact.ArtifactHits {
			r := hit.RepositoryId
			for _, link := range hit.ArtifactLinks {
				e := link.Extension
				c := link.Classifier

				artifacts = append(artifacts, &Artifact{g, a, v, c, e, r})
			}
		}
	}

	return artifacts
}

// a slight modification of Go's v, ok := m[key] idiom. has returns false for ok if value is "".
func has(m map[string]string, key string) (value string, ok bool) {
	value, ok = m[key]

	if ok && value == "" { // not ok!
		return value, false
	}

	return value, ok
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
	started := false
	artifacts := newArtifactSet() // accumulates the artifacts

	for offset != 0 || !started {
		started = true // do-while can sometimes be useful :)

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

		var payload *artifactSearchResponse
		err = json.Unmarshal(body, &payload)
		if err != nil {
			return nil, err
		}

		// extract the artifacts
		payloadArtifacts := extractArtifactsFrom(payload)

		// Nexus 2.x's search always returns the POMs, even when one filters specifically for the packaging or the
		// classifier. So we'll have to take them out here. Of course, if the user specifies "pom", she'll get POMs :)
		packaging, okPack := has(filter, "p") // using has instead of Go's idiom, since p="" still means no packaging
		_, okClass := has(filter, "c")

		if (okPack && packaging != "pom") || okClass { // remove the POMs
			for i := 0; i < len(payloadArtifacts); i++ {
				if payloadArtifacts[i].Extension == "pom" {
					payloadArtifacts = append(payloadArtifacts[:i], payloadArtifacts[i+1:]...)
				}
			}
		}

		// store the artifacts. The set ensures we ignore repetitions.
		artifacts.add(payloadArtifacts)

		// a lower bound for the number of artifacts returned, since every GAV in the payload holds at least one
		// artifact. There will be some repetitions, but artifacts takes care of that.
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

	//result := newArtifactSet()

	// 1)
	dirs, err := nexus.firstLevelDirsOf(repositoryId)
	if err != nil {
		return nil, err
	}

	return concurrentArtifactSearch(
		dirs,
		func(datum string) ([]*Artifact, error) {
			return nexus.readArtifactsWhere(map[string]string{"g": datum + "*", "repositoryId": repositoryId})
		})
}

// returns all artifacts visible by this Nexus.
func (nexus Nexus2x) readAllArtifacts() ([]*Artifact, error) {
	// there's no easy way to do this, so get the repos and search for all artifacts in each one (yup)
	repos, err := nexus.Repositories()
	if err != nil {
		return nil, err
	}

	// all we need for the search is the IDs
	ids := make([]string, len(repos))
	for i, repo := range repos {
		ids[i] = repo.Id
	}

	return concurrentArtifactSearch(
		ids,
		func(datum string) ([]*Artifact, error) { return nexus.readArtifactsFrom(datum) })
}

// InfoOf implements the Client interface, fetching extra information about the given artifact.
func (nexus Nexus2x) InfoOf(artifact *Artifact) (*ArtifactInfo, error) {
	// first resolve the artifact: building the URL by hand may fail in some situations (e.g. snapshot artifacts, odd
	// file names)
	path, err := nexus.repositoryPathOf(artifact)
	if err != nil {
		return nil, err
	}

	// now we can reliably build the proper URL
	resp, err := nexus.fetch("service/local/repositories/"+artifact.RepositoryId+"/content"+path,
		map[string]string{"describe": "info"})
	if err != nil {
		return nil, err
	}

	body, err := bodyToBytes(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload *infoSearchResponse
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	return extractInfoFrom(payload, artifact), nil
}

func (nexus Nexus2x) repositoryPathOf(artifact *Artifact) (string, error) {
	resp, err := nexus.fetch("service/local/artifact/maven/resolve",
		map[string]string{
			"g": artifact.GroupId,
			"a": artifact.ArtifactId,
			"v": artifact.Version,
			"e": artifact.Extension,
			"c": artifact.Classifier,
			"r": artifact.RepositoryId,
		})
	if err != nil {
		return "", err
	}

	body, err := bodyToBytes(resp.Body)
	if err != nil {
		return "", err
	}

	var payload *struct {
		Data struct {
			RepositoryPath string
		}
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		return "", err
	}

	return payload.Data.RepositoryPath, nil
}

type infoSearchResponse struct {
	Data struct {
		MimeType     string
		Uploader     string
		Uploaded     int64
		LastChanged  int64
		Size         int64
		Sha1Hash     string
		Repositories []struct {
			RepositoryId string
			ArtifactUrl  string
		}
	}
}

func extractInfoFrom(payload *infoSearchResponse, artifact *Artifact) *ArtifactInfo {
	url := ""
	for _, repo := range payload.Data.Repositories {
		if repo.RepositoryId == artifact.RepositoryId {
			url = repo.ArtifactUrl
			break
		}
	}

	return &ArtifactInfo{
		Artifact:    artifact,
		Uploader:    payload.Data.Uploader,
		Uploaded:    time.Unix(payload.Data.Uploaded, 0),
		LastChanged: time.Unix(payload.Data.LastChanged, 0),
		Sha1:        payload.Data.Sha1Hash,
		Size:        util.FileSize(payload.Data.Size),
		MimeType:    payload.Data.MimeType,
		Url:         url,
	}
}

// Repositories implements the Client interface, returning all repositories in this Nexus.
func (nexus Nexus2x) Repositories() ([]*Repository, error) {
	resp, err := nexus.fetch("service/local/repositories", nil)
	if err != nil {
		return nil, err
	}

	body, err := bodyToBytes(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload *repoSearchResponse
	err = json.Unmarshal(body, &payload)
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
