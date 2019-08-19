// Package util provides some utilities
package util

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

// ConvertStringSliceToIntSlice converts string slices to int slices.
func ConvertStringSliceToIntSlice(stringSlice []string) (intSlice []int, err error) {
	for _, s := range stringSlice {
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, xerrors.Errorf("failed to convert string slice to int slice: %w", err)
		}
		intSlice = append(intSlice, num)
	}
	return
}

func LogPrintfInOneLine(template string, v ...interface{}) {
	text := fmt.Sprintf(template, v...)
	log.Print(strings.Replace(text, "\n", " ", -1))
}
func LogPrintlnInOneLine(v ...interface{}) {
	log.Println(strings.Replace(fmt.Sprint(v...), "\n", " ", -1))
}
