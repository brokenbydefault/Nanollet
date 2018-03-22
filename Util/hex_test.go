package Util

import (
	"testing"
)


func TestUnsafeHexEncode(t *testing.T) {

	table := []struct{
		Message string
		Result string
	}{
		{"Testing the hex", "54657374696E672074686520686578"},
		{"Testing the hex?", "54657374696E6720746865206865783F"},
	}

	for _, v := range table {

		if r := UnsafeHexEncode([]byte(v.Message)); r != v.Result {
			t.Errorf("UnsafeHexEncode failed given %s expecting %s", r, v.Result)
		}
	}

}

func TestSecureHexEncode(t *testing.T) {

	table := []struct{
		Message string
		Result string
	}{
		{"Testing the hex", "54657374696E672074686520686578"},
		{"Testing the hex?", "54657374696E6720746865206865783F"},
	}

	for _, v := range table {

		if r := SecureHexEncode([]byte(v.Message)); r != v.Result {
			t.Errorf("SecureHexEncode failed given %s expecting %s", r, v.Result)
		}
	}

}

func TestSecureHexDecode(t *testing.T) {

	table := []struct{
		Message string
		Result string
	}{
		{ "54657374696E672074686520686578", "Testing the hex",},
		{"54657374696E6720746865206865783F", "Testing the hex?"},
		{"6E3F", "n?"},
		{"6e3f", "n?"},
	}

	for _, v := range table {

		r, ok := SecureHexDecode(v.Message)

		if string(r) != v.Result {
			t.Error( r, v.Result , ok)
		}


		if !ok {
			t.Error(r, v.Result , ok)
		}
	}

}


