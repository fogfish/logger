package logger

import "unicode/utf8"

const hex = "0123456789abcdef"

var noEscapeTable = [256]bool{}

func init() {
	for i := 0; i <= 0x7e; i++ {
		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
	}
}

func escape(s string) string {
	for i := 0; i < len(s); i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if !noEscapeTable[s[i]] {
			return appendStringComplex(s, i)
		}
	}
	// The string has no need for encoding and therefore is directly used
	return s
}

func appendStringComplex(s string, i int) string {
	dst := ""
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				// In case of error, first append previous simple characters to
				// the byte slice if any and append a replacement character code
				// in place of the invalid sequence.
				if start < i {
					dst = dst + s[start:i]
				}
				dst = dst + `\ufffd`
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if noEscapeTable[b] {
			i++
			continue
		}
		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			dst = dst + s[start:i]
		}
		switch b {
		case '"':
			dst = dst + `\"`
		case '\\':
			dst = dst + `\\`
		case '\b':
			dst = dst + `\b`
		case '\f':
			dst = dst + `\f`
		case '\n':
			dst = dst + `\n`
		case '\r':
			dst = dst + `\r`
		case '\t':
			dst = dst + `\t`
		default:
			dst = dst + `\u00` + string(hex[b>>4]) + string(hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		dst = dst + s[start:]
	}
	return dst
}
