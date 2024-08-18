package s3

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/google/uuid"
)

func (s *Server) Create(ctx context.Context, path string, flag int, perm os.FileMode, r io.Reader) (os.FileInfo, error) {
	if path = s.resolve(path); path == "" {
		return nil, os.ErrNotExist
	}

	isNotExist, err := s.fileExists(path)
	if err != nil {
		return nil, err
	}

	sr := &sizingReader{Reader: r}

	var p Pointer
	if isNotExist {
		p.ID = uuid.New().String()
	} else {
		p, err = s.readPointer(path, flag, perm)
		if err != nil {
			return nil, err
		}
	}

	if err := s.PhysicalStore.PutObjectLarge(ctx, p.ID, sr); err != nil {
		return nil, err
	}

	p.Size = sr.size
	if err := s.writePointer(path, flag, perm, p); err != nil {
		return nil, err
	}

	return s.getFileInfo(path, p.Size)
}

func (s *Server) fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return false, nil
	}
	if os.IsNotExist(err) {
		return true, nil
	}
	return false, err
}

func (s *Server) readPointer(path string, flag int, perm os.FileMode) (Pointer, error) {
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return Pointer{}, err
	}
	defer f.Close()

	return s.decodePointer(f)
}

func (s *Server) decodePointer(r io.Reader) (Pointer, error) {
	var p Pointer
	if err := json.NewDecoder(r).Decode(&p); err != nil {
		return Pointer{}, err
	}
	return p, nil
}

func (s *Server) writePointer(path string, flag int, perm os.FileMode, p Pointer) error {
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&p); err != nil {
		return err
	}
	return nil
}

func (s *Server) getFileInfo(path string, size int64) (os.FileInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		FileInfo: fi,
		size:     size,
	}, nil
}

type sizingReader struct {
	io.Reader
	size int64
}

func (r *sizingReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.size += int64(n)
	return
}
