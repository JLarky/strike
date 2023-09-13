package strike_test

import (
	"bytes"
	"testing"

	"github.com/JLarky/strike/internal/assert"
	. "github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/strike"
)

func TestOneChild(t *testing.T) {
	a := H("div", "Hello", "World", Props{"style": "color: red;"})
	buf := new(bytes.Buffer)
	_ = strike.RenderToString(buf, a)
	assert.Equal(t, `<div style="color: red;">HelloWorld</div>`, buf.String())
}

func TestChildDeep(t *testing.T) {
	a := H("div", "Hello", "World", Props{"style": "color: red;"})
	b := H("main", a)
	c := H("body", b, H("footer", "Copyright"))
	buf := new(bytes.Buffer)
	_ = strike.RenderToString(buf, c)
	assert.Equal(t, `<body><main><div style="color: red;">HelloWorld</div></main><footer>Copyright</footer></body>`, buf.String())
}

func TestNilProp(t *testing.T) {
	a := H("div", Props{"style": nil})
	buf := new(bytes.Buffer)
	err := strike.RenderToString(buf, a)
	assert.Equal(t, err, nil)
	assert.Equal(t, `<div style="null"></div>`, buf.String())
}

func TestNilStringProp(t *testing.T) {
	var n *string = nil
	a := H("div", Props{"style": n})
	buf := new(bytes.Buffer)
	err := strike.RenderToString(buf, a)
	assert.Equal(t, err, nil)
	assert.Equal(t, `<div style="null"></div>`, buf.String())
}
