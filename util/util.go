package util

import (
	"fmt"
	"strings"
	"time"
)

func BuildBulletinURL(month, year string) string {
	now := time.Now()
	if month == "" {
		month = strings.ToLower(now.Format(monthFormat))
	}
	if year == "" {
		year = now.Format(yearFormat)
	}
	return fmt.Sprintf(visaBulletinBaseURL, year, strings.ToLower(month), year)
}

func NormalizeCountry(country string) string {
	lower := strings.ToLower(country)
	for key, value := range CountryNormalization {
		if strings.Contains(lower, key) {
			return value
		}
	}
	return countryWorldwide
}

func ParseDate(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == shortCurrent || raw == dateCurrent {
		return dateCurrent
	}
	if raw == shortUnavailable || raw == dateUnavailable {
		return dateUnavailable
	}
	t, err := time.Parse(dateParseShort, raw)
	if err != nil {
		t, err = time.Parse(dateParseLong, raw)
		if err != nil {
			return raw
		}
	}
	return t.Format(dateOutputFormat)
}

func IsEmploymentBasedTable(table [][]string) bool {
	for _, row := range table {
		if len(row) == 0 {
			continue
		}
		label := strings.ToLower(strings.TrimSpace(row[0]))
		if strings.HasPrefix(label, "1st") || strings.HasPrefix(label, "2nd") || strings.HasPrefix(label, "3rd") {
			return true
		}
	}
	return false
}
