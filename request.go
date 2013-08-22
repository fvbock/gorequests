package gorequests

import (
	"bytes"
	"net/http"
)

// TYPE: http.Request
//      Name            | Type                 | Index   | Offset
// FIELDS               |                      |         |
//  F=> Method          | string               | [00000] | 0000000000
//  F=> URL             | *url.URL             | [00001] | 0000000016
//  F=> Proto           | string               | [00002] | 0000000024
//  F=> ProtoMajor      | int                  | [00003] | 0000000040
//  F=> ProtoMinor      | int                  | [00004] | 0000000044
//  F=> Header          | http.Header          | [00005] | 0000000048
//  F=> Body            | io.ReadCloser        | [00006] | 0000000056
//  F=> ContentLength   | int64                | [00007] | 0000000072
//  F=> TransferEncodin | []string             | [00008] | 0000000080
//  F=> Close           | bool                 | [00009] | 0000000096
//  F=> Host            | string               | [00010] | 0000000104
//  F=> Form            | url.Values           | [00011] | 0000000120
//  F=> MultipartForm   | *multipart.Form      | [00012] | 0000000128
//  F=> Trailer         | http.Header          | [00013] | 0000000136
//  F=> RemoteAddr      | string               | [00014] | 0000000144
//  F=> RequestURI      | string               | [00015] | 0000000160
//  F=> TLS             | *tls.ConnectionState | [00016] | 0000000176

type Request struct {
	HttpRequest *http.Request
	Body        *bytes.Buffer
	ContentType string
}

func (r *Request) Method() string {
	return r.HttpRequest.Method
}

func (r *Request) URL() string {
	return r.HttpRequest.URL.String()
}

func (r *Request) Params() map[string][]string {
	if r.HttpRequest.Form == nil {
		r.HttpRequest.ParseForm() // ! ...
	}
	return r.HttpRequest.Form
}

func (r *Request) Param(name string) (val []string) {
	if r.HttpRequest.Form == nil {
		r.HttpRequest.ParseForm() // ! ...
	}
	if val, _ = r.HttpRequest.Form[name]; true {
		val = append(val, "")
	}
	return
}
