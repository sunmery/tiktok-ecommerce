package constants

type Currency string

const (
	UNKNOWN Currency = "UNKNOWN" // 未知
	CNY     Currency = "CNY"     // 人民币
	USD     Currency = "USD"     // 美元
	EUR     Currency = "EUR"     // 欧元
)

func (c Currency) String(currency string) string {
	switch currency {
	case "CNY":
		return string(CNY)
	case "USD":
		return string(USD)
	case "EUR":
		return string(EUR)
	default:
		return string(UNKNOWN)
	}
}

func (c Currency) Currency(currency string) Currency {
	switch currency {
	case "CNY":
		return CNY
	case "USD":
		return USD
	case "EUR":
		return EUR
	default:
		return UNKNOWN
	}
}
