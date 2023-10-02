package strike_http

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/*
var static embed.FS

// NewAssetsHandler returns a http.Handler that serves the strike static assets. Example:
//
//	http.Handle("/_strike/", strike_http.NewAssetsHandler())
func NewAssetsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fSys, err := fs.Sub(static, "assets")
		if err != nil {
			panic(err)
		}
		http.StripPrefix("/_strike/", http.FileServer(http.FS(fSys))).ServeHTTP(w, r)
	})
}
