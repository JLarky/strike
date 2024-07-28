package framework

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/JLarky/strike/pkg/h"
	. "github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/strike"
)

func RscHandler(w http.ResponseWriter, r *http.Request, rsc h.Component) error {
	is_rsc := r.Header.Get("RSC")
	if is_rsc == "1" {
		w.Header().Set("Content-Type", "text/x-component; charset=utf-8")
		return RenderRscStream(w, rsc)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return RenderHtmlDocument(w, rsc)
	}
}

func RenderRscStream(w io.Writer, rsc h.Component) error {
	jsonData, err := json.Marshal(rsc)
	if err != nil {
		return err
	}
	w.Write(jsonData)
	return nil
}

func RenderHtmlDocument(w http.ResponseWriter, rsc h.Component) error {
	jsonBuf := new(bytes.Buffer)
	err := RenderRscStream(jsonBuf, rsc)

	if err != nil {
		return err
	}

	new_rsc := h.UpdateChildren(rsc, func(children []any) []any {
		for _, child := range children {
			child := child.(h.Component)
			if child.Tag_type == "head" {
				rewriteHead(child)
			}
		}
		return children
	})

	htmlStringBuf := new(bytes.Buffer)
	err = strike.RenderToString(htmlStringBuf, new_rsc)
	if err != nil {
		return err
	}

	const tpl = `<!DOCTYPE html>{{.HtmlString}}<script>self.__rsc=self.__rsc||[];__rsc.push({{.JsonData}})</script>`

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
		JsonData:   jsonBuf.String(),
		HtmlString: template.HTML(htmlStringBuf.String()),
	}

	err = t.Execute(w, data)

	if err != nil {
		return err
	}

	return nil
}

// add bootstrap script to head
func rewriteHead(head h.Component) {
	h.UpdateChildren(head, func(children []any) []any {
		for _, child := range Bootstrap() {
			children = append(children, child)
		}
		return children
	})
}

func Bootstrap() []Component {
	// "react": "https://esm.sh/react@canary?dev",
	// "react-dom/client": "https://esm.sh/react-dom@canary/client?dev",
	// "react/jsx-runtime": "https://esm.sh/react@canary/jsx-runtime?dev",

	// "react": "https://esm.sh/react@18.3.0-canary-2807d781a-20230918",
	// "react-dom/client": "https://esm.sh/react-dom@18.3.0-canary-2807d781a-20230918/client",
	// "react/jsx-runtime": "https://esm.sh/react@18.3.0-canary-2807d781a-20230918/jsx-runtime",
	return []Component{
		H("script", Props{"type": "importmap"}, []template.HTML{`
			{
				"imports": {
					"strike_islands": "/static/app/islands.js",
					"react": "https://esm.sh/react@0.0.0-experimental-9ba1bbd65-20230922?dev",
					"react-dom/client": "https://esm.sh/react-dom@0.0.0-experimental-9ba1bbd65-20230922/client?dev",
					"react-dom": "https://esm.sh/react-dom@0.0.0-experimental-9ba1bbd65-20230922?dev",
					"react/jsx-runtime": "https://esm.sh/react@0.0.0-experimental-9ba1bbd65-20230922/jsx-runtime?dev",
					"react-error-boundary": "https://esm.sh/react-error-boundary@4.0.11"
				}
			}`}),
		H("link", Props{"rel": "modulepreload", "href": "/_strike/bootstrap.js"}),
		H("link", Props{"rel": "modulepreload", "href": "https://esm.sh/v132/react-error-boundary@4.0.11/es2022/react-error-boundary.mjs"}),
		H("link", Props{"rel": "modulepreload", "href": "https://esm.sh/react-error-boundary@4.0.11"}),
		H("link", Props{"rel": "modulepreload", "href": "https://esm.sh/react@0.0.0-experimental-9ba1bbd65-20230922?dev"}),
		H("link", Props{"rel": "modulepreload", "href": "https://esm.sh/react-dom@0.0.0-experimental-9ba1bbd65-20230922"}),
		H("link", Props{"rel": "modulepreload", "href": "https://esm.sh/react-dom@0.0.0-experimental-9ba1bbd65-20230922/client"}),
		H("script", Props{"async": "async", "type": "module", "src": "/_strike/bootstrap.js"}),
	}
}
