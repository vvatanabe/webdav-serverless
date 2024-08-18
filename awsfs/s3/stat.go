package s3

import (
	"context"
	"os"
)

func (s *Server) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if name = s.resolve(name); name == "" {
		return nil, os.ErrNotExist
	}

	meta, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	if meta.IsDir() {
		return meta, nil
	}

	pointer, err := s.readPointer(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	return FileInfo{
		FileInfo: meta,
		size:     pointer.Size,
	}, nil
}
