package utils

import (
	"os"
	"strconv"
)


func GetEnvStr(key string) string {
	return os.Getenv(key)
}

func GetEnvInt(key string) int {
	value := os.Getenv(key)
	iValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return iValue
}
