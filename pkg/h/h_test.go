package h_test

import (
	"fmt"
	"testing"

	"github.com/JLarky/strike/internal/assert"
	. "github.com/JLarky/strike/pkg/h"
)

func TestEq(t *testing.T) {
	assert.Equal(t, "123", "123")
}

func TestOneChild(t *testing.T) {
	a := H("div", "Hello")
	assert.Equal(t, a.Tag_type, "div")
	assert.Equal(t, 1, len(a.Props))
	children := a.Props["children"].([]any)
	assert.Equal(t, "Hello", children[0])
	assert.Equal(t, 1, len(children))
}

func TestNoChildren(t *testing.T) {
	a := H("div")
	assert.Equal(t, a.Tag_type, "div")
	assert.Equal(t, 0, len(a.Props))
	children := a.Props["children"]
	assert.Equal(t, nil, children)
}

func TestNoChildrenWithProps(t *testing.T) {
	a := H("div", Props{"style": "color: red;"})
	assert.Equal(t, a.Tag_type, "div")
	assert.Equal(t, 1, len(a.Props))
	assert.Equal(t, "color: red;", a.Props["style"])
	children := a.Props["children"]
	assert.Equal(t, nil, children)
}

func TestCustomComponent(t *testing.T) {
	c := func(c Component) Component {
		return H("div", c.Props)
	}
	a := H(c, Props{"style": "color: red;"})
	fmt.Println(a.Props)
	assert.Equal(t, a.Tag_type, "div")
	assert.Equal(t, 1, len(a.Props))
	assert.Equal(t, "color: red;", a.Props["style"])
	children := a.Props["children"]
	assert.Equal(t, nil, children)
}

func TestNilInProps(t *testing.T) {
	a := H("div", Props{"value": nil})
	assert.Equal(t, a.Tag_type, "div")
	assert.Equal(t, 1, len(a.Props))
	assert.Equal(t, nil, a.Props["value"])
	children := a.Props["children"]
	assert.Equal(t, nil, children)
}

func TestSpec(t *testing.T) {
	H("div")
	H("div#id")
	H("div.class.class2")
	H("div", "Hello")
	MyComponent := func() Component {
		return H("div")
	}
	H(MyComponent, "Hello")
	MyComponent2 := func(c Component) Component {
		return H("div")
	}
	H(MyComponent2, "Hello")
	// H(MyComponent, P("test": "Hello"))
	// MyComponent("Hello")
	// H("div", P("style": "color: red;", "class": "test"))
	// H("div", P("style": "color: red;"), P("class": "test"))
	// H("div", P("style": "color: red;"), "Hello", P("class": "test"))
	H("div", Props{"style": "color: red;"}, "Hello", Props{"class": "test"})
}
