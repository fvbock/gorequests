gorequests
==========

A wrapper for the golang stf lib http request/response functions.

This is just a very first version that does only have GET and POST support. I will add other verbs, basic auth, SSL control over time.

If you want to help please fork and PR. Any help appreciated.

Documentation does also not exist yet... so here is a hint:

```
// GET + Response Text, Request params, Response status
url := "https://duckduckgo.com"
r, err := gorequests.Get(url, nil, -1)

data, err := gorequests.NewQueryData(
	map[string]string{
		"q": "golang http get request querystring",
	})

r, err = gorequests.Get(url, data, -1)

// response http status code
log.Println(r.Status)

// response body
log.Println(r.Text)

// request params (querystring)
log.Println(r.Request.Params())
log.Println(r.Request.Param("q"))
log.Println(r.Request.Param("foo"))



// GET + JSON
var someType Something

data, err := gorequests.NewQueryData(
    map[string]string{
        "token": TOKEN,
        "stuff": strings.Join(ids, ","),
    })

r, err := gorequests.Get("http://example.com/endpoint", data, -1)
err = r.UnmarshalJson(&someType)



// POST + JSON
var done bool

data := map[string]string{
	"token": TOKEN,
	"uuid":  id,
}

r, err := gorequests.Post("http://example.com/endpoint", data, nil, -1)
err = r.UnmarshalJson(&done)



// GET + save file
r, err := gorequests.Get("http://example.com/endpoint/file", data, -1)

// show reponse headers
log.Println(r.Headers())

err = r.IntoFile("/tmp/foobar")
```
