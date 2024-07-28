package routes

import (
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/JLarky/strike/pkg/framework"
	"github.com/JLarky/strike/pkg/h"
	. "github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/island"
	"github.com/JLarky/strike/pkg/strike_http"
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
	r.Method("GET", "/static/app/islands.js", Handler(staticHandler2))
	r.Method("GET", "/", Handler(rscHandler))
	r.Method("GET", "/about", Handler(rscHandler))
	r.Method("GET", "/_strike/*", strike_http.NewAssetsHandler())

	return r
}

//go:embed static/*
var static embed.FS

func staticHandler2(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/javascript")
	f, err := static.ReadFile("static/app.js")
	if err != nil {
		return err
	}
	w.Write(f)

	return nil
}

func Island(componentName string, props Props, fallback any) Component {
	return H(island.Island, props, Props{"component-export": componentName, "ssrFallback": fallback})
}

var serverCounter uint64

func rscHandler(w http.ResponseWriter, r *http.Request) error {
	b := make([]byte, 32)
	rand.Read(b)
	sha := sha256.New()
	sha.Write(b)
	sha256 := base64.StdEncoding.EncodeToString([]byte(sha.Sum(nil)))
	nav := H("nav",
		H("a", h.Props{"href": "/"}, "Home"), " ",
		H("a", h.Props{"href": "/about"}, "About"),
	)
	footer := H("footer",
		H("a", Props{"href": "https://github.com/JLarky/strike"}, "see source"),
	)
	var island Component

	if r.URL.Path == "/" {
		island = Island(
			"Counter",
			Props{"serverCounter": serverCounter},
			H("span", "Loading..."),
		)
	} else {
		c := atomic.AddUint64(&serverCounter, 1)
		island = Island(
			"Counter",
			Props{"serverCounter": c},
			H("button", fmt.Sprintf("Count: 0 (%d)", c)),
		)
	}
	body := H("div", h.Props{"id": "root"},
		nav,
		H("div",
			H("div", "My page is "+r.URL.Path),
			H("div", "and I generated this sha256 on the server: "+sha256),
			island,
		),
		footer,
	)

	// page := H("html", Props{"lang": "en", "suppressHydrationWarning": true},
	page := H("html", Props{"lang": "en"},
		H("head",
			H("title", "Title "+r.URL.Path),
		),
		H("body",
			body,
		))

	return framework.RscHandler(w, r, page)
}

// http hello world 60-80k rps
// my RSC 45k rps
// Fresh hello world 15k
// Fresh with islands 10k
// Next.js app dir 3-4k rps
