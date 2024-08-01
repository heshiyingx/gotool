package strext

import (
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestToJsonStr(t *testing.T) {
	type Person struct {
		Name string  `json:"name"`
		Age  int     `json:"age"`
		Son  *Person `json:"son"`
	}
	p0 := Person{
		Name: "John",
		Age:  30,
	}

	p := Person{
		Name: "John",
		Age:  30,
		Son:  &p0,
	}
	p1 := Person{
		Name: "Amy",
		Age:  1,
		Son:  &p,
	}

	jsonStr := ToJsonStr(p1)
	//jsonStr1 := ToJsonStr(&p)
	//jsonStr2 := ToJsonStr(nil)

	logx.Infof("jsonStr:%v", jsonStr)
}
