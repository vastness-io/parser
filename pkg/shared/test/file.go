package test

import (
	"gopkg.in/src-d/go-billy.v4"
	"io"
)

var _ billy.File = &MockFile{}

type MockFile struct {
	FileName string
	io.Reader
}

func (f *MockFile) Name() string {
	return f.FileName
}

func (*MockFile) Write(p []byte) (int, error) {
	return 0, nil
}

func (*MockFile) ReadAt(b []byte, off int64) (int, error) {
	return 0, nil
}

func (*MockFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (*MockFile) Close() error {
	return nil
}

func (*MockFile) Lock() error {
	return nil
}

func (*MockFile) Unlock() error {
	return nil
}

func (*MockFile) Truncate(size int64) error {
	return nil
}
