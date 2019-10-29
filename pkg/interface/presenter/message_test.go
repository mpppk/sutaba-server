package presenter

import (
	"fmt"
	"testing"

	"github.com/mpppk/messagen/messagen"
)

func Test_getMessagenDefinitions(t *testing.T) {
	classTypes := []string{ClassSutabaValue, ClassRamenValue, ClassOtherValue}
	confidenceTypes := []string{ConfidenceHighValue, ConfidenceMediumValue, ConfidenceLowValue}
	debugTypes := []string{DebugOnValue, DebugOffValue}

	var states []map[string]string
	for _, cType := range classTypes {
		for _, confType := range confidenceTypes {
			for _, dType := range debugTypes {
				states = append(states, map[string]string{
					ClassType:      cType,
					ConfidenceType: confType,
					DebugType:      dType,
					RuleNumType:    "999",
				})
			}
		}
	}

	definitions := getMessagenDefinitions()
	generator, err := messagen.New(nil)
	if err != nil {
		t.Errorf("error occurred when generate new messagen instance %v", err)
	}
	if err := generator.AddDefinition(definitions...); err != nil {
		t.Errorf("error occurred when add definitions: %v", err)
	}
	for _, state := range states {
		messages, err := generator.Generate("Root", state, 1)
		if err != nil {
			t.Errorf("error occurred when generate message: state: %v, err: %v", state, err)
		}
		fmt.Println(state, messages[0])
	}
}
