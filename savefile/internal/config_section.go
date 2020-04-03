package internal

import (
	"io"
)

type ConfigSection interface {
	Name() string
	NameValue() string
	AppendLine(line string)
	Write(w io.Writer, newLine string) (n int64, err error)
}