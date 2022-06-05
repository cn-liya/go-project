package types

import "strconv"

type Integer interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8 | ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

func Int2Str[T Integer](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func Float2Str[T float32 | float64](f T, prec int) string {
	return strconv.FormatFloat(float64(f), 'f', prec, 64)
}
