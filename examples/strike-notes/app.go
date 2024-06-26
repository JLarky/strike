package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"time"

	_ "net/http/pprof"

	"github.com/JLarky/strike-notes/server/db"
	"github.com/JLarky/strike/pkg/action"
	"github.com/JLarky/strike/pkg/framework"
	. "github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/island"
	"github.com/JLarky/strike/pkg/promise"
	"github.com/JLarky/strike/pkg/strike"
	"github.com/JLarky/strike/pkg/strike_http"
	"github.com/JLarky/strike/pkg/suspense"
)

var lastForm = ""

var useStreaming = true

//go:embed public/*
var static embed.FS

var serverActions = action.NewServerActions()

func main() {
	serverActions.Register("test123", action.ActionFunc(func(ctx context.Context, args url.Values) (any, error) {
		lastForm = fmt.Sprintf("%v", args)
		fmt.Println("test123", args)
		promise := promise.NewPromise[any](ctx)
		promise.ResolveAsync(func() any {
			time.Sleep(1000 * time.Millisecond)
			return "data"
		})
		return promise, nil
	}))

	http.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fSys, err := fs.Sub(static, "public")
		if err != nil {
			panic(err)
		}
		http.FileServer(http.FS(fSys)).ServeHTTP(w, r)
	}))
	http.Handle("/_strike/", strike_http.NewAssetsHandler())
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fSys, err := fs.Sub(static, "public")
		if err != nil {
			panic(err)
		}
		http.StripPrefix("/static/", http.FileServer(http.FS(fSys))).ServeHTTP(w, r)
	}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctxOriginal := r.Context()

		ctx, getChunkCh := promise.WithContext(ctxOriginal)
		if !useStreaming {
			ctx = ctxOriginal
		}

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

		page := Page(
			r.URL,
			App(r.URL, ctx),
		)

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
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Page(url *url.URL, children Component) Component {
	// time.Sleep(1000 * time.Millisecond)
	fmt.Println("Page", url)

	return H("html", Props{"lang": "en"},
		H("head",
			H("meta", Props{"charset": "utf-8"}),
			H("meta", Props{"name": "description", "content": "React with Server Components demo"}),
			H("meta", Props{"name": "viewport", "content": "width=device-width, initial-scale=1"}),
			H("link", Props{"rel": "stylesheet", "href": "/static/style.css"}),
			H("title", "React Notes"),
			framework.Bootstrap(),
		),
		H("body",
			children,
		),
	)
}

func App(url *url.URL, ctx context.Context) Component {
	props := Props{}
	if useStreaming {
		props["ctx"] = ctx
	}
	p := promise.NewPromise[Component](ctx)
	p2 := promise.NewPromise[Component](ctx)
	return H("div", Props{"class": "main"},
		H("section", Props{"class": "col sidebar"},
			H("section", Props{"class": "sidebar-header"},
				H("img", Props{"class": "logo", "src": "/static/logo.svg", "width": "22px", "height": "20px", "alt": "", "role": "presentation"}),
				H("strong", "React Notes"),
			),
			H("section", Props{"class": "sidebar-menu", "role": "menubar"},
				searchField(url),
				editButton(nil, "New"),
			),
			H("nav", H(suspense.Suspense,
				props,
				Props{"fallback": noteListSkeleton(), "p": p, "p2": p2},
				func() Component {
					return nodeList(url)
				},
			)),
		),
		H(action.Form, Props{"action": serverActions.GetOrFail("test123")},
			H("input", Props{"type": "text", "name": "test"}),
			H(island.Island, Props{"component-export": "SubmitButton"},
				Props{"myAct": serverActions.GetOrFail("test123")},
				H("button", Props{"type": "submit"}, "Submit"),
			),
		),
		H("section", Props{"class": "col note-viewer"},
			H("div", Props{"class": "note--empty-state"},
				H("span", Props{"class": "note-text--empty-state"}, "Click a note on the left to view something! 🥺"+url.String()+lastForm),
			),
		),
		// 	<Suspense fallback={<NoteSkeleton isEditing={isEditing} />}>
		// 	<Note selectedId={selectedId} isEditing={isEditing} />
		// </Suspense>
	)
}

func searchField(url *url.URL) Component {
	q := url.Query().Get("q")
	return H(island.Island, Props{
		"component-export": "SearchField",
		"ssrFallback": H("form", Props{"class": "search", "role": "search"},
			H("label", Props{"class": "offscreen"}),
			H("input", Props{"placeholder": "Search", "value": q, "disabled": "disabled"}),
		)})
}

func editButton(noteId *string, title string) Component {
	ssrFallback := H("button", Props{"class": "edit-button edit-button--solid", "role": "menuitem"}, title)
	return H(island.Island, Props{
		"component-export": "EditButton",
		"noteId":           noteId,
		"ssrFallback":      ssrFallback,
	}, title)
}

func noteListSkeleton() Component {
	return H("div", H("ul", Props{"class": "notes-list skeleton-container"},
		H("li", Props{"class": "v-stack"},
			H("div", Props{"class": "sidebar-note-list-item skeleton", "style": "height:84px"}),
		),
		H("li", Props{"class": "v-stack"},
			H("div", Props{"class": "sidebar-note-list-item skeleton", "style": "height:84px"}),
		),
		H("li", Props{"class": "v-stack"},
			H("div", Props{"class": "sidebar-note-list-item skeleton", "style": "height:84px"}),
		),
	))
}

func nodeList(url *url.URL) Component {
	q := url.Query().Get("q")
	notes, err := db.SearchNotes(q)
	if err != nil {
		panic(fmt.Sprintf("Error searching notes: %v", err))
	}
	if (len(notes)) == 0 {
		text := "No notes created yet!"
		if q != "" {
			text = fmt.Sprintf(`Couldn't find any notes titled "%s".`, q)
		}
		return H("div", Props{"class": "notes-empty"}, text)
	}
	noteComponents := make([]Component, len(notes))
	for i, note := range notes {
		noteComponents[i] = H("li", Props{"key": note.Id}, sidebarNote(note))
	}
	return H("ul", Props{"class": "notes-list"}, noteComponents)
}

func sidebarNote(note db.Note) Component {
	isToday := func(t time.Time) bool {
		now := time.Now()
		return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
	}

	lastEdited := ""

	if isToday(note.UpdatedAt) {
		lastEdited = note.UpdatedAt.Format("3:04 PM")
	} else {
		lastEdited = note.UpdatedAt.Format("1/_2/06")
	}

	children :=
		H("header", Props{"class": "sidebar-note-header"},
			H("strong", note.Title),
			H("small", lastEdited),
		)

	return H(island.Island,
		Props{
			"id": note.Id, "title": note.Title,
			"component-export": "SidebarNoteContent",
			"ssrFallback": H("div",
				Props{"class": "sidebar-note-list-item"},
				children,
				H("button", Props{"class": "sidebar-note-open"}),
			),
			"expandedChildren": H("p", Props{"class": "sidebar-note-excerpt"}, H("i", "(No content)")),
		},
		children,
	)
}
