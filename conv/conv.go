// Package conv contains predefined convert functions from string to standard types.
package conv

import (
	"errors"
	"strconv"
)

const (
	intBase  = 10
	bitSize  = 0
	bitSize8 = 2 << iota
	bitSize16
	bitSize32
	bitSize64
)

// Bool converts string `s` to bool.
func Bool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// Uint converts string `s` to uint.
func Uint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, intBase, bitSize)
	if err != nil {
		return 0, err
	}

	return uint(val), nil
}

// Uint8 converts string `s` to uint8.
func Uint8(s string) (uint8, error) {
	val, err := strconv.ParseUint(s, intBase, bitSize8)
	if err != nil {
		return 0, err
	}

	return uint8(val), nil
}

// Uint16 converts string `s` to uint16.
func Uint16(s string) (uint16, error) {
	val, err := strconv.ParseUint(s, intBase, bitSize16)
	if err != nil {
		return 0, err
	}

	return uint16(val), nil
}

// Uint32 converts string `s` to uint32.
func Uint32(s string) (uint32, error) {
	val, err := strconv.ParseUint(s, intBase, bitSize32)
	if err != nil {
		return 0, err
	}

	return uint32(val), nil
}

// Uint64 converts string `s` to uint64.
func Uint64(s string) (uint64, error) {
	val, err := strconv.ParseUint(s, intBase, bitSize64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Int converts string `s` to int.
func Int(s string) (int, error) {
	val, err := strconv.ParseInt(s, intBase, bitSize)
	if err != nil {
		return 0, err
	}

	return int(val), nil
}

// Int8 converts string `s` to int8.
func Int8(s string) (int8, error) {
	val, err := strconv.ParseInt(s, intBase, bitSize8)
	if err != nil {
		return 0, err
	}

	return int8(val), nil
}

// Int16 converts string `s` to int16.
func Int16(s string) (int16, error) {
	val, err := strconv.ParseInt(s, intBase, bitSize16)
	if err != nil {
		return 0, err
	}

	return int16(val), nil
}

// Int32 converts string `s` to int32.
func Int32(s string) (int32, error) {
	val, err := strconv.ParseInt(s, intBase, bitSize32)
	if err != nil {
		return 0, err
	}

	return int32(val), nil
}

// Int64 converts string `s` to int64.
func Int64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, intBase, bitSize64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// String converts string `s` to string.
func String(s string) (string, error) {
	return s, nil
}

const maxRune = 1

// Rune converts string `s` to rune.
func Rune(s string) (rune, error) {
	runes := []rune(s)
	if len(runes) != maxRune {
		return 0, errors.New("must be not empty")
	}

	return runes[0], nil
}

// Float32 converts string `s` to float32.
func Float32(s string) (float32, error) {
	val, err := strconv.ParseFloat(s, bitSize32)
	if err != nil {
		return 0, err
	}

	return float32(val), nil
}

// Float64 converts string `s` to float64.
func Float64(s string) (float64, error) {
	val, err := strconv.ParseFloat(s, bitSize64)
	if err != nil {
		return 0, err
	}

	return val, nil
}
