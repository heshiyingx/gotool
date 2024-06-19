package gen

import (
	"github.com/heshiyingx/gotool/dbext/sql/parser"
	"github.com/heshiyingx/gotool/dbext/sql/template"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/pathext"
	stringx "github.com/heshiyingx/gotool/util/stringext"
)

func genModelCustom(in parser.Table, withCache bool, pkg string) (string, error) {
	text, err := pathext.LoadTemplate(category, modelCustomTemplateFile, template.ModelCustom)
	if err != nil {
		return "", err
	}

	t := util.With("model-custom").
		Parse(text).
		GoFmt(true)
	output, err := t.Execute(map[string]any{
		"pkg":                   pkg,
		"withCache":             withCache,
		"upperStartCamelObject": in.Name.ToCamel(),
		"lowerStartCamelObject": stringx.From(in.Name.ToCamel()).Untitle(),
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
