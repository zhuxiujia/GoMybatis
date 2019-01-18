package utils

import (
	"bytes"
	"errors"
)

func MakeErrors(args ...string) error {
	var buffer bytes.Buffer
	for _, v := range args {
		buffer.WriteString(v)
	}
	return errors.New(buffer.String())
}
