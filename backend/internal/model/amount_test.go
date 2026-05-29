package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAmountFromCents(t *testing.T) {
	a := NewAmountFromCents(100)
	assert.Equal(t, Amount(100), a)
	assert.Equal(t, int64(100), a.Cents())
}

func TestNewAmountFromFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected Amount
	}{
		{0.01, 1},
		{0.10, 10},
		{1.00, 100},
		{10.50, 1050},
		{99.99, 9999},
		{123456.78, 12345678},
		{0.00, 0},
	}
	for _, tt := range tests {
		a := NewAmountFromFloat(tt.input)
		assert.Equal(t, tt.expected, a)
	}
}

func TestAmountToFloat(t *testing.T) {
	a := NewAmountFromCents(1250)
	assert.Equal(t, 12.50, a.ToFloat())
}

func TestAmountZero(t *testing.T) {
	a := NewAmountFromCents(0)
	assert.Equal(t, int64(0), a.Cents())
	assert.Equal(t, 0.0, a.ToFloat())
}

func TestAmountNegative(t *testing.T) {
	a := NewAmountFromCents(-500)
	assert.Equal(t, Amount(-500), a)
	assert.Equal(t, int64(-500), a.Cents())
	assert.Equal(t, -5.0, a.ToFloat())
}
