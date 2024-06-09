package strategy

import "github.com/common"

type Amplitude struct {
	Status    bool
	Length    int
	KlineData common.KlineData
}

func (signal *Amplitude) Run() {
	if signal.Length < 4 {
		signal.Length = 4
	}
}
