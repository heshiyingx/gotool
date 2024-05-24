package gen

import (
	"fmt"
	"github.com/heshiyingx/gotool/dbext/sql/template"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/pathext"
	goformat "go/format"
	"strings"
)

func genImports(table Table, withCache, timeImport bool) (string, error) {
	var thirdImports []string
	var m = map[string]struct{}{}
	for _, c := range table.Fields {
		if len(c.ThirdPkg) > 0 {
			if _, ok := m[c.ThirdPkg]; ok {
				continue
			}
			m[c.ThirdPkg] = struct{}{}
			thirdImports = append(thirdImports, fmt.Sprintf("%q", c.ThirdPkg))
		}
	}

	if withCache {
		text, err := pathext.LoadTemplate(category, importsTemplateFile, template.Imports)
		if err != nil {
			return "", err
		}

		buffer, err := util.With("import").Parse(text).Execute(map[string]any{
			"time":  timeImport,
			"data":  table,
			"third": strings.Join(thirdImports, "\n"),
		})
		if err != nil {
			return "", err
		}
		source, err := goformat.Source(buffer.Bytes())
		if err != nil {
			return "", err
		}
		return string(source), err

	}

	text, err := pathext.LoadTemplate(category, importsWithNoCacheTemplateFile, template.ImportsNoCache)
	if err != nil {
		return "", err
	}

	buffer, err := util.With("import").Parse(text).Execute(map[string]any{
		"time":  timeImport,
		"data":  table,
		"third": strings.Join(thirdImports, "\n"),
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
