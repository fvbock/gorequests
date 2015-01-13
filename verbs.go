package gorequests

import (
	"bytes"
	"strings"
	// "crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

// func Options(requestURL string, headers, timeout time.Duration) (r *Response) {
// 	r = do("OPTIONS", requestURL, headers, "", nil)
// 	return
// }

/*
Get
*/
func Get(requestURL string, headers http.Header, data map[string][]string, timeout time.Duration) (r *Response) {
	requestURL = prepareQuery(requestURL, data)
	r = do("GET", requestURL, headers, nil)
	return
}

/*
Delete
*/
func Delete(requestURL string, headers http.Header, data map[string][]string, timeout time.Duration) (r *Response) {
	requestURL = prepareQuery(requestURL, data)
	r = do("DELETE", requestURL, headers, nil)
	return
}

/*
Post
*/
func Post(requestURL string, headers http.Header, data interface{}, files map[string]map[string]io.ReadCloser, timeout time.Duration) (r *Response) {
	return postOrPut("POST", requestURL, headers, data, files, timeout)
}

/*
Post Form
*/
func PostForm(requestURL string, headers http.Header, data map[string]string, timeout time.Duration) (r *Response) {
	if headers == nil {
		headers = http.Header{}
	}
	values := make(url.Values)
	for key, val := range data {
		values.Set(key, val)
	}
	headers.Add("Content-Type", "application/x-www-form-urlencoded")
	r = do("POST", requestURL, headers, bytes.NewBuffer([]byte(values.Encode())))
	return
}

/*
Put
*/
func Put(requestURL string, headers http.Header, data interface{}, files map[string]map[string]io.ReadCloser, timeout time.Duration) (r *Response) {
	return postOrPut("PUT", requestURL, headers, data, files, timeout)
}

func postOrPut(verb string, requestURL string, headers http.Header, data interface{}, files map[string]map[string]io.ReadCloser, timeout time.Duration) (r *Response) {
	if headers == nil {
		headers = http.Header{}
	}
	var bodyBuf *bytes.Buffer
	var contentType string
	switch data.(type) {
	case []byte:
		bodyBuf = bytes.NewBuffer(data.([]byte))
	case string:
		bodyBuf = bytes.NewBuffer([]byte(data.(string)))
	case map[string]string:
		contentType, bodyBuf = prepareBody(data.(map[string]string), files)
		headers.Add("Content-Type", contentType)
	}

	r = do(verb, requestURL, headers, bodyBuf)
	return
}

func prepareQuery(requestURL string, data map[string][]string) (reqURL string) {
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
		reqURL = fmt.Sprintf("%s?%s", requestURL, params.Encode())
	} else {
		reqURL = requestURL
	}
	return
}

func prepareBody(data map[string]string, files map[string]map[string]io.ReadCloser) (contentType string, bodyBuf *bytes.Buffer) {
	bodyBuf = bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	for key, val := range data {
		err := bodyWriter.WriteField(key, val)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error writing to multipart buffer: %v", err))
			bodyWriter.Close()
			return
		}
	}

	for key, filesForKey := range files {
		for filename, fs := range filesForKey {
			fileWriter, err := bodyWriter.CreateFormFile(key, filename)
			if err != nil {
				err = errors.New(fmt.Sprintf("Error writing file to buffer: %v", err))
				bodyWriter.Close()
				return
			}
			io.Copy(fileWriter, fs)
			err = fs.Close()
		}
	}

	contentType = bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return
}

// TODO: instead of checking the buffer do something like this
// var encodeUrlMethods = set.NewStringSet([]string{"DELETE", "GET", "HEAD", "OPTIONS"}...)
// var encodeBodyMethods = set.NewStringSet([]string{"PATCH", "POST", "PUT", "TRACE"}...)

func do(method string, requestUrl string, headers http.Header, bodyBuffer *bytes.Buffer) (r *Response) {
	var client = &http.Client{Jar: &CookieJar{}}
	if len(requestUrl) > 5 && strings.ToLower(requestUrl[:5]) == "https" {
		client.Transport = HttpsTransport
	}
	var bufCopy string
	var req *http.Request
	var err error

	if bodyBuffer != nil {
		// this is clunky. i need the stuff to retry a request.
		// do should eventually know that i might want to directly
		// retry the request in case of failures and _only_ then
		// keep this stuff
		bufCopy = bodyBuffer.String()
		req, err = http.NewRequest(method, requestUrl, bodyBuffer)
	} else {
		req, err = http.NewRequest(method, requestUrl, nil)
	}

	for k, vals := range headers {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	// user agent
	req.Header.Set("User-Agent", GR_USER_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		r = &Response{
			Error: errors.New(fmt.Sprintf("Error sending request: %v", err)),
		}
		return
	}

	r = &Response{
		HttpResponse: resp,
		Status:       resp.StatusCode,
		Request: &Request{
			HttpRequest: resp.Request,
			ContentType: headers.Get("Content-Type"),
		},
		Error: err,
	}

	// same as above:
	// should eventually know that i might want to directly
	// retry the request in case of failures and _only_ then
	// keep this stuff
	if bodyBuffer != nil {
		r.Request.Body = bytes.NewBufferString(bufCopy)
	}
	return
}
