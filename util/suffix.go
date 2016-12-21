package util

func unitsDigit(num int) (digit int) {
	digit = num % 10
	return
}

func tensDigit(num int) (digit int) {
	digit = (num % 100) / 10
	return
}

func Suffix(num int) (suffix string) {
	if tensDigit(num) == 1 {
		suffix = "th"
		return
	}

	switch unitsDigit(num) {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	default:
		suffix = "th"
	}
	return
}
