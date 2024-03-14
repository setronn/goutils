package goutils

import "mime/multipart"

type FormString struct {
	name  string
	value string
}

func NewFormString(name string, value string) *FormString {
	return &FormString{name, value}
}

func (fs FormString) CreateForm(mw *multipart.Writer) {
	mw.WriteField(fs.name, fs.value)
}
