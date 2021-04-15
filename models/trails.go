package models

import (
	"bytes"
	"compress/gzip"
	"io"

	"gorm.io/gorm"
)

type Trails struct {
	gorm.Model
	Owner   *User
	OwnerID int
	Name    string
	Title   string
	GPX     []byte
}

func (t *Trails) SetGPX(gpx *bytes.Buffer) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(gpx.Bytes())
	writer.Flush()
	t.GPX = buf.Bytes()
}

func (t *Trails) GetGPX() *bytes.Buffer {
	var decompressed bytes.Buffer
	zr, _ := gzip.NewReader(bytes.NewBuffer(t.GPX))
	io.Copy(&decompressed, zr)
	zr.Close()
	return &decompressed
}
