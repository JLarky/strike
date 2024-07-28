package framework

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/JLarky/strike/pkg/h"
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

func Bootstrap() []h.Component {
	H := h.H
	bootstrap := "_strike/bootstrap.js"
	react := "https://esm.sh/react@0.0.0-experimental-9ba1bbd65-20230922"
	react_client := "https://esm.sh/react-dom@0.0.0-experimental-9ba1bbd65-20230922/client"
	react_dom := "https://esm.sh/react-dom@0.0.0-experimental-9ba1bbd65-20230922"
	react_jsx := "https://esm.sh/react@0.0.0-experimental-9ba1bbd65-20230922/jsx-runtime"
	react_error_boundary := "https://esm.sh/react-error-boundary@4.0.11"

	// in production PORT will be set by the hosting environment
	is_dev := os.Getenv("PORT") == ""
	if is_dev {
		react += "?dev"
		react_client += "?dev"
		react_dom += "?dev"
		react_jsx += "?dev"
	}

	return []h.Component{
		H("script", h.Props{"type": "importmap"}, []template.HTML{`
			{
				"imports": {
					"strike_islands": "/static/app/islands.js",
					"react": "`, template.HTML(react), `",
					"react-dom/client": "`, template.HTML(react_client), `",
					"react-dom": "`, template.HTML(react_dom), `",
					"react/jsx-runtime": "`, template.HTML(react_jsx), `",
					"react-error-boundary": "`, template.HTML(react_error_boundary), `"
				}
			}`}),
		H("link", h.Props{"rel": "modulepreload", "href": bootstrap}),
		H("link", h.Props{"rel": "modulepreload", "href": react}),
		H("link", h.Props{"rel": "modulepreload", "href": react_client}),
		H("link", h.Props{"rel": "modulepreload", "href": react_dom}),
		H("link", h.Props{"rel": "modulepreload", "href": react_jsx}),
		H("link", h.Props{"rel": "modulepreload", "href": react_error_boundary}),
		H("script", h.Props{"async": "async", "type": "module", "src": bootstrap}),
	}
}
