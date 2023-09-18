package strike_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/JLarky/strike/internal/assert"
	"github.com/JLarky/strike/pkg/async"
	. "github.com/JLarky/strike/pkg/h"
	. "github.com/JLarky/strike/pkg/island"
	"github.com/JLarky/strike/pkg/strike"
	"github.com/JLarky/strike/pkg/suspense"
)

func TestOneChild(t *testing.T) {
	a := H("div", "Hello", "World", Props{"style": "color: red;"})
	assert.Equal(t, `<div style="color: red;">HelloWorld</div>`, renderToString(a))
}

func TestChildDeep(t *testing.T) {
	a := H("div", "Hello", "World", Props{"style": "color: red;"})
	b := H("main", a)
	c := H("body", b, H("footer", "Copyright"))
	assert.Equal(t, `<body><main><div style="color: red;">HelloWorld</div></main><footer>Copyright</footer></body>`, renderToString(c))
}

func TestNilProp(t *testing.T) {
	a := H("div", Props{"style": nil})
	assert.Equal(t, `<div style="null"></div>`, renderToString(a))
}

func TestNilStringProp(t *testing.T) {
	var n *string = nil
	a := H("div", Props{"style": n})
	assert.Equal(t, `<div style="null"></div>`, renderToString(a))
}

func TestPromise(t *testing.T) {
	longRunningTask := func() <-chan string {
		r := make(chan string)

		go func() {
			defer close(r)

			// Simulate a workload.
			time.Sleep(time.Millisecond * 100)
			r <- "Hello"
		}()

		return r
	}

	a := H("div", "World", Props{"title": longRunningTask()})
	buf := new(bytes.Buffer)
	err := strike.RenderToString(buf, a)
	assert.Equal(t, err, nil)
	assert.Equal(t, `<div title="Hello">World</div>`, buf.String())
}

func TestPromiseComponent(t *testing.T) {
	longRunningTask := func() <-chan Component {
		r := make(chan Component)

		go func() {
			defer close(r)

			// Simulate a workload.
			time.Sleep(time.Millisecond * 100)
			r <- H("div", "World")
		}()

		return r
	}

	a := H("div", "Hello", longRunningTask())
	assert.Equal(t, `<div>Hello<div>World</div></div>`, renderToString(a))
}

func renderToString(c Component) string {
	buf := new(bytes.Buffer)
	err := strike.RenderToString(buf, c)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func renderToStream(c Component) string {
	buf := new(bytes.Buffer)
	stream := strike.NewStream(buf)
	err := strike.RenderToStream(stream, c)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// func TestAsync(t *testing.T) {
// 	a := H("div", async.Async(func() Component {
// 		time.Sleep(time.Millisecond * 100)
// 		return H("div", "Hello")
// 	}))
// 	assert.Equal(t, `<div><div>Hello</div></div>`, renderToString(a))
// }

func TestSuspenseAsync(t *testing.T) {
	a := H("div", H(
		suspense.Suspense,
		Props{"fallback": H("div", "Loading...")},
		async.Async(func() Component {
			panic("This should not be called")
			time.Sleep(time.Millisecond * 100)
			return H("div", "Hello")
		})),
	)
	assert.Equal(t, `<div><div>Loading...</div></div>`, renderToString(a))
}

func TestPromiseComponent2(t *testing.T) {
	a := H(Island, Props{"ssrFallback": H("div", "Loading...")}, H("div", "Hello"))

	jsonData, err := json.Marshal(a)
	if err != nil {
		panic(fmt.Sprintf("Error serializing data: %v", err))
	}

	assert.Equal(t, `{"$strike":"component","$type":"strike-island","props":{"children":[{"$strike":"component","$type":"div","props":{"children":["Hello"]}}],"ssrFallback":{"$strike":"component","$type":"div","props":{"children":["Loading..."]}}}}`, string(jsonData))

	assert.Equal(t, `<div>Loading...</div>`, renderToString(a))
}

func TestSuspenseComponent(t *testing.T) {
	a := H("div", H(
		suspense.Suspense,
		Props{"fallback": H("div", "Loading...")},
		async.Async(func() Component { panic("Should not be called") }),
	))
	assert.Equal(t, `<div><div>Loading...</div></div>`, renderToString(a))
}

func TestSuspenseComponentToStream(t *testing.T) {
	a := H("div", H(
		suspense.Suspense,
		Props{"fallback": H("div", "Loading...")},
		async.Async(func() Component { panic("Should not be called") }),
	))
	assert.Equal(t, `<div><div>Loading...</div></div>`, renderToStream(a))
}
