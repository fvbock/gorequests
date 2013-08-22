package gorequests

import (
	"github.com/fvbock/gorequests"
	"io"
	"os"
	"testing"
)

var (
	// requestbinURL = "http://requestb.in/151i6401"
	requestbinURL = "https://k2xgquot4wuv.runscope.net"
)

type CrocFile struct {
	UUID string `json:uuid`
}

func TestGet(t *testing.T) {
	url := requestbinURL
	r := gorequests.Get(url, nil, -1)
	if r.Error != nil {
		t.Fail()
		t.Log(r.Error)
	}
	t.Log(r.Request.Params())
	t.Log(r.Request.URL())

	data, err := gorequests.NewQueryData(
		map[string]string{
			"q": "golang http get request querystring",
		})

	if err != nil {
		t.Fail()
		t.Log(err)
	}

	// data, err := gorequests.NewQueryData(
	// 	map[string][]string{
	// 		"q": []string{"golang http get request querystring"},
	// 		"t": []string{"canonical", "foobar"},
	// 	})

	r = gorequests.Get(url, data, -1)
	t.Log(r.Status)
	t.Log(r.Request.URL())
	t.Log(r.Request.Params())
	t.Log(r.Request.Param("q"))
	t.Log(r.Request.Param("foo"))
	t.Log(r.Error)
}

func TestPost(t *testing.T) {
	url := requestbinURL
	data := map[string]string{
		"token": "foobar",
	}

	// files, err := gorequests.NewFilesMap(map[string]string{
	// 	"file": "Sample.doc",
	// })
	// if err != nil {
	// 	t.Fail()
	// 	t.Log(err)
	// }

	fh, err := os.Open("testfile.txt")
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	files := map[string]map[string]io.ReadCloser{
		"file": map[string]io.ReadCloser{"file": fh},
	}

	r := gorequests.Post(url, data, files, -1)
	if r.Error != nil {
		t.Fail()
		t.Log(r.Error)
	}

	t.Log(r)
	// t.Log(r.HttpR)
	// t.Log(r.Text)
	// t.Log(r.HttpR.Request.TransferEncoding)
	// var cf CrocFile
	// r.UnmarshalJson(&cf)
	// t.Log(cf)
	t.Log(r.Status)
	t.Log(r.Headers())
	t.Log(r.Header("Date"))
	t.Log(r.Header("Connection"))
	t.Log(r.Header("Server"))
	t.Log(r.Header("Content-Type"))
	t.Log(r.Header("FOO"))
}

func TestRetry(t *testing.T) {
	url := requestbinURL
	data := map[string]string{
		"token":      "foobar",
		"other_data": "1234567890",
	}

	fh, err := os.Open("testfile.txt")
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	files := map[string]map[string]io.ReadCloser{
		"file": map[string]io.ReadCloser{"file": fh},
	}

	r := gorequests.Retry(gorequests.Post(url, data, files, -1), 3, 5, []int{200})
	if r.Error != nil {
		t.Fail()
		t.Log(r.Error)
	}

	t.Log(r)
	t.Log(r.Status)
	t.Log(r.Headers())
}

func TestPut(t *testing.T) {
	url := requestbinURL
	data := map[string]string{
		"token": "foobar",
	}

	fh, err := os.Open("testfile.txt")
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	files := map[string]map[string]io.ReadCloser{
		"file": map[string]io.ReadCloser{"file": fh},
	}

	r := gorequests.Put(url, data, files, -1)
	if r.Error != nil {
		t.Fail()
		t.Log(r.Error)
	}

	t.Log(r)
	t.Log(r.Status)
	t.Log(r.Headers())
	t.Log(r.Header("Date"))
	t.Log(r.Header("Connection"))
	t.Log(r.Header("Server"))
	t.Log(r.Header("Content-Type"))
	t.Log(r.Header("FOO"))
}

func TestDelete(t *testing.T) {
	url := requestbinURL
	data, err := gorequests.NewQueryData(
		map[string]string{
			"id": "foobar",
		})
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	r := gorequests.Delete(url, data, -1)
	if r.Error != nil {
		t.Fail()
		t.Log(r.Error)
	}

	t.Log(r)
	t.Log(r.Status)
	t.Log(r.Headers())
	t.Log(r.Header("Date"))
	t.Log(r.Header("Connection"))
	t.Log(r.Header("Server"))
	t.Log(r.Header("Content-Type"))
}

// func TestOptions(t *testing.T) {
// 	url := requestbinURL

// 	r := gorequests.Options(url, -1)
// 	if r.Error != nil {
// 		t.Fail()
// 		t.Log(r.Error)
// 	}

// 	t.Log(r)
// 	t.Log(r.Status)
// 	t.Log(r.Headers())
// 	t.Log(r.Header("Date"))
// 	t.Log(r.Header("Connection"))
// 	t.Log(r.Header("Server"))
// 	t.Log(r.Header("Content-Type"))
// }
