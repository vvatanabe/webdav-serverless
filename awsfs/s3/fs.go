package s3

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/webdav-serverless/webdav-serverless/awsfs"
)

type Server struct {
	Root          string
	PhysicalStore awsfs.PhysicalStore
	TempDir       string
}

type Pointer struct {
	ID   string `json:"id"`
	Size int64  `json:"size"`
}

// slashClean is equivalent to but slightly more efficient than
// path.Clean("/" + name).
func slashClean(name string) string {
	if name == "" || name[0] != '/' {
		name = "/" + name
	}
	return path.Clean(name)
}

func (s *Server) resolve(name string) string {
	// This implementation is based on Dir.Open's code in the standard net/http package.
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return ""
	}
	dir := s.Root
	if dir == "" {
		dir = "."
	}
	return filepath.Join(dir, filepath.FromSlash(slashClean(name)))
}

type FileInfo struct {
	os.FileInfo
	size int64
}

func (fi FileInfo) Size() int64 {
	return fi.size
}
