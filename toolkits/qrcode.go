// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package toolkits

import (
	"bytes"
	"fmt"
	"github.com/qjw/kelly"
	"github.com/skip2/go-qrcode"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"strconv"
)

const (
	QrcodeLow     = int(qrcode.Low)
	QrcodeMedium  = int(qrcode.Medium)
	QrcodeHigh    = int(qrcode.High)
	QrcodeHighest = int(qrcode.Highest)
)

type Qrcode struct {
	*qrcode.QRCode
}

func (q *Qrcode) Image(size int) image.Image {
	return q.QRCode.Image(size)
}

func (q *Qrcode) Write(size int, out io.Writer) error {
	return q.QRCode.Write(size, out)
}

func (q *Qrcode) WriteFile(size int, filename string) error {
	return q.QRCode.WriteFile(size, filename)
}

func (q *Qrcode) WriteKelly(size int, c *kelly.Context) error {
	img := q.Image(size)
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		c.WriteString(http.StatusInternalServerError, err.Error())
		return err
	}

	c.SetHeader("Content-Type", "image/png")
	c.SetHeader("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := c.Write(buffer.Bytes()); err != nil {
		c.WriteString(http.StatusInternalServerError, err.Error())
		return err
	}
	return nil
}

func NewQRCode(content string, level int) (*Qrcode, error) {
	if level < int(QrcodeLow) || level > int(QrcodeHighest) {
		panic(fmt.Errorf("invalid level %d", level))
	}

	var q *qrcode.QRCode
	q, err := qrcode.New(content, qrcode.RecoveryLevel(level))

	if err != nil {
		return nil, err
	}

	return &Qrcode{
		QRCode: q,
	}, nil
}
