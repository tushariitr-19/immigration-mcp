package models

// VisaBulletin represents the parsed Visa Bulletin data
type VisaBulletin struct {
	Month           string                   `json:"month"`
	Year            string                   `json:"year"`
	EmploymentBased map[string]CategoryDates `json:"employment_based"`
	PublishedURL    string                   `json:"published_url"`
}

// CategoryDates holds priority dates for each EB category per country
type CategoryDates struct {
	EB1 string `json:"eb1"`
	EB2 string `json:"eb2"`
	EB3 string `json:"eb3"`
}

// PriorityDateResult represents the result of a priority date check
type PriorityDateResult struct {
	Country    string `json:"country"`
	Category   string `json:"category"`
	YourDate   string `json:"your_date"`
	CutoffDate string `json:"cutoff_date"`
	Eligible   bool   `json:"eligible"`
	DaysBehind int    `json:"days_behind,omitempty"`
	Message    string `json:"message"`
}

// USCISNewsItem represents a single USCIS news article
type USCISNewsItem struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	URL     string `json:"url"`
	Summary string `json:"summary"`
}

// ImmigrationTerm represents an explained immigration term
type ImmigrationTerm struct {
	Term         string   `json:"term"`
	Simple       string   `json:"simple"`
	Detail       string   `json:"detail"`
	RelatedTerms []string `json:"related_terms"`
	Forms        []string `json:"applicable_forms,omitempty"`
	Source       string   `json:"source,omitempty"`
}

type ImmigrationTermResult struct {
	Term      ImmigrationTerm `json:"term"`
	MatchType string          `json:"match_type"`
}
