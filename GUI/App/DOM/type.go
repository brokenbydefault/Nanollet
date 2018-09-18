package DOM

import "errors"

type Application interface {
	Pages() []Page
	Name() string
	HaveSidebar() bool
}

type Page interface {
	OnContinue(w *Window, dom *DOM, action string)
	OnView(w *Window, dom *DOM)
	Name() string
}

type ActionMethod int

const (
	Click = iota
)

var (
	ErrInvalidActionMethod = errors.New("invalid action method")
)

type ReplaceMethod int

const (
	InnerReplaceContent ReplaceMethod = iota
	InnerPrepend
	InnerAppend

	OuterReplace
	OuterPrepend
	OuterAppend
)

var (
	ErrInvalidReplaceMethod = errors.New("invalid replace method")
)
