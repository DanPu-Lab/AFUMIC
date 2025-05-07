package core

func GetCode(symbol uint8) int {
	switch symbol {
	case 'A':
		return 0
	case 'G':
		return 1
	case 'C':
		return 2
	case 'T':
		return 3
	default:
		return 4
	}
}

func GetSymbol(code int) uint8 {
	switch code {
	case 0:
		return 'A'
	case 1:
		return 'G'
	case 2:
		return 'C'
	case 3:
		return 'T'
	default:
		return 'N'
	}
}
