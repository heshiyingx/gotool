package snoyflakeid

import (
	"strconv"
	"testing"
)

func TestNextID(t *testing.T) {
	id, _ := NextID()
	t.Log(strconv.FormatUint(id, 10))
}
