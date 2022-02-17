package util

const (
	USD = "USD"
	EUR = "EUR"
	PLN = "PLN"
)

func IsSupportedCurency(currency string) bool {
	switch currency {
	case USD, EUR, PLN:
		return true
	}
	return false
}
