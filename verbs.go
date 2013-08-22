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
	// "os"
	"log"
	"time"
)

// func Options(requestURL string, timeout time.Duration) (r *Response) {
// 	r = do("OPTIONS", requestURL, "", nil)
// 	return
// }

func Get(requestURL string, data map[string][]string, timeout time.Duration) (r *Response) {
	requestURL = prepareQuery(requestURL, data)
	r = do("GET", requestURL, "", nil)
	return
}

func Delete(requestURL string, data map[string][]string, timeout time.Duration) (r *Response) {
	requestURL = prepareQuery(requestURL, data)
	r = do("DELETE", requestURL, "", nil)
	return
}

func Post(url string, data map[string]string, files map[string]map[string]io.ReadCloser, timeout time.Duration) (r *Response) {
	contentType, bodyBuf := prepareBody(data, files)
	r = do("POST", url, contentType, bodyBuf)
	return
}

func Put(url string, data map[string]string, files map[string]map[string]io.ReadCloser, timeout time.Duration) (r *Response) {
	contentType, bodyBuf := prepareBody(data, files)
	r = do("PUT", url, contentType, bodyBuf)
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

func do(method string, url string, bodyType string, bodyBuffer *bytes.Buffer) (r *Response) {
	var client = &http.Client{nil, nil, &CookieJar{}}
	var bufCopy string
	var req *http.Request
	var err error
	if bodyBuffer != nil {
		// this is clunky. i need the stuff to retry a request.
		// do should eventually know that i might want to directly
		// retry the request in case of failures and _only_ then
		// keep this stuff
		bufCopy = bodyBuffer.String()
		req, err = http.NewRequest(method, url, bodyBuffer)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if len(bodyType) > 0 {
		req.Header.Set("Content-Type", bodyType)
	}
	req.Header.Set("User-Agent", GR_USER_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error sending request: %v", err))
		log.Println("err != nil", err)
	}

	r = &Response{
		HttpResponse: resp,
		Status:       resp.StatusCode,
		Request: &Request{
			HttpRequest: resp.Request,
			ContentType: bodyType,
		},
		Error: err,
	}

	// same as above:
	// do should eventually know that i might want to directly
	// retry the request in case of failures and _only_ then
	// keep this stuff
	if bodyBuffer != nil {
		r.Request.Body = bytes.NewBufferString(bufCopy)
	}
	return
}

// POST TODO: form
// func postExample() {
//         values := make(url.Values)
//         values.Set("foo", "bar")
//         r, err := http.PostForm("http://foo.bar.com/ep", values)
//         if err != nil {
//             log.Printf("post error: %s", err)
//             return
//         }
//         body, _ := ioutil.ReadAll(r.Body)
//         r.Body.Close()
//         log.Printf("post repsponse body: %s", body)
// }
