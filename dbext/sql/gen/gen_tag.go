package gen

import (
	"github.com/heshiyingx/gotool/dbext/sql/template"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/pathext"
)

func genTag(table Table, in string) (string, error) {
	if in == "" {
		return in, nil
	}

	text, err := pathext.LoadTemplate(category, tagTemplateFile, template.Tag)
	if err != nil {
		return "", err
	}

	output, err := util.With("tag").Parse(text).Execute(map[string]any{
		"field": in,
		"data":  table,
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
