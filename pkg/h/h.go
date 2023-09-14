package h

import (
	"fmt"
	"html/template"

	"github.com/JLarky/strike/pkg/strike"
)

type Component = strike.Component
type Props = strike.Props

func H(tag any, rest ...any) Component {
	props := make(map[string]any)
	children := make([]any, 0)
	for _, item := range rest {
		switch item_type := item.(type) {
		case []template.HTML:
			for _, v := range item_type {
				children = append(children, v)
			}
		case []Component:
			for _, v := range item_type {
				children = append(children, v)
			}
		case []any:
			children = append(children, item_type...)
		case Props:
			for k, v := range item_type {
				props[k] = v
			}
		default:
			children = append(children, item)
		}
	}

	comp := Component{
		Tag_type: "",
		Props:    props,
		Children: children,
	}

	switch tag_type := tag.(type) {
	case string:
		comp.Tag_type = tag_type
		return comp
	case func() Component:
		return tag_type()
	case func(Component) Component:
		return tag_type(comp)
	default:
		panic(fmt.Sprintf("Unsupported type: %T", tag_type))
	}
}
