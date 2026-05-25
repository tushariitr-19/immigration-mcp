package tests

import (
	"testing"

	"github.com/tushariitr-19/immigration-mcp/util"
)

func TestBuildBulletinURL(t *testing.T) {
	tests := []struct {
		month    string
		year     string
		expected string
	}{
		{
			month:    "may",
			year:     "2026",
			expected: "https://travel.state.gov/content/travel/en/legal/visa-law0/visa-bulletin/2026/visa-bulletin-for-may-2026.html",
		},
		{
			month:    "january",
			year:     "2025",
			expected: "https://travel.state.gov/content/travel/en/legal/visa-law0/visa-bulletin/2025/visa-bulletin-for-january-2025.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.month+"-"+tt.year, func(t *testing.T) {
			got := util.BuildBulletinURL(tt.month, tt.year)
			if got != tt.expected {
				t.Errorf("expected %s got %s", tt.expected, got)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"01APR23", "2023-04-01"},
		{"15JUL14", "2014-07-15"},
		{"C", "Current"},
		{"Current", "Current"},
		{"U", "Unavailable"},
		{"Unavailable", "Unavailable"},
		{"01JAN2026", "2026-01-01"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := util.ParseDate(tt.input)
			if got != tt.expected {
				t.Errorf("input %s: expected %s got %s", tt.input, tt.expected, got)
			}
		})
	}
}

func TestNormalizeCountry(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"INDIA", "India"},
		{"CHINA-mainland born", "China"},
		{"MEXICO", "Mexico"},
		{"PHILIPPINES", "Philippines"},
		{"All Chargeability Areas", "Worldwide"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := util.NormalizeCountry(tt.input)
			if got != tt.expected {
				t.Errorf("input %s: expected %s got %s", tt.input, tt.expected, got)
			}
		})
	}
}

func TestIsEmploymentBasedTable(t *testing.T) {
	tests := []struct {
		name     string
		table    [][]string
		expected bool
	}{
		{
			name: "valid EB table",
			table: [][]string{
				{"Employment", "Worldwide", "INDIA", "CHINA"},
				{"1st", "Current", "2023-04-01", "2023-04-01"},
				{"2nd", "Current", "2014-07-15", "2021-09-01"},
				{"3rd", "Current", "2013-11-15", "2021-06-15"},
			},
			expected: true,
		},
		{
			name: "family sponsored table",
			table: [][]string{
				{"Family", "Worldwide", "INDIA"},
				{"F1", "01SEP17", "01SEP17"},
				{"F2A", "01AUG24", "01AUG24"},
			},
			expected: false,
		},
		{
			name:     "empty table",
			table:    [][]string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.IsEmploymentBasedTable(tt.table)
			if got != tt.expected {
				t.Errorf("expected %v got %v", tt.expected, got)
			}
		})
	}
}
