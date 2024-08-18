package s3

import (
	"context"
	"io"
	"io/fs"
	"os"

	"github.com/webdav-serverless/webdav-serverless/awsfs"
	"github.com/webdav-serverless/webdav-serverless/webdav"
)

func (s *Server) OpenFile(ctx context.Context, path string, flag int, perm os.FileMode) (webdav.File, error) {
	if path = s.resolve(path); path == "" {
		return nil, os.ErrNotExist
	}

	meta, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	fi, err := meta.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return &FileReader{
			meta: meta,
		}, nil
	}

	p, err := s.decodePointer(meta)
	if err != nil {
		return nil, err
	}

	r, err := s.PhysicalStore.GetObject(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	real, err := os.CreateTemp(s.TempDir, "webdav-temp-")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(real, r)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		real: real,
		meta: meta,
		size: p.Size,
	}, nil
}

type FileReader struct {
	real *os.File
	meta *os.File
	size int64
}

func (f FileReader) Close() error {
	if f.real == nil {
		return nil
	}
	realErr := f.real.Close()
	metaErr := f.meta.Close()
	_ = os.Remove(f.real.Name())
	if realErr != nil {
		return realErr
	}
	if metaErr != nil {
		return metaErr
	}
	return nil
}

func (f FileReader) Read(p []byte) (n int, err error) {
	return f.real.Read(p)
}

func (f FileReader) Seek(offset int64, whence int) (int64, error) {
	return f.real.Seek(offset, whence)
}

func (f FileReader) Readdir(count int) ([]fs.FileInfo, error) {
	return f.meta.Readdir(count)
}

func (f FileReader) Stat() (fs.FileInfo, error) {
	fi, err := f.meta.Stat()
	if err != nil {
		return nil, err
	}
	return FileInfo{
		FileInfo: fi,
		size:     f.size,
	}, nil
}

func (f FileReader) Write(p []byte) (n int, err error) {
	return 0, awsfs.ErrNotSupported
}
