package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/JLarky/strike-notes/server/db"
	"github.com/JLarky/strike/pkg/async"
	. "github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/island"
	"github.com/JLarky/strike/pkg/promise"
	"github.com/JLarky/strike/pkg/strike"
	"github.com/JLarky/strike/pkg/suspense"
)

func main() {
	http.Handle("/favicon.ico", http.FileServer(http.Dir("public")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx, taskChannel, wg := promise.WithContext(ctx)
		fmt.Println("ctx", ctx)
		flush := func() {
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
				time.Sleep(1 * time.Millisecond)
			}
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
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<!doctype html>"))
		// flush()
		s := strike.NewStream(w)
		err := strike.RenderToStream(s, page)
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

		go func() {
			time.Sleep(100 * time.Millisecond)
			close(s.DoneChan)
		}()

		go func() {
			for i := 0; i < 5; i++ {
				wg.Add(1)

				go func(id int) {
					defer wg.Done()

					// Simulating work.
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
					taskChannel <- promise.Task{
						ID:     id,
						Result: fmt.Sprintf("Task %d completed", id),
					}
				}(i)
			}
			wg.Wait()
			close(taskChannel)
		}()

		{
			for task := range taskChannel {
				fmt.Println(task.Result)
				w.Write([]byte("\n<script>console.log('"))
				w.Write([]byte(task.Result))
				w.Write([]byte("')</script>"))
				flush()
			}
		}

		<-s.Done()
	})
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Page(url *url.URL, children Component) Component {
	// time.Sleep(1000 * time.Millisecond)
	fmt.Println("Page", url)
	nav := H("nav",
		H("a", Props{"href": "/"}, "Home"), " ",
		H("a", Props{"href": "/about"}, "About"),
	)

	return H("html", Props{"lang": "en"},
		H("head",
			H("meta", Props{"charset": "utf-8"}),
			H("meta", Props{"name": "description", "content": "React with Server Components demo"}),
			H("meta", Props{"name": "viewport", "content": "width=device-width, initial-scale=1"}),
			H("link", Props{"rel": "stylesheet", "href": "/static/style.css"}),
			H("title", "React Notes"),
			H("style", Props{"type": "text/css"}, []template.HTML{`
			strike-slot {
				display: none;
			}
			strike-island {
				display: contents;
			}
			`}),
		),
		H("body",
			children,
			H("div", Props{"id": "root"}, nav, "Loading... "+url.Path),
			bootstrap(),
		),
	)
}

func bootstrap() Component {
	return H("div",
		H("script", Props{"type": "importmap"}, []template.HTML{`
			{
				"imports": {
					"react": "https://esm.sh/react@canary?dev",
					"react-dom/client": "https://esm.sh/react-dom@canary/client?dev",
					"react/jsx-runtime": "https://esm.sh/react@canary/jsx-runtime?dev",
					"react-error-boundary": "https://esm.sh/react-error-boundary"
				}
			}`}),
		H("script", Props{"async": "async", "type": "module", "src": "/static/strike/bootstrap.js"}),
	)
}

func Wait(c string, ms int) string {
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return c
}

func App2(url *url.URL) Component {
	a := H("div", "Hello", H("div", "World"))
	return a
	// return H(suspense.Suspense, Props{"class": "main", "fallback": a}, Wait("hello", 1000))
	// return H("nav", H(Suspense,
	// 	Props{"fallback": NoteListSkeleton()},
	// 	NodeList(url),
	// ))
}

func App(url *url.URL, ctx context.Context) Component {
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
				Props{"fallback": noteListSkeleton()},
				async.Async(
					ctx,
					func() Component {
						c := NodeList(url)
						fmt.Println("c", c)
						// panic(123)
						return c
					},
				),
			)),
			// func() Component { return NodeList(url) },
			// <Suspense fallback={<NoteListSkeleton />}>
			// 	<NoteList searchText={searchText} />
			// </Suspense>

		),
		H("section", Props{"class": "col note-viewer"}), // 	<Suspense fallback={<NoteSkeleton isEditing={isEditing} />}>
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
			H("div", Props{"class": "sidebar-note-list-item skeleton", "style": "height:5em"}),
		),
		H("li", Props{"class": "v-stack"},
			H("div", Props{"class": "sidebar-note-list-item skeleton", "style": "height:5em"}),
		),
		H("li", Props{"class": "v-stack"},
			H("div", Props{"class": "sidebar-note-list-item skeleton", "style": "height:5em"}),
		),
	))
}

func NodeList(url *url.URL) Component {
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
		noteComponents[i] = H("li", Props{"key": note.Id}, SidebarNote(note))
	}
	return H("ul", Props{"class": "notes-list"}, noteComponents)
}

func SidebarNote(note db.Note) Component {
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

	return Island("SidebarNoteContent", Props{"id": note.Id, "title": note.Title},
		H("div", Props{"class": "sidebar-note-list-item"},
			H("header", Props{"class": "sidebar-note-header"},
				H("strong", note.Title),
				H("small", lastEdited),
			),
			H("button", Props{"class": "sidebar-note-open"}),
		),
		H("strike-slot",
			H("p", Props{"class": "sidebar-note-excerpt"}, H("i", "(No content)")),
		),
	)
}

func Island(componentName string, props Props, children ...any) Component {
	return H("strike-island", props, Props{"component-export": componentName}, children)
}
