// Package util store helper code that's useful but not... business logic, so to speak.
package util

import "fmt"

// FileSize represents an amount of bytes.
type FileSize int64

const (
	Byte     = FileSize(1)
	Kilobyte = FileSize(1 << 10)
	Megabyte = FileSize(1 << 20)
	Gigabyte = FileSize(1 << 30)
)

// String implements the Stringer interface, for easy printing.
func (size FileSize) String() string {
	switch true {
	case size <= Kilobyte:
		return fmt.Sprintf("%d B", int(size))
	case size <= Megabyte:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(Kilobyte))
	case size <= Gigabyte:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(Megabyte))
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(Gigabyte))
	}
}
