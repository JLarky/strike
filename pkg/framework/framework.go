package framework

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"

	. "github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/strike"
)

func RscHandler(w http.ResponseWriter, r *http.Request, page Component) error {
	jsonData, err := json.Marshal(page)

	if err != nil {
		return err
	}

	htmlStringBuf := new(bytes.Buffer)
	err = strike.RenderToString(htmlStringBuf, page)
	if err != nil {
		return err
	}

	rsc := r.Header.Get("RSC")
	if rsc == "1" {
		w.Header().Set("Content-Type", "text/x-component; charset=utf-8")
		w.Write(jsonData)
		return nil
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
		JsonData:   string(jsonData),
		HtmlString: template.HTML(htmlStringBuf.String()),
	}

	err = t.Execute(w, data)

	if err != nil {
		return err
	}

	return nil
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
