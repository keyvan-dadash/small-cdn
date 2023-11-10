package cache

import (
	"bytes"
	"io"

	"github.com/tdewolff/minify/v2"
)

type CSSCacheFile struct {
	minifier *minify.M
	reader   io.Reader
	buffer   *bytes.Buffer
}

func CreateCSSCacheFile(minifier *minify.M) *CSSCacheFile {
	return &CSSCacheFile{
		minifier: minifier,
		buffer:   new(bytes.Buffer),
	}
}

func (c *CSSCacheFile) Write(writer io.Writer) error {
	_, err := writer.Write(c.buffer.Bytes())
	return err
}

func (c *CSSCacheFile) Read(reader io.Reader) error {
	c.reader = reader
	return nil
}

func (c *CSSCacheFile) Process() error {
	if err := c.minifier.Minify("text/css", c.buffer, c.reader); err != nil {
		return err
	}

	return nil
}

func (c *CSSCacheFile) GetCacheFileType() uint8 {
	return kCSSFile
}
