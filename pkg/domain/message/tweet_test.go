package message

import (
	"fmt"
	"testing"
)

func TestGetRandomSutabaHighConfidenceText(t *testing.T) {
	fmt.Println(GetRandomSutabaHighConfidenceText())
}
func TestGetRandomSutabaMiddleConfidenceText(t *testing.T) {
	fmt.Println(GetRandomSutabaMiddleConfidenceText())
}
func TestGetRandomSutabaLowConfidenceText(t *testing.T) {
	fmt.Println(GetRandomSutabaLowConfidenceText())
}

func TestGetRandomRamenHighConfidenceText(t *testing.T) {
	fmt.Println(GetRandomRamenHighConfidenceText(0.5))
}
func TestGetRandomRamenMiddleConfidenceText(t *testing.T) {
	fmt.Println(GetRandomRamenMiddleConfidenceText(0.5))
}
func TestGetRandomRamenLowConfidenceText(t *testing.T) {
	fmt.Println(GetRandomRamenLowConfidenceText(0.5))
}

func TestGetRandomOtherHighConfidenceText(t *testing.T) {
	fmt.Println(GetRandomOtherHighConfidenceText(0.25))
}
func TestGetRandomOtherMiddleConfidenceText(t *testing.T) {
	fmt.Println(GetRandomOtherMiddleConfidenceText(0.5))
}
func TestGetRandomOtherLowConfidenceText(t *testing.T) {
	fmt.Println(GetRandomOtherLowConfidenceText())
}
