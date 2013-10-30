package gorequests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	// "log"
	"net/http"
	"os"
	"strings"
	"time"
)

// TYPE: http.Response
//      Name            | Type                 | Index   | Offset
// FIELDS               |                      |         |
//  F=> Status          | string               | [00000] | 0000000000
//  F=> StatusCode      | int                  | [00001] | 0000000016
//  F=> Proto           | string               | [00002] | 0000000024
//  F=> ProtoMajor      | int                  | [00003] | 0000000040
//  F=> ProtoMinor      | int                  | [00004] | 0000000044
//  F=> Header          | http.Header          | [00005] | 0000000048
//  F=> Body            | io.ReadCloser        | [00006] | 0000000056
//  F=> ContentLength   | int64                | [00007] | 0000000072
//  F=> TransferEncodin | []string             | [00008] | 0000000080
//  F=> Close           | bool                 | [00009] | 0000000096
//  F=> Trailer         | http.Header          | [00010] | 0000000104
//  F=> Request         | *http.Request        | [00011] | 0000000112

type Response struct {
	HttpResponse        *http.Response
	Status              int
	Request             *Request
	Body                []byte
	BufferReadAndClosed bool // Body could be len() 0...
	Error               error
}

func (r *Response) File(filename string) (fd *os.File, err error) {
	// return an os.File ?
	err = r.IntoFile(filename)
	if err != nil {
		return
	}
	fd, err = os.Open(filename)
	return
}

func (r *Response) IntoFile(filename string) (err error) {
	// check r.BufferReadAndClosed
	// TODO: check content disposition?
	fmt.Println(r.Header("Content-Disposition"))
	out, err := os.Create(filename)
	if err != nil {
		err = errors.New(fmt.Sprintf("could not open outfile (%s): %v", filename, err))
		return
	}
	defer out.Close()

	defer r.HttpResponse.Body.Close()
	_, err = io.Copy(out, r.HttpResponse.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not write outfile (%s): %v", filename, err))
		return
	}
	r.BufferReadAndClosed = true
	return
}

func (r *Response) Buffer() io.ReadCloser {
	// access to the raw buffer with checks whether its closed
	return r.HttpResponse.Body
}

func (r *Response) Text() (body string, err error) {
	err = r.ReadHttpResponse()
	if err != nil {
		err = errors.New(fmt.Sprintf("could not read response body: %v", err))
	}
	body = string(r.Body)
	return
}

func (r *Response) UnmarshalJson(v interface{}) (err error) {
	// TODO: 1. make Text() instead and only use the reader when called
	// TODO: 2. use Text() or the
	if strings.ToLower(r.HttpResponse.Header.Get("Content-Type"))[:16] != "application/json" {
		err = errors.New(fmt.Sprintf("Response body is not JSON: %v", r.HttpResponse.Header.Get("Content-Type")))
		return
	}
	err = r.ReadHttpResponse()
	if err != nil {
		err = errors.New(fmt.Sprintf("json decoding error/could not read response body: %v", err))
		return
	}
	err = json.Unmarshal(r.Body, &v)
	if err != nil {
		err = errors.New(fmt.Sprintf("json decoding error: %v", err))
	}
	return
}

func (r *Response) Headers() http.Header {
	return r.HttpResponse.Header
}

func (r *Response) Header(name string) string {
	return r.HttpResponse.Header.Get(name)
}

func (r *Response) ReadHttpResponse() (err error) {
	if r.BufferReadAndClosed {
		return
	}
	r.Body, err = ioutil.ReadAll(r.HttpResponse.Body)
	defer r.HttpResponse.Body.Close()
	if err != nil {
		err = errors.New(fmt.Sprintf("Error reading response body: %v", err))
	}
	r.BufferReadAndClosed = true
	return
}

func Retry(r *Response, retryCount int, retryTimeout int, retryOnHttpStatus []int) (rr *Response) {
	if retryCount == 0 {
		rr = r
		return
	}

	if retryOnHttpStatus == nil {
		if r.Status == 200 {
			rr = r
			return
		}
	} else {
		for _, s := range retryOnHttpStatus {
			if r.Status == s {
				// log.Println("Status", s, "retry in", retryTimeout, "seconds")
				if retryTimeout > 0 {
					time.Sleep(time.Duration(retryTimeout) * time.Second)
				}
				rr = Retry(do(r.Request.Method(), r.Request.URL(), r.Headers(), r.Request.Body), retryCount-1, retryTimeout, retryOnHttpStatus)
				return
			}
		}
		// none of the statuses for which we want to retry - pass the response on as is
		rr = r
	}

	return
}
