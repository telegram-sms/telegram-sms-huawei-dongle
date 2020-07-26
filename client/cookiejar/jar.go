package cookiejar

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// This is jar is to workaround an issue with current version of go.
// For the bug report, see: https://github.com/golang/go/issues/40414

type Jar struct {
	*cookiejar.Jar
}

func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		cookie.Path = "/"
	}

	j.Jar.SetCookies(u, cookies)
}

func New(opts *cookiejar.Options) (*Jar, error) {
	var err error
	jar := &Jar{}
	jar.Jar, err = cookiejar.New(opts)
	return jar, err
}
