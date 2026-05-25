package tests

import (
	"strings"
	"testing"

	"github.com/tushariitr-19/immigration-mcp/util"
)

func TestExplainTermExactMatch(t *testing.T) {
	tests := []struct {
		input        string
		expectedTerm string
	}{
		{"priority date", "Priority Date"},
		{"h1b", "H1B Visa"},
		{"ead", "EAD (Employment Authorization Document)"},
		{"perm", "PERM Labor Certification"},
		{"i-485", "Form I-485"},
		{"i-140", "Form I-140"},
		{"ac21", "AC21 / H1B Portability"},
		{"green card", "Green Card (Permanent Resident Card)"},
		{"rfe", "RFE (Request for Evidence)"},
		{"eb2 niw", "EB2 NIW (National Interest Waiver)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			term, ok := util.ImmigrationTerms[strings.ToLower(tt.input)]
			if !ok {
				t.Errorf("term not found: %s", tt.input)
				return
			}
			if term.Term != tt.expectedTerm {
				t.Errorf("expected term %s got %s", tt.expectedTerm, term.Term)
			}
			if term.Simple == "" {
				t.Errorf("expected simple explanation to be non-empty for %s", tt.input)
			}
			if term.Detail == "" {
				t.Errorf("expected detail to be non-empty for %s", tt.input)
			}
		})
	}
}

func TestExplainTermHasRelatedTerms(t *testing.T) {
	tests := []string{
		"priority date",
		"h1b",
		"i-485",
		"perm",
		"ac21",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			term, ok := util.ImmigrationTerms[input]
			if !ok {
				t.Errorf("term not found: %s", input)
				return
			}
			if len(term.RelatedTerms) == 0 {
				t.Errorf("expected related terms for %s", input)
			}
		})
	}
}

func TestExplainTermNotFound(t *testing.T) {
	_, ok := util.ImmigrationTerms["nonexistent term xyz"]
	if ok {
		t.Error("expected term not to be found")
	}
}

func TestAllTermsHaveRequiredFields(t *testing.T) {
	for key, term := range util.ImmigrationTerms {
		t.Run(key, func(t *testing.T) {
			if term.Term == "" {
				t.Errorf("term %s has empty Term field", key)
			}
			if term.Simple == "" {
				t.Errorf("term %s has empty Simple field", key)
			}
			if term.Detail == "" {
				t.Errorf("term %s has empty Detail field", key)
			}
		})
	}
}
