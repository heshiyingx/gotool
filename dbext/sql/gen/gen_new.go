package gen

import (
	"fmt"
	"github.com/heshiyingx/gotool/dbext/sql/template"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/pathext"
)

func genNew(table Table, withCache bool) (string, error) {
	text, err := pathext.LoadTemplate(category, modelNewTemplateFile, template.New)
	if err != nil {
		return "", err
	}

	t := fmt.Sprintf(`"%s"`, wrapWithRawString(table.Name.Source()))

	output, err := util.With("new").
		Parse(text).
		Execute(map[string]any{
			"table":                 t,
			"withCache":             withCache,
			"upperStartCamelObject": table.Name.ToCamel(),
			"data":                  table,
			"pkType":                table.PrimaryKey.Fields[0].DataType,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
