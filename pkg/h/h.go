package h

import (
	"encoding/json"
	"fmt"
	"html/template"
)

type Component struct {
	Tag_type string `json:"tag_type"`
	Props    Props  `json:"props"`
}

func (c Component) MarshalJSON() ([]byte, error) {
	if c.Tag_type == "strike-island" {
		return json.Marshal([]any{
			"$strike:island-go",
			c.Props,
		})
	}
	return json.Marshal([]any{
		"$strike:element",
		c.Tag_type,
		c.Props,
	})
}

// this is just for logs, this does not render to HTML
func (c Component) String() string {
	props := ""
	if len(c.Props) > 0 {
		for k, v := range c.Props {
			props += fmt.Sprintf(" %v={%v}", k, v)
		}
	}
	return fmt.Sprintf("<%v%s></%v>", c.Tag_type, props, c.Tag_type)
}

type Props map[string]any

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

	if len(children) > 0 {
		props["children"] = children
	}

	comp := Component{
		Tag_type: "",
		Props:    props,
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
