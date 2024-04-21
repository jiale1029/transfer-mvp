package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	CentUnit = 10000
)

// Assuming transactions is in dollar and cents
// Support precision up to 1/100000 (e.g. 0.56123)
type Money string

func (m Money) Dollar() int64 {
	dollar, _ := strconv.ParseInt(strings.Split(string(m), ".")[0], 10, 64)
	return dollar
}

// 0.1   -> returns 10000
// 0.56  -> returns 56000
// 0.566 -> returns 56600
func (m Money) Cent() int64 {
	components := strings.Split(string(m), ".")
	if len(components) == 1 { // only dollar
		return 0
	}
	// Pad cents with zeros at the right
	centsStr := components[1]
	for len(centsStr) < 5 {
		centsStr += "0"
	}
	v, _ := strconv.ParseInt(centsStr[:5], 10, 64)
	return v
}

func ToMoneyString(dollar, cent int64) string {
	return fmt.Sprintf("%v.%05s", dollar, strconv.Itoa(int(cent)))
}

type Account struct {
	Id            int64  `gorm:"primaryKey"`
	AccountId     string `gorm:"unique"`
	BalanceDollar int64
	BalanceCent   int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
