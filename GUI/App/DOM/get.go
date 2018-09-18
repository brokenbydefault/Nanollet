// +build !js

package DOM

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

