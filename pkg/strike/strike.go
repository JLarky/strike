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
		switch fallback := comp.Props["fallback"].(type) {
		case Component:
			return RenderToString(wr, fallback)
		default:
			fmt.Printf("Suspense component %v is missing fallback prop (got %v instead)", comp, fallback)
			return nil
		}
	}
	if island.IsIsland(comp) {
		switch fallback := comp.Props["ssrFallback"].(type) {
		case Component:
			return RenderToString(wr, fallback)
		default:
			fmt.Printf("Suspense component %v is missing fallback prop (got %v instead)", comp, fallback)
			return nil
		}
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
			return fmt.Errorf("cannot convert prop %s (%v %T) to string", prop, value, value)
		}

		wr.Write([]byte(fmt.Sprintf(` %s="%s"`, html.EscapeString(prop), html.EscapeString(strValue))))
	}
	wr.Write([]byte(">"))
	childTpl, err := template.New("htmlString").Parse("{{.}}")
	if err != nil {
		return err
	}
	if (comp.Props["children"]) != nil {
		children := comp.Props["children"].([]any)
		for _, child := range children {
			if child != nil {
				switch childComp := child.(type) {
				case Component:
					err = RenderToString(wr, childComp)
					if err != nil {
						return err
					}
				case func() Component:
					err = RenderToString(wr, childComp())
					if err != nil {
						return err
					}
				case <-chan Component:
					err = RenderToString(wr, <-childComp)
					if err != nil {
						return err
					}
				case string:
					err = childTpl.Execute(wr, child)
					if err != nil {
						return err
					}
				case template.HTML:
					err = childTpl.Execute(wr, child)
					if err != nil {
						return err
					}
				default:
					return fmt.Errorf("cannot convert prop (%v %T) to string for %v", childComp, childComp, comp)
				}
			}
		}
	}
	wr.Write([]byte("</" + comp.Tag_type + ">"))
	return nil
}

type Stream struct {
	io.Writer
	DoneChan chan bool
}

func (s Stream) Done() <-chan bool {
	return s.DoneChan
}

func NewStream(wr io.Writer) Stream {
	return Stream{wr, make(chan bool)}
}

func RenderToStream(wr Stream, comp Component) error {
	if suspense.IsSuspense(comp) {
		switch fallback := comp.Props["fallback"].(type) {
		case Component:
			return RenderToStream(wr, fallback)
		default:
			fmt.Printf("Suspense component %v is missing fallback prop (got %v instead)", comp, fallback)
			return nil
		}
	}
	if island.IsIsland(comp) {
		switch fallback := comp.Props["ssrFallback"].(type) {
		case Component:
			return RenderToString(wr, fallback)
		default:
			fmt.Printf("Suspense component %v is missing fallback prop (got %v instead)", comp, fallback)
			return nil
		}
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
		default:
			return fmt.Errorf("cannot convert prop %s (%v %T) to string", prop, value, value)
		}

		wr.Write([]byte(fmt.Sprintf(` %s="%s"`, html.EscapeString(prop), html.EscapeString(strValue))))
	}
	wr.Write([]byte(">"))
	childTpl, err := template.New("htmlString").Parse("{{.}}")
	if err != nil {
		return err
	}
	if (comp.Props["children"]) != nil {
		children := comp.Props["children"].([]any)
		for _, child := range children {
			if child != nil {
				switch childComp := child.(type) {
				case Component:
					err = RenderToStream(wr, childComp)
					if err != nil {
						return err
					}
				case func() Component:
					err = RenderToStream(wr, childComp())
					if err != nil {
						return err
					}
				case <-chan Component:
					err = RenderToStream(wr, <-childComp)
					if err != nil {
						return err
					}
				case string:
					err = childTpl.Execute(wr, child)
					if err != nil {
						return err
					}
				case template.HTML:
					err = childTpl.Execute(wr, child)
					if err != nil {
						return err
					}
				default:
					return fmt.Errorf("cannot convert prop (%v %T) to string for %v", childComp, childComp, comp)
				}
			}
		}
	}
	wr.Write([]byte("</" + comp.Tag_type + ">"))
	return nil
}
