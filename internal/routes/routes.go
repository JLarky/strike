package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// log the error
		fmt.Println(err)
		// handle returned error here.
		w.WriteHeader(503)
		w.Write([]byte("bad"))
	}
}

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world 1123"))
	})
	r.Method("GET", "/index.html", Handler(staticHandler))
	r.Method("GET", "/client.js", Handler(staticHandler2))
	r.Method("GET", "/", Handler(rscHandler))
	r.Method("GET", "/about", Handler(rscHandler))

	return r
}

func staticHandler(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("err")

	if q != "" {
		return errors.New(q)
	}

	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "data/index.html")

	return nil
}

func staticHandler2(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("err")

	if q != "" {
		return errors.New(q)
	}

	w.Header().Set("Content-Type", "text/javascript")
	http.ServeFile(w, r, "data/client.js")

	return nil
}

type HTML string

type Component struct {
	Tag_type string      `json:"tag_type"`
	Props    interface{} `json:"props"`
}

func Text(v interface{}) HTML {
	return HTML("\n" + html.EscapeString(fmt.Sprint(v)))
}

func rscHandler(w http.ResponseWriter, r *http.Request) error {
	page := Component{
		Tag_type: "div",
		Props:    map[string]interface{}{"id": "root", "children": Text("My page" + fmt.Sprint(r.Context().Value(middleware.RequestIDKey)) + " " + r.URL.Path)},
	}

	jsonData, err := json.Marshal(page)

	if err != nil {
		return err
	}

	const htmlStringTpl = `<div id="root">{{.Props.children}}</div>`

	htmlString, err := template.New("htmlString").Parse(htmlStringTpl)

	if err != nil {
		return err
	}

	htmlStringBuf := new(bytes.Buffer)

	err = htmlString.Execute(htmlStringBuf, page)

	if err != nil {
		return err
	}

	q := r.URL.Query().Get("err")

	if q != "" {
		return errors.New(q)
	}

	rsc := r.Header.Get("RSC")
	if rsc == "1" {
		w.Header().Set("Content-Type", "text/x-component; charset=utf-8")
		w.Write(jsonData)
		return nil
	}

	const tpl = `<!DOCTYPE html>
	<html lang="en">
	  <head>
	  <title>{{.Title}}</title>
	  </head>
	  <body>
		{{.HtmlString}}
		<script type="module">
		  import {renderPage} from "./client.js";
		  const x = JSON.parse("{{.JsonData}}");
		  renderPage(x);
		</script>
	  </body>
    </html>`

	t, err := template.New("webpage").Parse(tpl)

	if err != nil {
		return err
	}

	data := struct {
		Title      string
		JsonData   string
		HtmlString template.HTML
	}{
		// My page plus current time in ms
		Title:      "Title",
		JsonData:   string(jsonData),
		HtmlString: template.HTML(htmlStringBuf.String()),
	}

	err = t.Execute(w, data)

	if err != nil {
		return err
	}

	return nil
}

// http hello world 60-80k rps
// my RSC 45k rps
// Fresh hello world 15k
// Fresh with islands 10k
// Next.js app dir 3-4k rps
