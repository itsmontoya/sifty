package language

func isLetter(char rune) bool {
	switch {
	case char >= 'a' && char <= 'z':
		return true
	case char >= 'A' && char <= 'Z':
		return true
	default:
		return false
	}
}

func isNumber(chars ...rune) (ok bool) {
	for _, char := range chars {
		switch char {
		case '0':
		case '1':
		case '2':
		case '3':
		case '4':
		case '5':
		case '6':
		case '7':
		case '8':
		case '9':

		default:
			return false
		}
	}

	return true
}

func isWhitespace(ch rune) bool {
	switch ch {
	case ' ':
	case '\t':
	case '\n':
	case '\r':
	default:
		return false
	}

	return true
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// Keep these aligned with your identifier policy.
func isWordStart(ch rune) bool {
	return isLetter(ch)
}

func isWordPart(ch rune) bool {
	switch {
	case isLetter(ch):
		return true
	case isDigit(ch):
		return true
	}

	switch ch {
	case '.':
	case '[':
	case ']':
	default:
		return false
	}

	return true
}
