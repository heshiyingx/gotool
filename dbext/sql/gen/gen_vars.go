package gen

import (
	"fmt"
	"github.com/heshiyingx/gotool/dbext/sql/template"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/collection"
	"github.com/heshiyingx/gotool/util/pathext"
	stringx "github.com/heshiyingx/gotool/util/stringext"
	"sort"
	"strings"
)

func genVars(table Table, withCache, postgreSql bool) (string, error) {
	keys := make([]string, 0)
	keys = append(keys, table.PrimaryCacheKey.VarExpression)
	for _, v := range table.UniqueCacheKey {
		keys = append(keys, v.VarExpression)
	}

	camel := table.Name.ToCamel()
	text, err := pathext.LoadTemplate(category, varTemplateFile, template.Vars)
	if err != nil {
		return "", err
	}

	output, err := util.With("var").Parse(text).
		GoFmt(true).Execute(map[string]any{
		"lowerStartCamelObject": stringx.From(camel).Untitle(),
		"upperStartCamelObject": camel,
		"cacheKeys":             strings.Join(keys, "\n"),
		"autoIncrement":         table.PrimaryKey.AutoIncrement,
		"originalPrimaryKey":    wrapWithRawString(table.PrimaryKey.Fields[0].Name.Source()),
		"withCache":             withCache,
		"postgreSql":            postgreSql,
		"data":                  table,
		"ignoreColumns": func() string {
			var set = collection.NewSet[string]()
			for _, c := range table.ignoreColumns {
				if postgreSql {
					set.Add(fmt.Sprintf(`"%s"`, c))
				} else {
					set.Add(fmt.Sprintf("\"`%s`\"", c))
				}
			}
			list := set.Elems()
			sort.Strings(list)
			return strings.Join(list, ", ")
		}(),
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
func wrapWithRawString(v string) string {
	if v == "`" {
		return v
	}

	if !strings.HasPrefix(v, "`") {
		v = "`" + v
	}

	if !strings.HasSuffix(v, "`") {
		v = v + "`"
	} else if len(v) == 1 {
		v = v + "`"
	}

	return v
}
