package rpctypes

import "io"

type Connection interface {
	SendRequest(b []byte) ([]byte, error)
	SendRequestReader(b []byte) (io.ReadCloser, error)
	SendRequestJSON(request interface{}, response interface{}, try ...interface{})error
}