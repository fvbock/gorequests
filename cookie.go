package gorequests

import (
	"net/http"
	"net/url"
	"sync"
)

type CookieJar struct {
	sync.Mutex
	cookies map[string][]*http.Cookie
}

func (cj *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	cj.Lock()
	defer cj.Unlock()

	if cj.cookies == nil {
		cj.Reset()
	}
	cj.cookies[u.Host] = cookies
}

func (cj *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	return cj.cookies[u.Host]
}

func (cj *CookieJar) Reset() {
	cj.cookies = map[string][]*http.Cookie{}
}
