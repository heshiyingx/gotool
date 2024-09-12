package snoyflakeid

import (
	"strconv"
	"testing"
	"time"
)

func TestNextID(t *testing.T) {
	for {
		id, _ := NextID()
		t.Log(strconv.FormatUint(id, 10))
		time.Sleep(time.Second)
	}

}
