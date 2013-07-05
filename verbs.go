package gorequests

import (
	"bytes"
	// "crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Get(request_url string, data map[string][]string, timeout time.Duration) (r Response, err error) {
	if len(data) > 0 {
		params := url.Values{}
		for key, vals := range data {
			if len(vals) == 1 {
				params.Set(key, vals[0])
			} else {
				for _, v := range vals {
					params.Add(key, v)
				}
			}
		}
		request_url = fmt.Sprintf("%s?%s", request_url, params.Encode())
	}

	var client = &http.Client{nil, nil, &CookieJar{}}
	resp, err := client.Get(request_url)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error sending request: %v", err))
		return r, err
	}
	r.HttpResponse = resp
	r.Status = r.HttpResponse.StatusCode
	r.Request = &Request{HttpRequest: r.HttpResponse.Request}
	return
}

func Post(url string, data map[string]string, files map[string]*os.File, timeout time.Duration) (r Response, err error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	for key, val := range data {
		err := body_writer.WriteField(key, val)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error writing to multipart buffer: %v", err))
			return r, err
		}
	}

	for key, fh := range files {
		file_writer, err := body_writer.CreateFormFile(key, fh.Name())
		if err != nil {
			err = errors.New(fmt.Sprintf("Error writing file to buffer: %v", err))
			return r, err
		}
		io.Copy(file_writer, fh)
		fh.Close()
	}

	content_type := body_writer.FormDataContentType()
	body_writer.Close()

	var client = &http.Client{nil, nil, &CookieJar{}}
	resp, err := client.Post(url, content_type, body_buf)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error sending request: %v", err))
		return r, err
	}
	r.HttpResponse = resp
	r.Status = r.HttpResponse.StatusCode
	r.Request = &Request{HttpRequest: r.HttpResponse.Request}
	return
}
