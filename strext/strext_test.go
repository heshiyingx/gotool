package strext

import (
	"fmt"
	"testing"
)

func TestToJsonStr(t *testing.T) {
	type Person struct {
		Name string  `json:"name"`
		Age  int     `json:"age"`
		Son  *Person `json:"son"`
	}

	p := Person{
		Name: "John",
		Age:  30,
	}

	jsonStr := ToJsonStr(p)
	jsonStr1 := ToJsonStr(&p)
	jsonStr2 := ToJsonStr(nil)

	fmt.Println(jsonStr)
	fmt.Println(jsonStr1)
	fmt.Println(jsonStr2)
}
