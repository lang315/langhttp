package LangHttp

import (
	"net/http/cookiejar"
	"net/http"
	"time"
	"net/url"
)

type Client struct {
	c           *http.Client
	FixedHeader map[string]string
	Cookie      map[string]string
}

var (
	urlCookie *url.URL
)

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	//proxyUrl, _ := url.Parse("http://127.0.0.1:8888")
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	return &Client{
		c: &http.Client{
			Timeout:   time.Second * 60,
			Jar:       jar,
			//Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		},
	}
}

func (self *Client) Get(urlGet string) (*Response, error) {
	req, err := http.NewRequest("GET", urlGet, nil)
	req.TransferEncoding = []string{"gzip"}
	if err != nil {
		return nil, err
	}

	for name, value := range self.FixedHeader {
		req.Header.Set(name, value)
	}

	for name, value := range self.Cookie {
		req.Header.Set("Cookie", name+"="+value)
	}

	res, err := self.c.Do(req)

	if err != nil {
		return nil, err
	}

	r, err := newResponse(res)
	if err != nil {
		return nil, err
	}

	urlCookie, _ = url.Parse(urlGet)

	return r, nil
}

func (self *Client) Cookies() []*http.Cookie {

	return self.c.Jar.Cookies(urlCookie)
}
