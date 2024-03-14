package goutils

import "mime/multipart"

type Form interface {
	CreateForm(*multipart.Writer)
}
