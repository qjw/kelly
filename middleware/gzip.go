// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

// https://github.com/gin-contrib/gzip

package middleware

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"github.com/qjw/kelly"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

const (
	BestCompression = iota
	BestSpeed
	DefaultCompression
	NoCompression
	MaxCompressionLevel
)

const (
	GzipMethod = iota
	DeflateMethod
	MaxCompressionMethod
)

var (
	gzipLevels = [4]int{gzip.BestCompression, gzip.BestSpeed, gzip.DefaultCompression, gzip.NoCompression}
	zlibLevels = [4]int{zlib.BestCompression, zlib.BestSpeed, zlib.DefaultCompression, zlib.NoCompression}
)

func Gzip(level int, method int) kelly.HandlerFunc {

	var gzPool sync.Pool
	var methodStr = ""
	if method == GzipMethod {
		methodStr = gzipWriter{}.Name()
		gzPool.New = func() interface{} {
			return newGzipWriter(level)
		}
	} else if method == DeflateMethod {
		methodStr = deflateWriter{}.Name()
		gzPool.New = func() interface{} {
			return newDeflateWriter(level)
		}
	} else {
		panic(fmt.Errorf("invalid method %d", method))
	}

	return func(c *kelly.Context) {
		if !shouldCompress(c.Request(), methodStr) {
			c.InvokeNext()
			return
		}

		gz := gzPool.Get().(compressWriter)
		defer func() {
			gz.SetCompressionWriter(ioutil.Discard)
			gzPool.Put(gz)
		}()
		gz.SetCompressionWriter(c.ResponseWriter)
		gz.SetResponseWriter(c.ResponseWriter)

		c.SetHeader("Content-Encoding", methodStr)
		c.SetHeader("Vary", "Accept-Encoding")
		c.ResponseWriter = gz
		defer func() {
			gz.Close()
		}()
		c.InvokeNext()
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type compressWriter interface {
	http.ResponseWriter
	SetCompressionWriter(w io.Writer)
	SetResponseWriter(w http.ResponseWriter)
	Close()
	Name() string
}

// ---------------------------------------------------------------------------------------------------------------------

type gzipWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (g gzipWriter) Name() string {
	return "gzip"
}

func (g *gzipWriter) SetResponseWriter(w http.ResponseWriter) {
	g.ResponseWriter = w
}

func (g *gzipWriter) SetCompressionWriter(w io.Writer) {
	g.writer.Reset(w)
}

func (g *gzipWriter) Close() {
	g.writer.Close()
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func newGzipWriter(level int) compressWriter {
	if level < BestCompression || level >= MaxCompressionLevel {
		panic(fmt.Errorf("invalid level %d", level))
	}
	gz, err := gzip.NewWriterLevel(ioutil.Discard, gzipLevels[level])
	if err != nil {
		panic(err)
	}
	return &gzipWriter{
		writer: gz,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
type deflateWriter struct {
	http.ResponseWriter
	writer *zlib.Writer
}

func (g deflateWriter) Name() string {
	return "deflate"
}

func (g *deflateWriter) SetResponseWriter(w http.ResponseWriter) {
	g.ResponseWriter = w
}

func (g *deflateWriter) SetCompressionWriter(w io.Writer) {
	g.writer.Reset(w)
}

func (g *deflateWriter) Close() {
	g.writer.Close()
}

func (g *deflateWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func newDeflateWriter(level int) compressWriter {
	if level < BestCompression || level >= MaxCompressionLevel {
		panic(fmt.Errorf("invalid level %d", level))
	}

	gz, err := zlib.NewWriterLevel(ioutil.Discard, zlibLevels[level])
	if err != nil {
		panic(err)
	}
	return &deflateWriter{
		writer: gz,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func shouldCompress(req *http.Request, method string) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), method) {
		return false
	}
	extension := filepath.Ext(req.URL.Path)
	if len(extension) < 4 { // fast path
		return true
	}

	switch extension {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	default:
		return true
	}
}
