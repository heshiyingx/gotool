package parser

import "testing"

func TestParse(t *testing.T) {
	parse, err := Parse("/Users/john/study/code/gocode/test2/sqld/user.sql", "database", true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(parse)
}
