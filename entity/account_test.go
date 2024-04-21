package entity

import (
	"math"
	"testing"
)

func TestMoney_Dollar(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want int64
	}{
		{
			name: "dollar only",
			m:    "123",
			want: 123,
		},
		{
			name: "cent only",
			m:    "0.123",
			want: 0,
		},
		{
			name: "dollar and cent",
			m:    "99.123",
			want: 99,
		},
		{
			name: "max",
			m:    "18446744073709551615.123",
			want: math.MaxInt64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Dollar(); got != tt.want {
				t.Errorf("Money.Dollar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Cent(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want int64
	}{
		{
			name: "dollar only",
			m:    "123",
			want: 0,
		},
		{
			name: "cent only",
			m:    "0.123",
			want: 12300,
		},
		{
			name: "dollar and cent",
			m:    "99.12345",
			want: 12345,
		},
		{
			name: "exceed precision",
			m:    "99.12399999",
			want: 12399,
		},
		{
			name: "smallest precision",
			m:    "0.00001",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Cent(); got != tt.want {
				t.Errorf("Money.Cent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToMoneyString(t *testing.T) {
	type args struct {
		dollar int64
		cent   int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{dollar: 100, cent: 50000},
			want: "100.50000",
		},
		{
			name: "normal 1",
			args: args{dollar: 100, cent: 5},
			want: "100.00005",
		},
		{
			name: "cent only",
			args: args{dollar: 0, cent: 5},
			want: "0.00005",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToMoneyString(tt.args.dollar, tt.args.cent); got != tt.want {
				t.Errorf("ToMoneyString() = %v, want %v", got, tt.want)
			}
		})
	}
}
