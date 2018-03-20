package Storage

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
)

var Amount *Numbers.RawAmount

func SetAmount(amm *Numbers.RawAmount) {
	Amount = amm
}

func AddAmount(amm *Numbers.RawAmount) {
	Amount = Amount.Add(amm)
}

func SubtractAmount(amm *Numbers.RawAmount) {
	Amount = Amount.Subtract(amm)
}