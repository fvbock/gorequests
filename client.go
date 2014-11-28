package gorequests

import (
	"crypto/tls"
	"net/http"
)

// var DefaultClient = NewClient()

// func NewClient() http.Client {
// 	o := cookiejar.Options{
// 		PublicSuffixList: publicsuffix.List,
// 	}
// 	jar, _ := cookiejar.New(&o)
// 	return http.Client{Jar: jar}
// }

// TLSClientConfig: &tls.Config{RootCAs: x509.CertPool},
var HttpsTransport = &http.Transport{
	TLSClientConfig: &tls.Config{},
}
