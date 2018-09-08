package controllers

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/schema"
)

const (
	lenSnippet = 20
)

var (
	// reg is the regexp to match the first paragraph of an article
	reg = regexp.MustCompile(`<p.*?>(.*?)</p>`)
)

func parseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}

func parseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(dst, values); err != nil {
		return err
	}
	return nil
}

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}

func generateSnippet(body string) string {
	p := strings.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(p)
	if err != nil {
		return "Sorry, could not generate snippet..."
	}
	a := doc.Find("p").First().Text()
	if len(a) < 150 {
		return "Sorry, I need more text before I can generate a snippet..."
	}
	return string(a[:150]) + "..."
}
