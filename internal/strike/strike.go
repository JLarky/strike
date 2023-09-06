package strike

import (
	"fmt"
	"html"
	"html/template"
	"io"

	"github.com/JLarky/strike/internal/h"
)

func RenderToString(wr io.Writer, comp h.Component) error {
	wr.Write([]byte("<" + comp.Tag_type))
	for prop, value := range comp.Props {
		// Perform a type assertion to convert `value` to a string
		strValue, ok := value.(string)
		if !ok {
			return fmt.Errorf("cannot convert prop %s to string", prop)
		}

		wr.Write([]byte(fmt.Sprintf(` %s="%s"`, html.EscapeString(prop), html.EscapeString(strValue))))
	}
	wr.Write([]byte(">"))
	childTpl, err := template.New("htmlString").Parse("{{.}}")
	if err != nil {
		return err
	}
	for _, child := range comp.Children {
		if child != nil {
			switch childComp := child.(type) {
			case h.Component:
				err = RenderToString(wr, childComp)
				if err != nil {
					return err
				}
			default:
				err = childTpl.Execute(wr, child)
				if err != nil {
					return err
				}
			}
		}
	}
	wr.Write([]byte("</" + comp.Tag_type + ">"))
	return nil
}
