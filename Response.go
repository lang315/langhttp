package LangHttp

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"bytes"
	"github.com/thanhps42/tlib/regex"

)

type Response struct {
	*http.Response

	Bytes []byte

	doc *goquery.Document
	j   *simplejson.Json
}

func newResponse(res *http.Response) (*Response, error) {
	r := &Response{Response: res}
	defer r.Body.Close()

	var err error
	r.Bytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (self *Response) String() string {
	str := string(self.Bytes)
	if str == "" {
		return  ""
	}
	return str
}

func (res *Response) getDoc() (*goquery.Document, error) {
	var err error
	if res.doc == nil {
		res.doc, err = goquery.NewDocumentFromReader(bytes.NewReader(res.Bytes))
		if err != nil {
			return nil, err
		}
	}

	return res.doc, nil
}


func (res *Response) HTMLFindFirst(selector string) *goquery.Selection {
	doc, err := res.getDoc()
	if err != nil {
		return nil
	}
	element := doc.Find(selector)
	if element.Size() == 0 {
		return nil
	}
	return element.First()
}

func (res *Response) HTMLFindLast(selector string) *goquery.Selection {
	doc, err := res.getDoc()
	if err != nil {
		return nil
	}
	element := doc.Find(selector)
	if element.Size() == 0 {
		return nil
	}
	return element.Last()
}

func (res *Response) HTMLFindIndex(selector string, index int) *goquery.Selection {
	doc, err := res.getDoc()
	if err != nil {
		return nil
	}
	element := doc.Find(selector)
	if element.Size() == 0 {
		return nil
	}
	return element.Eq(index)
}

func (res *Response) HTMLFindAll(selector string) []*goquery.Selection {
	doc, err := res.getDoc()
	if err != nil {
		return nil
	}
	element := doc.Find(selector)
	if element.Size() == 0 {
		return nil
	}

	var elems []*goquery.Selection
	element.Each(func(i int, selection *goquery.Selection) {
		elems = append(elems, selection)
	})
	return elems
}



func (res *Response) JSON() (*simplejson.Json, error) {
	var err error
	if res.j == nil {
		res.j, err = simplejson.NewJson(res.Bytes)
		if err != nil {
			return nil, err
		}
	}

	return res.j, nil
}

func (res *Response) Match(expr string) []string {
	value, err := tregex.Match(res.String(), expr)
	if err != nil {
		return nil
	}
	return value
}

func (res *Response) Matches(expr string) ([][]string, error) {
	return tregex.Matches(res.String(), expr)
}

func (res *Response) IsMatch(expr string) bool {
	return tregex.IsMatch(res.String(), expr)
}