package goutils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
)

type FormFile struct {
	fieldname   string
	filename    string
	contentType string
	fileBytes   []byte
}

func NewFormFile(fieldname string, filename string, contentType string, fileBytes []byte) *FormFile {
	return &FormFile{fieldname, filename, contentType, fileBytes}
}

func (ff FormFile) CreateForm(mw *multipart.Writer) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, ff.fieldname, ff.filename))
	h.Set("Content-Type", ff.contentType)
	fw, _ := mw.CreatePart(h)

	fr := bytes.NewReader(ff.fileBytes)
	io.Copy(fw, fr)
}
