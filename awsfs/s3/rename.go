package s3

import (
	"context"
	"os"
	"path/filepath"
)

func (s *Server) Rename(ctx context.Context, oldName, newName string) error {
	if oldName = s.resolve(oldName); oldName == "" {
		return os.ErrNotExist
	}
	if newName = s.resolve(newName); newName == "" {
		return os.ErrNotExist
	}
	if root := filepath.Clean(s.Root); root == oldName || root == newName {
		// Prohibit renaming from or to the virtual root directory.
		return os.ErrInvalid
	}
	return os.Rename(oldName, newName)
}
