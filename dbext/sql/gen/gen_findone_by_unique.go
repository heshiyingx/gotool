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

type findOneCode struct {
	findOneMethod          string
	findOneInterfaceMethod string
	cacheExtra             string
}

func genFindAndUpdateOneByUniqueKey(table Table, withCache bool) (*findOneCode, error) {
	text, err := pathext.LoadTemplate(category, findOneByFieldTemplateFile, template.FindOneByField)
	if err != nil {
		return nil, err
	}

	t := util.With("findOneByField").Parse(text)
	var list []string
	camelTableName := table.Name.ToCamel()
	for _, key := range table.UniqueCacheKey {
		keyExpressionSet := collection.NewSet[string]()
		keyNameSet := collection.NewSet[string]()
		//keySet.Add(table.PrimaryCacheKey.DataKeyExpression)
		//keyVariableSet.Add(table.PrimaryCacheKey.KeyLeft)
		//pkCacheKey := table.PrimaryCacheKey.KeyLeft
		keyExpressionSet.Add(table.PrimaryCacheKey.KeyLeft)
		//pkCacheKeyExpression := table.PrimaryCacheKey.DataKeyExpression
		keyNameSet.Add(table.PrimaryCacheKey.DataKeyExpression)
		for _, key := range table.UniqueCacheKey {
			keyNameSet.Add(key.DataKeyExpression)
			keyExpressionSet.Add(key.KeyLeft)
		}
		keys := keyNameSet.Elems()
		sort.Strings(keys)
		keyVars := keyExpressionSet.Elems()
		sort.Strings(keyVars)

		in, paramJoinString, originalFieldString := convertJoin(key)

		output, err := t.Execute(map[string]any{
			"upperStartCamelObject":     camelTableName,
			"upperField":                key.FieldNameJoin.Camel().With("").Source(),
			"in":                        in,
			"withCache":                 withCache,
			"cacheKey":                  key.KeyExpression,
			"cacheKeyVariable":          key.KeyLeft,
			"pkKey":                     fmt.Sprintf(`%v:=fmt.Sprintf("%s%v", %v, %v)`, table.PrimaryCacheKey.KeyLeft, "%s", "%v", table.PrimaryCacheKey.VarLeft, table.PrimaryKey.Fields[0].Name.ToCamel()),
			"pkNameWrap":                wrapWithRawString(table.PrimaryKey.Fields[0].NameOriginal),
			"pkNameType":                table.PrimaryKey.Fields[0].DataType,
			"pkCacheKeyName":            table.PrimaryCacheKey.KeyLeft,
			"lowerStartCamelObject":     stringx.From(camelTableName).Untitle(),
			"lowerStartCamelField":      paramJoinString,
			"upperStartCamelPrimaryKey": table.PrimaryKey.Fields[0].Name.ToCamel(),
			"originalField":             originalFieldString,
			"data":                      table,
			"keys":                      strings.Join(keys, "\n"),
			"keyNames":                  strings.Join(keyVars, ","),
		})
		if err != nil {
			return nil, err
		}

		list = append(list, output.String())
	}

	text, err = pathext.LoadTemplate(category, findOneByFieldMethodTemplateFile, template.FindOneByFieldMethod)
	if err != nil {
		return nil, err
	}

	t = util.With("findOneByFieldMethod").Parse(text)
	var listMethod []string
	for _, key := range table.UniqueCacheKey {
		var inJoin, paramJoin Join
		for _, f := range key.Fields {
			param := util.EscapeGolangKeyword(stringx.From(f.Name.ToCamel()).Untitle())
			inJoin = append(inJoin, fmt.Sprintf("%s %s", param, getNotNullType(f.DataType)))
			paramJoin = append(paramJoin, param)
		}

		var in string
		if len(inJoin) > 0 {
			in = inJoin.With(", ").Source()
		}
		output, err := t.Execute(map[string]any{
			"upperStartCamelObject": camelTableName,
			"upperField":            key.FieldNameJoin.Camel().With("").Source(),
			"in":                    in,
			"data":                  table,
			"pkKey":                 fmt.Sprintf(`%v:=fmt.Sprintf("%s%v", %v, %v)`, table.PrimaryCacheKey.KeyLeft, "%s", "%v", table.PrimaryCacheKey.VarLeft, table.PrimaryCacheKey.Fields[0].NameOriginal),
			"pkNameWrap":            wrapWithRawString(table.PrimaryKey.Fields[0].NameOriginal),
			"pkNameType":            table.PrimaryKey.Fields[0].DataType,
			"pkCacheKeyName":        table.PrimaryCacheKey.KeyLeft,
		})
		if err != nil {
			return nil, err
		}

		listMethod = append(listMethod, output.String())
	}

	if withCache {
		text, err := pathext.LoadTemplate(category, findOneByFieldExtraMethodTemplateFile,
			template.FindOneByFieldExtraMethod)
		if err != nil {
			return nil, err
		}

		out, err := util.With("findOneByFieldExtraMethod").Parse(text).Execute(map[string]any{
			"upperStartCamelObject": camelTableName,
			"primaryKeyLeft":        table.PrimaryCacheKey.VarLeft,
			"lowerStartCamelObject": stringx.From(camelTableName).Untitle(),
			"originalPrimaryField":  wrapWithRawString(table.PrimaryKey.Fields[0].Name.Source()),
			"data":                  table,
		})
		if err != nil {
			return nil, err
		}

		return &findOneCode{
			findOneMethod:          strings.Join(list, "\n"),
			findOneInterfaceMethod: strings.Join(listMethod, "\n"),
			cacheExtra:             out.String(),
		}, nil
	}

	return &findOneCode{
		findOneMethod:          strings.Join(list, "\n"),
		findOneInterfaceMethod: strings.Join(listMethod, "\n"),
	}, nil
}

func convertJoin(key Key) (in, paramJoinString, originalFieldString string) {
	var inJoin, paramJoin, argJoin Join
	for _, f := range key.Fields {
		param := util.EscapeGolangKeyword(stringx.From(f.Name.ToCamel()).Untitle())
		inJoin = append(inJoin, fmt.Sprintf("%s %s", param, getNotNullType(f.DataType)))
		paramJoin = append(paramJoin, param)

		//if postgreSql {
		//	argJoin = append(argJoin, fmt.Sprintf("%s = $%d", wrapWithRawString(f.Name.Source()), index+1))
		//} else {
		argJoin = append(argJoin, fmt.Sprintf("%s = ?", wrapWithRawString(f.Name.Source())))
		//}
	}
	if len(inJoin) > 0 {
		in = inJoin.With(", ").Source()
	}

	if len(paramJoin) > 0 {
		paramJoinString = paramJoin.With(",").Source()
	}

	if len(argJoin) > 0 {
		originalFieldString = argJoin.With(" and ").Source()
	}
	return in, paramJoinString, originalFieldString
}
