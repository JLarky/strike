package action

import (
	"log"

	"github.com/JLarky/strike/pkg/h"
)

// for fallback if form was not hydrated

func Form(comp h.Component) h.Component {
	comp.Tag_type = "form"

	if comp.Props["method"] == nil {
		comp.Props["method"] = "POST"
	}
	if comp.Props["enctype"] == nil {
		comp.Props["enctype"] = "multipart/form-data"
	}
	if comp.Props["action"] == nil {
		log.Println("action is required for a form")
	}
	if val, ok := comp.Props["action"].(Action); ok {
		delete(comp.Props, "action")
		if children, ok := comp.Props["children"].([]any); ok {
			comp.Props["children"] = append(children, h.H(HiddenInput, h.Props{"name": val}))
		} else {
			log.Println("prop `children` must be an array")
		}
	} else {
		log.Println("prop `action` must be an action")
	}

	return comp
}
