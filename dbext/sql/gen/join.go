package gen

import (
	stringx "github.com/heshiyingx/gotool/util/stringext"
	"strings"
)

type Join []string

// Title convert items into Title and return
func (j Join) Title() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Title())
	}

	return join
}

// Camel convert items into Camel and return
func (j Join) Camel() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).ToCamel())
	}
	return join
}

// Snake convert items into Snake and return
func (j Join) Snake() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).ToSnake())
	}

	return join
}

// Untitle converts items into Untitle and return
func (j Join) Untitle() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Untitle())
	}

	return join
}

// Upper convert items into Upper and return
func (j Join) Upper() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Upper())
	}

	return join
}

// Lower convert items into Lower and return
func (j Join) Lower() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Lower())
	}

	return join
}

// With convert items into With and return
func (j Join) With(sep string) stringx.String {
	return stringx.From(strings.Join(j, sep))
}
