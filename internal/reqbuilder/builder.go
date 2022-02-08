package reqbuilder

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/bytebufferpool"
)

type RequestBuilder struct {
	method             string
	host               string
	path               string
	headerKeyAndValues []string
	queryKeyAndValues  []string
	body               interface{}
}

func New() *RequestBuilder {
	return &RequestBuilder{}
}

func (b *RequestBuilder) Method(method string) *RequestBuilder {
	b.method = method
	return b
}

func (b *RequestBuilder) Host(host string) *RequestBuilder {
	b.host = host
	return b
}

func (b *RequestBuilder) Path(path string) *RequestBuilder {
	b.path = path
	return b
}

func (b *RequestBuilder) Headers(keyAndValues ...string) *RequestBuilder {
	b.headerKeyAndValues = append(b.headerKeyAndValues, keyAndValues...)
	return b
}

func (b *RequestBuilder) Queries(keyAndValues ...string) *RequestBuilder {
	b.queryKeyAndValues = append(b.queryKeyAndValues, keyAndValues...)
	return b
}

func (b *RequestBuilder) Body(body interface{}) *RequestBuilder {
	b.body = body
	return b
}

func (b *RequestBuilder) Build() (*http.Request, error) {
	pathBuf := bytebufferpool.Get()
	pathBuf.B = append(pathBuf.B, b.host...)
	pathBuf.B = append(pathBuf.B, b.path...)
	var (
		n    = len(pathBuf.B)
		k, v string
	)
	for i := 1; i < len(b.queryKeyAndValues); i += 2 {
		k = b.queryKeyAndValues[i-1]
		v = b.queryKeyAndValues[i]
		if len(v) == 0 {
			continue
		}
		pathBuf.B = append(pathBuf.B, '&')
		pathBuf.B = append(pathBuf.B, k...)
		pathBuf.B = append(pathBuf.B, '=')
		pathBuf.B = append(pathBuf.B, v...)
	}
	if len(pathBuf.B) > n {
		pathBuf.B[n] = '?'
	}

	var bodyReader io.Reader
	if r, ok := b.body.(io.Reader); ok {
		bodyReader = r
	} else if b.body != nil {
		data, err := json.Marshal(b.body)
		if err != nil {
			return nil, errors.Wrap(err, "api: body encode error")
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(
		b.method,
		string(pathBuf.B),
		bodyReader,
	)
	if err != nil {
		return nil, errors.Wrap(err, "api: failed to create request object")
	}

	for i := 1; i < len(b.headerKeyAndValues); i += 2 {
		req.Header.Add(b.headerKeyAndValues[i-1], b.headerKeyAndValues[i])
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	return req, nil
}
