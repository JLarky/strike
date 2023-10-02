package strike_http

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/JLarky/strike/pkg/action"
	"github.com/JLarky/strike/pkg/promise"
	"github.com/JLarky/strike/pkg/strike"
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

// NewRscHandler is a default handler, you don't have to use it. Example:
//
//	http.Handle("/", strike_http.NewRscHandler(serverActions, func(w http.ResponseWriter, r *http.Request, ctx context.Context) strike.Component {
//		return H("div", "Hello, World!")
//	}))
func NewRscHandler(serverActions *action.ServerActions, handler func(w http.ResponseWriter, r *http.Request, ctx context.Context) strike.Component) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxOriginal := r.Context()

		ctx, getChunkCh := promise.WithContext(ctxOriginal)

		flush := func() {
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}

		// debug POST
		if r.Method == "POST" {
			r.ParseMultipartForm(10 << 20)
			action, err := serverActions.ConsumeForm(r.PostForm)
			if err != nil {
				fmt.Printf("Error consuming form: %v", err)
				return
			}
			promise := promise.NewPromise[any](ctx)
			promise.PromiseId = action.ToActionName()
			data, err := action.Action(ctx, r.PostForm)
			promise.ResolveAsync(func() any { return data })
			// FIXME: send errors to client
			fmt.Println(action, data, err)
		}

		page := handler(w, r, ctx)

		rsc := r.Header.Get("RSC")
		if rsc == "1" {
			jsonData, err := json.MarshalIndent(page, "", "  ")
			if err != nil {
				fmt.Printf("Error serializing data: %v", err)
				return
			}

			w.Header().Set("Content-Type", "text/x-component; charset=utf-8")
			w.Write(jsonData)
			w.Write([]byte("\n\n"))
			flush()

			{
				for chunk := range getChunkCh() {
					newEncoder := json.NewEncoder(w)
					newEncoder.SetEscapeHTML(false) // TODO: check if this is safe
					err := newEncoder.Encode(chunk)
					if err != nil {
						fmt.Printf("Error encoding chunk: %v", err)
						return
					}
					w.Write([]byte("\n"))
					flush()
				}
			}

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<!doctype html>"))
		// flush()
		err := strike.RenderToString(w, page)
		if err != nil {
			fmt.Printf("Error rendering page: %v", err)
			return
		}

		// <-skeletonDone

		flush()
		// w.Write([]byte("Hello, World!"))
		jsonData, err := json.MarshalIndent(page, "", "  ")

		if err != nil {
			fmt.Printf("Error serializing data: %v", err)
			return
		}

		const tpl = `<script>self.__rsc=self.__rsc||[];__rsc.push({{.}})</script>`

		t, err := template.New("webpage").Parse(tpl)

		if err != nil {
			fmt.Printf("Error parsing template: %v", err)
			return
		}

		err = t.Execute(w, string(jsonData))

		if err != nil {
			fmt.Printf("Error parsing template: %v", err)
			return
		}

		{ // debug JSX
			w.Write([]byte("\n<template>"))
			w.Write(jsonData)
			w.Write([]byte("</template>"))
		}

		{
			for chunk := range getChunkCh() {
				jsonData, err := json.MarshalIndent(chunk, "", "  ")
				if err != nil {
					fmt.Printf("Error serializing data: %v", err)
					return
				}
				err = t.Execute(w, string(jsonData))
				if err != nil {
					fmt.Printf("Error parsing template: %v", err)
					return
				}
				flush()
			}
		}
	})
}
