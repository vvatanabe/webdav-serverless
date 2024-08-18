package s3

import (
	"context"
	"os"
)

func (s *Server) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	if name = s.resolve(name); name == "" {
		return os.ErrNotExist
	}
	return os.Mkdir(name, perm)
}
