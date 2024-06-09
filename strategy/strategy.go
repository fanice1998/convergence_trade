package strategy

import (
	"fmt"

	"github.com/common"
)

type Strategy struct {
	Signal []Signal
}

// type Signal struct {
// 	Amplitude Amplitude
// 	EMA       Ema
// }

type Signal struct {
	Status bool
	Data   []common.KlineData
}

type Ema struct {
	Status bool
}

func (t *Strategy) Run() {
	fmt.Println()
}
