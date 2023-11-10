package cache

import (
	"bytes"
	"io"
)

type SimpleCacheFile struct {
	reader io.Reader
	buffer *bytes.Buffer
}

func CreateSimpleCacheFile() *SimpleCacheFile {
	return &SimpleCacheFile{
		buffer: new(bytes.Buffer),
	}
}

func (s *SimpleCacheFile) Write(writer io.Writer) error {
	_, err := writer.Write(s.buffer.Bytes())
	return err
}

func (s *SimpleCacheFile) Read(reader io.Reader) error {
	_, err := s.buffer.ReadFrom(reader)
	return err
}

func (s *SimpleCacheFile) Process() error {
	return nil
}

func (s *SimpleCacheFile) GetCacheFileType() uint8 {
	return kJSFile
}
