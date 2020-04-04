package sections

import (
	"bytes"
	"io"
)

type RawConfigSection struct {
	name      string
	nameValue string // nameless
	raw       bytes.Buffer
}

func NewRawConfigSection(name string, nameValue string) *RawConfigSection {
	sect := RawConfigSection{
		name:      name,
		nameValue: nameValue,
		raw:       bytes.Buffer{},
	}

	return &sect
}

func (r *RawConfigSection) NameValue() string {
	return r.nameValue
}

func (r *RawConfigSection) Write(w io.Writer, newLine string) (n int64, err error) {
	return r.raw.WriteTo(w)
}

func (r *RawConfigSection) AppendLine(line string) {
	r.raw.WriteString(line)
}

func (r *RawConfigSection) Name() string {
	return r.name
}
