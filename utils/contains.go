package utils

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func MapToObject(input interface{}, output interface{}) (interface{}, error) {
	c := mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           output,
		WeaklyTypedInput: true,
	}

	dec, err := mapstructure.NewDecoder(&c)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err = dec.Decode(input); err != nil {
		return nil, errors.WithStack(err)
	}

	return output, nil
}
