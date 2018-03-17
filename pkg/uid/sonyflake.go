package uid

import (
	"github.com/sony/sonyflake"
	"strconv"
)

type SonyFlakeWrapper struct {
	sf *sonyflake.Sonyflake
}

func NewSonyFlake(st sonyflake.Settings) *SonyFlakeWrapper {
	return &SonyFlakeWrapper{sonyflake.NewSonyflake(st)}
}

func (w SonyFlakeWrapper) NextID() (string, error) {
	idUInt, err := w.sf.NextID()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(idUInt, 36), nil
}
