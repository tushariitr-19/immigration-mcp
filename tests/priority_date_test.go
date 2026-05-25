package tests

import (
	"testing"

	"github.com/tushariitr-19/immigration-mcp/util"
)

func TestNormalizeInputCountry(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"India", "India"},
		{"india", "India"},
		{"INDIA", "India"},
		{"China", "China"},
		{"china", "China"},
		{"Mexico", "Mexico"},
		{"Philippines", "Philippines"},
		{"Germany", "Worldwide"},
		{"", "Worldwide"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := util.NormalizeInputCountry(tt.input)
			if got != tt.expected {
				t.Errorf("input %s: expected %s got %s", tt.input, tt.expected, got)
			}
		})
	}
}

func TestNormalizeCategory(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"EB1", "EB1"},
		{"eb1", "EB1"},
		{"Eb1", "EB1"},
		{"EB2", "EB2"},
		{"eb2", "EB2"},
		{"EB3", "EB3"},
		{"eb3", "EB3"},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := util.NormalizeCategory(tt.input)
			if got != tt.expected {
				t.Errorf("input %s: expected %s got %s", tt.input, tt.expected, got)
			}
		})
	}
}
