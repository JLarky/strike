package framework

import (
	"html/template"

	. "github.com/JLarky/strike/pkg/h"
)

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
					"react": "https://esm.sh/react@0.0.0-experimental-9ba1bbd65-20230922",
					"react-dom/client": "https://esm.sh/react-dom@0.0.0-experimental-9ba1bbd65-20230922/client",
					"react/jsx-runtime": "https://esm.sh/react@0.0.0-experimental-9ba1bbd65-20230922/jsx-runtime",
					"react-error-boundary": "https://esm.sh/react-error-boundary@4.0.11"
				}
			}`}),
		H("link", Props{"rel": "modulepreload", "href": "/_strike/bootstrap.js"}),
		H("link", Props{"rel": "modulepreload", "href": "https://esm.sh/v132/react-error-boundary@4.0.11/es2022/react-error-boundary.mjs"}),
		H("link", Props{"rel": "modulepreload", "href": "https://esm.sh/react-error-boundary@4.0.11"}),
		H("link", Props{"rel": "modulepreload", "href": "https://esm.sh/react@0.0.0-experimental-9ba1bbd65-20230922?dev"}),
		H("link", Props{"rel": "modulepreload", "href": "https://esm.sh/react-dom@0.0.0-experimental-9ba1bbd65-20230922/client"}),
		H("script", Props{"async": "async", "type": "module", "src": "/_strike/bootstrap.js"}),
	}
}
