package util

const (
	visaBulletinBaseURL = "https://travel.state.gov/content/travel/en/legal/visa-law0/visa-bulletin/%s/visa-bulletin-for-%s-%s.html"

	dateCurrent     = "Current"
	dateUnavailable = "Unavailable"

	monthFormat      = "January"
	yearFormat       = "2006"
	dateParseShort   = "02Jan06"
	dateParseLong    = "02Jan2006"
	dateOutputFormat = "2006-01-02"

	shortCurrent     = "C"
	shortUnavailable = "U"

	countryWorldwide = "Worldwide"
)

var TrackedCountries = []string{
	"All Chargeability Areas Except Those Listed",
	"CHINA-mainland born",
	"INDIA",
	"MEXICO",
	"PHILIPPINES",
}

var CountryNormalization = map[string]string{
	"india":       "India",
	"china":       "China",
	"mexico":      "Mexico",
	"philippines": "Philippines",
}
