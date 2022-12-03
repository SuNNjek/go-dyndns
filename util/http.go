package util

import (
	"github.com/google/wire"
	"net/http"
)

var DefaultHttpClientValue = wire.InterfaceValue(new(HttpClient), http.DefaultClient)

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}
