package strike_test

import (
	"bytes"
	"testing"

	"github.com/JLarky/goReactServerComponents/internal/assert"
	. "github.com/JLarky/goReactServerComponents/internal/h"
	"github.com/JLarky/goReactServerComponents/internal/strike"
)

func TestOneChild(t *testing.T) {
	a := H("div", "Hello", "World", Props{"style": "color: red;"})
	buf := new(bytes.Buffer)
	_ = strike.RenderToString(buf, a)
	assert.Equal(t, `<div style="color: red;">HelloWorld</div>`, buf.String())
}

func TestOneChild(t *testing.T) {
	a := H("div", "Hello", "World", Props{"style": "color: red;"})
	buf := new(bytes.Buffer)
	_ = strike.RenderToString(buf, a)
	assert.Equal(t, `<div style="color: red;">HelloWorld</div>`, buf.String())
}
