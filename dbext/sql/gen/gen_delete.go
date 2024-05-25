package gen

import (
	"github.com/heshiyingx/gotool/dbext/sql/template"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/collection"
	"github.com/heshiyingx/gotool/util/pathext"
	stringx "github.com/heshiyingx/gotool/util/stringext"
	"sort"
	"strings"
)

func genDeleteByPK(table Table, withCache bool) (string, string, error) {
	keySet := collection.NewSet[string]()
	keyVariableSet := collection.NewSet[string]()
	keySet.Add(table.PrimaryCacheKey.KeyExpression)
	keyVariableSet.Add(table.PrimaryCacheKey.KeyLeft)
	for _, key := range table.UniqueCacheKey {
		keySet.Add(key.DataKeyExpression)
		keyVariableSet.Add(key.KeyLeft)
	}
	keys := keySet.Elems()
	sort.Strings(keys)
	keyVars := keyVariableSet.Elems()
	sort.Strings(keyVars)

	camel := table.Name.ToCamel()
	text, err := pathext.LoadTemplate(category, deleteTemplateFile, template.Delete)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("delete").
		Parse(text).
		Execute(map[string]any{
			"upperStartCamelObject":     camel,
			"withCache":                 withCache,
			"containsIndexCache":        table.ContainsUniqueCacheKey,
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Untitle()),
			"titlePrimaryKey":           table.PrimaryKey.Fields[0].Name.Title(),
			"dataType":                  table.PrimaryKey.Fields[0].DataType,
			"keys":                      strings.Join(keys, "\n"),
			"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Fields[0].Name.Source()),
			"keyValues":                 strings.Join(keyVars, ", "),
			"data":                      table,
		})
	if err != nil {
		return "", "", err
	}

	// interface method
	text, err = pathext.LoadTemplate(category, deleteMethodTemplateFile, template.DeleteMethod)
	if err != nil {
		return "", "", err
	}

	deleteMethodOut, err := util.With("deleteMethod").
		Parse(text).
		Execute(map[string]any{
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Untitle()),
			"dataType":                  table.PrimaryKey.Fields[0].DataType,
			"data":                      table,
			"titlePrimaryKey":           table.PrimaryKey.Fields[0].Name.Title(),
		})
	if err != nil {
		return "", "", err
	}

	return output.String(), deleteMethodOut.String(), nil
}
