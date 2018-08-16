package controllers

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/schema"
)

const (
	lenSnippet = 40
)

var (
	// reg is the regexp to match the first paragraph of an article
	reg = regexp.MustCompile(`<p.*?>(.*?)</p>`)
)

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

// THIS WILL NEED ADJUSTING.....
func generateSnippet(bodyHTML template.HTML) template.HTML {
	p := reg.FindString(string(bodyHTML))
	// fmt.Println(p)
	words := strings.Split(p, " ")
	switch {
	case len(words) <= 1:
		return template.HTML("<p class=\"has-text-grey-light\">" +
			"Sorry, couldn't create a snippet " + "for this article we are " +
			"working on improving this... </p>")
	case len(words) <= lenSnippet:
		return template.HTML("<span class=\"has-text-grey-light\">" +
			strings.Join(words, " ") + "...")
	case len(words) > lenSnippet:
		return template.HTML("<span class=\"has-text-grey-light\">" +
			strings.Join(words[:lenSnippet], " ") + "..." + "</span>")
		// default:
		// 	return template.HTML("<p> Sorry, couldn't create a snippet for this article " +
		// 		"we are working on improving this... </p>")
	}
	return template.HTML("<p class=\"has-text-grey-light\">" +
		"Sorry, couldn't create a snippet " + "for this article we are " +
		"working on improving this... </p>")
}
