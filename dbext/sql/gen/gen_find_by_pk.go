package gen

import (
	"github.com/heshiyingx/gotool/dbext/sql/template"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/pathext"
	stringx "github.com/heshiyingx/gotool/util/stringext"
	goformat "go/format"
)

func genFindPK(table Table, withCache bool) (string, string, error) {
	camel := table.Name.ToCamel()
	text, err := pathext.LoadTemplate(category, findOneTemplateFile, template.FindByPK)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("find-one-by-pk").
		Parse(text).
		Execute(map[string]any{
			"withCache":                 withCache,
			"upperStartCamelObject":     camel,
			"lowerStartCamelObject":     stringx.From(camel).Untitle(),
			"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Fields[0].Name.Source()),
			"titlePrimaryKey":           table.PrimaryKey.Fields[0].Name.Title(),
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Untitle()),
			"dataType":                  table.PrimaryKey.Fields[0].DataType,
			"cacheKey":                  table.PrimaryCacheKey.KeyExpression,
			"cacheKeyVariable":          table.PrimaryCacheKey.KeyLeft,
			"data":                      table,
		})
	if err != nil {
		return "", "", err
	}

	text, err = pathext.LoadTemplate(category, findOneMethodTemplateFile, template.FindOneMethod)
	if err != nil {
		return "", "", err
	}

	findOneMethod, err := util.With("findOneMethod").
		Parse(text).
		Execute(map[string]any{
			"upperStartCamelObject":     camel,
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Untitle()),
			"titlePrimaryKey":           table.PrimaryKey.Fields[0].Name.Title(),
			"dataType":                  table.PrimaryKey.Fields[0].DataType,
			"data":                      table,
		})
	if err != nil {
		return "", "", err
	}
	findOneCode, err := goformat.Source(output.Bytes())
	if err != nil {
		return "", "", err
	}
	return string(findOneCode), findOneMethod.String(), nil
}
