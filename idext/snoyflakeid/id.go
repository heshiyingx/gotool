package snoyflakeid

import (
	"github.com/sony/sonyflake"
	"time"
)

var sf *sonyflake.Sonyflake

func init() {
	var st = sonyflake.Settings{
		StartTime: time.Now(),
	}
	sf = sonyflake.NewSonyflake(st)
}
func NextID() (uint64, error) {
	return sf.NextID()
}
