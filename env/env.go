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
	r, err := strconv.Atoi(val)
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
