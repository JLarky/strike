package action

import (
	"log"

	"github.com/JLarky/strike/pkg/h"
)

const formComponentName = "strike-form"

// for fallback if form was not hydrated

func Form(comp h.Component) h.Component {
	comp.Tag_type = formComponentName

	if comp.Props["method"] == nil {
		comp.Props["method"] = "POST"
	}
	if comp.Props["encType"] == nil {
		comp.Props["encType"] = "multipart/form-data"
	}
	if comp.Props["action"] == nil {
		log.Println("action is required for a form")
	}
	if val, ok := comp.Props["action"].(Action); ok {
		delete(comp.Props, "action")
		comp.Props["data-$strike-action"] = val.ToActionName()
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

func IsForm(comp h.Component) bool {
	return comp.Tag_type == formComponentName
}
