// +build !js

package DOM

import (
	"bytes"
	"os"
	"io"
)

func (el *Element) GetAttr(name string) (result string, err error) {
	return el.el.Attr(name)
}

func (dom *DOM) GetAttrOf(name string, css string) (result string, err error) {
	input, err := dom.SelectFirstElement(css)
	if err != nil {
		return
	}

	return input.GetAttr(name)
}

func (el *Element) GetText() (result string, err error) {
	return el.el.Text()
}

func (el *Element) GetStringValue() (result string, err error) {
	value, err := el.el.GetValue()
	if err != nil {
		return "", err
	}

	return value.String(), nil
}

func (dom *DOM) GetStringValueOf(css string) (result string, err error) {
	input, err := dom.SelectFirstElement(css)
	if err != nil {
		return
	}

	return input.GetStringValue()
}

func (el *Element) GetBytesValue() (result []byte, err error) {
	value, err := el.el.GetValue()
	if err != nil {
		return nil, err
	}

	return []byte(value.String()), nil
}

func (dom *DOM) GetBytesValueOf(css string) (result []byte, err error) {
	input, err := dom.SelectFirstElement(css)
	if err != nil {
		return
	}

	return input.GetBytesValue()
}

func (el *Element) GetFile() (io.Reader, error) {
	input, err := el.GetStringValue()
	if err != nil || input == "" {
		return nil, ErrInvalidElement
	}

	file, err := os.Open(input[7:])
	if err != nil {
		return nil, err
	}

	defer file.Close()

	if stat, err := file.Stat(); err != nil || stat.IsDir() {
		return nil, ErrInvalidElement
	}

	r := bytes.NewBuffer(nil)
	if _, err := io.Copy(r, file); err != nil {
		return nil, err
	}

	return r, nil
}

func (dom *DOM) GetFileOf(css string) (io.Reader, error) {
	input, err := dom.SelectFirstElement(css)
	if err != nil {
		return nil, err
	}

	return input.GetFile()
}
