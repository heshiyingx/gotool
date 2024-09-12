package snoyflakeid

import (
	"github.com/sony/sonyflake"
	"time"
)

var sf *sonyflake.Sonyflake

func init() {
	var st = sonyflake.Settings{
		StartTime: time.Date(2024, 8, 1, 1, 1, 1, 0, time.UTC),
	}
	sf = sonyflake.NewSonyflake(st)
}
func NextID() (uint64, error) {
	return sf.NextID()
}
