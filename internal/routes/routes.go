package routes

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/JLarky/goReactServerComponents/internal/h"
	. "github.com/JLarky/goReactServerComponents/internal/h"
	"github.com/JLarky/goReactServerComponents/internal/strike"
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
	r.Method("GET", "/client.js", Handler(staticHandler2))
	r.Method("GET", "/", Handler(rscHandler))
	r.Method("GET", "/about", Handler(rscHandler))

	return r
}

//go:embed static/*
var static embed.FS

func staticHandler2(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("err")

	if q != "" {
		return errors.New(q)
	}

	w.Header().Set("Content-Type", "text/javascript")
	f, err := static.ReadFile("static/client.js")
	if err != nil {
		return err
	}
	w.Write(f)

	return nil
}

func rscHandler(w http.ResponseWriter, r *http.Request) error {
	nav := H("nav",
		H("a", h.Props{"href": "/"}, "Home"), " ",
		H("a", h.Props{"href": "/about"}, "About"),
	)
	page := H("div", h.Props{"id": "root"},
		nav,
		H("div",
			H("div", "My page is "+r.URL.Path),
			H("div", "and your IP is "+r.RemoteAddr+" (intention is to show that this is server rendered)"),
		),
	)

	jsonData, err := json.Marshal(page)

	if err != nil {
		return err
	}

	htmlStringBuf := new(bytes.Buffer)
	err = strike.RenderToString(htmlStringBuf, page)
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
