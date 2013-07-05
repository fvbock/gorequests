package gorequests

import (
	"gorequests"
	"testing"
)

type CrocFile struct {
	UUID string `json:uuid`
}

func TestGet(t *testing.T) {
	url := "https://duckduckgo.com"
	r, err := gorequests.Get(url, nil, -1)
	if err != nil {
		t.Fail()
		t.Log(err)
	}
	t.Log(r.Request.Params())

	data, err := gorequests.NewQueryData(
		map[string]string{
			"q": "golang http get request querystring",
		})

	r, err = gorequests.Get(url, data, -1)
	t.Log(r.Status)
	// t.Log(string(r.Text))
	t.Log(r.Request.Params())
	t.Log(r.Request.Param("q"))
	t.Log(r.Request.Param("foo"))
}
