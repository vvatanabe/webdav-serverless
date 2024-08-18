package s3

import (
	"context"
	"os"
	"path/filepath"
)

func (s *Server) RemoveAll(ctx context.Context, name string) error {
	if name = s.resolve(name); name == "" {
		return os.ErrNotExist
	}
	if name == filepath.Clean(s.Root) {
		// Prohibit removing the virtual root directory.
		return os.ErrInvalid
	}
	return os.RemoveAll(name)
}
