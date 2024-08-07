package strike

import (
	"fmt"
	"html"
	"html/template"
	"io"

	"github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/island"
	"github.com/JLarky/strike/pkg/suspense"
)

type Component = h.Component

type Island struct {
	Component
	ComopnentName string `json:"componentName"`
	Props         Props  `json:"props"`
}

type Props = h.Props

func RenderToString(wr io.Writer, comp Component) error {
	if suspense.IsSuspense(comp) {
		if suspense.CanStream(comp) {
			switch fallback := comp.Props["fallback"].(type) {
			case Component:
				wr.Write([]byte("<!-- Suspense Starts -->"))
				err := RenderToString(wr, fallback)
				wr.Write([]byte("<!-- Suspense Ends -->"))
				return err
			default:
				fmt.Printf("warning: Suspense component is missing fallback prop (got %v instead) in %v", fallback, comp)
				return nil
			}
		} else {
			return RenderChildren(wr, comp.Props["children"].([]any))
		}
	}
	if island.IsIsland(comp) {
		switch fallback := comp.Props["ssrFallback"].(type) {
		case Component:
			return RenderToString(wr, fallback)
		default:
			fmt.Printf("Island component %v is missing ssrFallback prop (got %v instead)", comp, fallback)
			return nil
		}
	}
	if comp.Tag_type == "strike-form" {
		comp.Tag_type = "form"
		delete(comp.Props, "data-$strike-action") // fixes hydration mismatch
		return RenderToString(wr, (comp))
	}
	wr.Write([]byte("<" + comp.Tag_type))
	for prop, value := range comp.Props {
		// Perform a type assertion to convert `value` to a string
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int:
			strValue = fmt.Sprintf("%d", v)
		case uint64:
			strValue = fmt.Sprintf("%d", v)
		case float64:
			strValue = fmt.Sprintf("%f", v)
		case nil:
			strValue = "null"
		case *string:
			if v == nil {
				strValue = "null"
			} else {
				strValue = *v
			}
		case <-chan string:
			strValue = <-v
		case []any:
			continue
		case Component:
			return fmt.Errorf("you can only pass a component as a prop to Island or Suspense components. Error rendering component %v", comp)
		default:
			return fmt.Errorf("cannot convert prop \"%s\" (%v %T) to string", prop, value, value)
		}

		wr.Write([]byte(fmt.Sprintf(` %s="%s"`, html.EscapeString(prop), html.EscapeString(strValue))))
	}
	if comp.Tag_type == "meta" {
		wr.Write([]byte(" />"))
		return nil
	} else {
		wr.Write([]byte(">"))
	}
	if (comp.Props["children"]) != nil {
		err := RenderChildren(wr, comp.Props["children"].([]any))
		if err != nil {
			return fmt.Errorf("cannot run RenderChildren %v %v", err, comp)
		}
	}
	wr.Write([]byte("</" + comp.Tag_type + ">"))
	return nil
}

var childTpl = template.Must(template.New("htmlString").Parse("{{.}}"))

func RenderChildren(wr io.Writer, children []any) error {
	for _, child := range children {
		if child != nil {
			switch childComp := child.(type) {
			case Component:
				err := RenderToString(wr, childComp)
				if err != nil {
					return err
				}
			case func() Component:
				err := RenderToString(wr, childComp())
				if err != nil {
					return err
				}
			case <-chan Component:
				err := RenderToString(wr, <-childComp)
				if err != nil {
					return err
				}
			case string:
				err := childTpl.Execute(wr, child)
				if err != nil {
					return err
				}
			case template.HTML:
				err := childTpl.Execute(wr, child)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("cannot convert child (%v %T) to string", childComp, childComp)
			}
		}
	}
	return nil
}
