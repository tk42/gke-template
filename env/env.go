package env

import (
	"math"
	"os"
	"strconv"
	"strings"
)

func GetBoolean(key string, defaultValue bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseBool(val)
	if err != nil {
		panic(err)
	}
	return r
}

func GetString(key string, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func GetInt(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(r)
}

func GetInt64(key string, defaultValue int64) int64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return r
}

func GetInt32(key string, defaultValue int32) int32 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		panic(err)
	}
	return r
}

func GetInt16(key string, defaultValue int16) int16 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseInt(val, 10, 16)
	if err != nil {
		panic(err)
	}
	return r
}

func GetInt8(key string, defaultValue int8) int8 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseInt(val, 10, 8)
	if err != nil {
		panic(err)
	}
	return r
}

func GetUint(key string, defaultValue uint) uint {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(r)
}

func GetUint64(key string, defaultValue uint64) uint64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return r
}

func GetUint32(key string, defaultValue uint32) uint32 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		panic(err)
	}
	return r
}

func GetUint16(key string, defaultValue uint16) uint16 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		panic(err)
	}
	return r
}

func GetUint8(key string, defaultValue uint8) uint8 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		panic(err)
	}
	return r
}

func GetInts(key string, delimiter string, defaultValue []int) []int {
	vals, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	var ints []int
	for _, s := range strings.Split(vals, delimiter) {
		val, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		ints = append(ints, val)
	}
	return ints
}

func GetFloat64(key string, defaultValue float64) float64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return math.NaN()
	}
	return r
}

func GetFloat32(key string, defaultValue float32) float32 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	r, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return math.NaN()
	}
	return r
}
