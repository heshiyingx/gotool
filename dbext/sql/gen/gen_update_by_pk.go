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

func genUpdateByPK(table Table, withCache bool) (
	string, string, error,
) {
	expressionValues := make([]string, 0)
	pkg := "data."
	if table.ContainsUniqueCacheKey {
		pkg = "newData."
	}
	for _, field := range table.Fields {
		camel := util.SafeString(field.Name.ToCamel())
		if table.isIgnoreColumns(field.Name.Source()) {
			continue
		}

		if field.Name.Source() == table.PrimaryKey.Fields[0].Name.Source() {
			continue
		}

		expressionValues = append(expressionValues, pkg+camel)
	}

	keySet := collection.NewSet[string]()
	keyNamesSet := collection.NewSet[string]()
	keySet.Add(table.PrimaryCacheKey.DataKeyExpression)
	keyNamesSet.Add(table.PrimaryCacheKey.KeyLeft)
	for _, key := range table.UniqueCacheKey {
		keySet.Add(key.DataKeyExpression)
		keyNamesSet.Add(key.KeyLeft)
	}
	keys := keySet.Elems()
	sort.Strings(keys)
	keyNames := keyNamesSet.Elems()
	sort.Strings(keyNames)

	expressionValues = append(
		expressionValues, pkg+table.PrimaryKey.Fields[0].Name.ToCamel(),
	)

	camelTableName := table.Name.ToCamel()
	text, err := pathext.LoadTemplate(category, updateTemplateFile, template.Update)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("update").Parse(text).Execute(
		map[string]any{
			"withCache":                 withCache,
			"containsIndexCache":        table.ContainsUniqueCacheKey,
			"upperStartCamelObject":     camelTableName,
			"keys":                      strings.Join(keys, "\n"),
			"keyNames":                  strings.Join(keyNames, ", "),
			"primaryCacheKey":           table.PrimaryCacheKey.DataKeyExpression,
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Untitle()),
			"dataType":                  table.PrimaryKey.Fields[0].DataType,
			"primaryKeyVariable":        table.PrimaryCacheKey.KeyLeft,
			"titlePrimaryKey":           table.PrimaryKey.Fields[0].Name.Title(),
			"lowerStartCamelObject":     stringx.From(camelTableName).Untitle(),
			"upperStartCamelPrimaryKey": util.EscapeGolangKeyword(
				stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Title(),
			),
			"originalPrimaryKey": wrapWithRawString(
				table.PrimaryKey.Fields[0].Name.Source()),
			"expressionValues": strings.Join(
				expressionValues, ", ",
			),
			"data": table,
		},
	)
	if err != nil {
		return "", "", nil
	}

	// update interface method
	text, err = pathext.LoadTemplate(category, updateMethodTemplateFile, template.UpdateMethod)
	if err != nil {
		return "", "", err
	}

	updateMethodOutput, err := util.With("updateMethod").Parse(text).Execute(
		map[string]any{
			"withCache":                 withCache,
			"containsIndexCache":        table.ContainsUniqueCacheKey,
			"upperStartCamelObject":     camelTableName,
			"keys":                      strings.Join(keys, "\n"),
			"keyNames":                  strings.Join(keyNames, ", "),
			"primaryCacheKey":           table.PrimaryCacheKey.DataKeyExpression,
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Untitle()),
			"dataType":                  table.PrimaryKey.Fields[0].DataType,
			"primaryKeyVariable":        table.PrimaryCacheKey.KeyLeft,
			"titlePrimaryKey":           table.PrimaryKey.Fields[0].Name.Title(),
			"lowerStartCamelObject":     stringx.From(camelTableName).Untitle(),
			"upperStartCamelPrimaryKey": util.EscapeGolangKeyword(
				stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Title(),
			),
			"originalPrimaryKey": wrapWithRawString(
				table.PrimaryKey.Fields[0].Name.Source()),
			"expressionValues": strings.Join(
				expressionValues, ", ",
			),
			"data": table,
		},
	)
	if err != nil {
		return "", "", nil
	}

	return output.String(), updateMethodOutput.String(), nil
}
