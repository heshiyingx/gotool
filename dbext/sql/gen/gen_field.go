package gen

import (
	"github.com/heshiyingx/gotool/dbext/sql/parser"
	"github.com/heshiyingx/gotool/dbext/sql/template"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/pathext"
	"strings"
)

func genFields(table Table, fields []*parser.Field) (string, error) {
	var list []string

	for _, field := range fields {
		result, err := genField(table, field)
		if err != nil {
			return "", err
		}

		list = append(list, result)
	}

	return strings.Join(list, "\n"), nil
}

func genField(table Table, field *parser.Field) (string, error) {
	tag, err := genTag(table, field.NameOriginal)
	if err != nil {
		return "", err
	}

	text, err := pathext.LoadTemplate(category, fieldTemplateFile, template.Field)
	if err != nil {
		return "", err
	}

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]any{
			"name":       util.SafeString(field.Name.ToCamel()),
			"type":       field.DataType,
			"tag":        tag,
			"hasComment": field.Comment != "",
			"comment":    field.Comment,
			"data":       table,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
