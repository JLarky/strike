package strike

import (
	"fmt"
	"html"
	"html/template"
	"io"
)

type Component struct {
	Tag_type string `json:"tag_type"`
	Props    Props  `json:"props"`
	Children []any  `json:"children"`
}

type Props map[string]any

func RenderToString(wr io.Writer, comp Component) error {
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
	for _, child := range comp.Children {
		if child != nil {
			switch childComp := child.(type) {
			case Component:
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
