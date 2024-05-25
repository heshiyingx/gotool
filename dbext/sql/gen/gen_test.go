package gen

import (
	"github.com/heshiyingx/gotool/dbext/sql/parser"
	"testing"
)

func TestNewDefaultGenerator(t *testing.T) {
	generator, err := NewDefaultGenerator(".", nil)
	if err != nil {
		t.Error(err)
		return
	}
	m := make(map[string]*codeTuple)
	filename := "/Users/john/study/code/gocode/test2/sqld/user.sql"
	database := "database"
	strict := true
	withCache := false
	tables, err := parser.Parse(filename, database, strict)
	if err != nil {
		t.Error(err)
		return
	}
	for _, e := range tables {
		gencode, customerCode, err := generator.genModel(*e, withCache)
		if err != nil {
			t.Error(err)
			return
		}
		//customCode, err := generator.genModelCustom(*e, withCache)
		//if err != nil {
		//	return nil, err
		//}

		m[e.Name.Source()] = &codeTuple{
			modelCode:       gencode,
			modelCustomCode: customerCode,
		}
	}
	t.Log(m)
}
func TestDefaultGenerator_StartFromDDL(t *testing.T) {
	generator, err := NewDefaultGenerator(".", nil)
	if err != nil {
		t.Error(err)
		return
	}
	filename := "/Users/john/study/code/gocode/test2/sqld/user.sql"
	//database := "database"
	strict := false
	withCache := true
	err = generator.StartFromDDL(filename, withCache, strict, "ab")
	if err != nil {
		t.Error(err)
		return
	}
}
