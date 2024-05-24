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

func genInsert(table Table, withCache bool) (string, string, error) {
	keySet := collection.NewSet[string]()
	uniqueCacheKeySet := collection.NewSet[string]()
	//keySet.Add(table.PrimaryCacheKey.DataKeyExpression)
	//keyVariableSet.Add(table.PrimaryCacheKey.KeyLeft)
	pkCacheKey := table.PrimaryCacheKey.KeyLeft
	pkCacheKeyExpression := table.PrimaryCacheKey.DataKeyExpression
	for _, key := range table.UniqueCacheKey {
		keySet.Add(key.DataKeyExpression)
		uniqueCacheKeySet.Add(key.KeyLeft)
	}
	keys := keySet.Elems()
	sort.Strings(keys)
	keyVars := uniqueCacheKeySet.Elems()
	sort.Strings(keyVars)

	//expressions := make([]string, 0)
	//expressionValues := make([]string, 0)
	//var count int
	//for _, field := range table.Fields {
	//camel := util.SafeString(field.Name.ToCamel())
	//if table.isIgnoreColumns(field.Name.Source()) {
	//	continue
	//}
	//
	//if field.Name.Source() == table.PrimaryKey.Fields[0].Name.Source() {
	//	if table.PrimaryKey.AutoIncrement {
	//		continue
	//	}
	//}

	//count += 1
	//expressions = append(expressions, "?")
	//expressionValues = append(expressionValues, "data."+camel)
	//}

	camel := table.Name.ToCamel()
	text, err := pathext.LoadTemplate(category, insertTemplateFile, template.Insert)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("insert").
		Parse(text).
		Execute(map[string]any{
			"withCache":             withCache,
			"upperStartCamelObject": camel,
			"lowerStartCamelObject": stringx.From(camel).Untitle(),
			"pkCacheKey":            pkCacheKey,
			"pkCacheKeyExpression":  pkCacheKeyExpression,
			"keysLen":               len(keys) + 1,
			//"expression":            strings.Join(expressions, ", "),
			//"expressionValues":      strings.Join(expressionValues, ", "),
			"keys":            strings.Join(keys, "\n"),
			"uniqueCacheKeys": strings.Join(keyVars, ","),
			"data":            table,
		})
	if err != nil {
		return "", "", err
	}

	// interface method
	text, err = pathext.LoadTemplate(category, insertTemplateMethodFile, template.InsertMethod)
	if err != nil {
		return "", "", err
	}

	insertMethodOutput, err := util.With("insertMethod").Parse(text).Execute(map[string]any{
		"upperStartCamelObject": camel,
		"data":                  table,
	})
	if err != nil {
		return "", "", err
	}

	return output.String(), insertMethodOutput.String(), nil
}
