// This module contains some helper functions for string data.

package str

import "unicode/utf8"

const (
	// The default lowest and highest continuation byte.
	locb = 0x80 // 1000 0000
	hicb = 0xBF // 1011 1111

	// These names of these constants are chosen to give nice alignment in the
	// table below. The first nibble is an index into acceptRanges or F for
	// special one-byte cases. The second nibble is the Rune length or the
	// Status for the special one-byte case.
	xx = 0xF1 // invalid: size 1
	as = 0xF0 // ASCII: size 1
	s1 = 0x02 // accept 0, size 2
	s2 = 0x13 // accept 1, size 3
	s3 = 0x03 // accept 0, size 3
	s4 = 0x23 // accept 2, size 3
	s5 = 0x34 // accept 3, size 4
	s6 = 0x04 // accept 0, size 4
	s7 = 0x44 // accept 4, size 4
)

// first is information about the first byte in a UTF-8 sequence.
var first = [256]uint8{
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x00-0x0F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x10-0x1F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x20-0x2F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x30-0x3F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x40-0x4F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x50-0x5F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x60-0x6F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x70-0x7F
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x80-0x8F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x90-0x9F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xA0-0xAF
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xB0-0xBF
	xx, xx, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xC0-0xCF
	s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xD0-0xDF
	s2, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s4, s3, s3, // 0xE0-0xEF
	s5, s6, s6, s6, s7, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xF0-0xFF
}

// acceptRange gives the range of valid values for the second byte in a UTF-8
// sequence.
type acceptRange struct {
	lo uint8 // lowest value for second byte.
	hi uint8 // highest value for second byte.
}

// acceptRanges has size 16 to avoid bounds checks in the code that uses it.
var acceptRanges = [16]acceptRange{
	0: {locb, hicb},
	1: {0xA0, hicb},
	2: {locb, 0x9F},
	3: {0x90, hicb},
	4: {locb, 0x8F},
}

// Get the index of `n`th rune character.
// If n is bigger than the count of characters, the second response will be false.
func RuneIndex(p []byte, n int) (int, bool) {
	var i int
	for ; n > 0 && i < len(p); n-- {
		if p[i] < utf8.RuneSelf {
			// ASCII fast path
			i++
			continue
		}
		x := first[p[i]]
		if x == xx {
			i++ // invalid.
			continue
		}
		size := int(x & 7)
		if i+size > len(p) {
			i++ // Short or invalid.
			continue
		}
		accept := acceptRanges[x>>4]
		if c := p[i+1]; c < accept.lo || accept.hi < c {
			size = 1
		} else if size == 2 {
		} else if c := p[i+2]; c < locb || hicb < c {
			size = 1
		} else if size == 3 {
		} else if c := p[i+3]; c < locb || hicb < c {
			size = 1
		}
		i += size
	}
	return i, n <= 0
}

// Get the index of `n`th rune character.
// If n is bigger than the count of characters, the second response will be false.
func RuneIndexInString(s string, n int) (int, bool) {
	var i int
	for ; n > 0 && i < len(s); n-- {
		if s[i] < utf8.RuneSelf {
			// ASCII fast path
			i++
			continue
		}
		x := first[s[i]]
		if x == xx {
			i++ // invalid.
			continue
		}
		size := int(x & 7)
		if i+size > len(s) {
			i++ // Short or invalid.
			continue
		}
		accept := acceptRanges[x>>4]
		if c := s[i+1]; c < accept.lo || accept.hi < c {
			size = 1
		} else if size == 2 {
		} else if c := s[i+2]; c < locb || hicb < c {
			size = 1
		} else if size == 3 {
		} else if c := s[i+3]; c < locb || hicb < c {
			size = 1
		}
		i += size
	}
	return i, n <= 0
}

// Get a part from string
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), 0, 3))
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), -3, 3))
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), 0, -3))
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), -6, -3))
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), 100, 10))
// "ABC"
// "OPQ"
// "ABCDEFGHIJKLMN"
// "LMN"
// ""
func RuneSub(p []byte, start, length int) []byte {
	if len(p) == 0 {
		return []byte{}
	}
	if start < 0 {
		start = utf8.RuneCount(p) + start
	}
	if start < 0 {
		return []byte{}
	}
	if start > 0 {
		n, _ := RuneIndex(p, start)
		p = p[n:]
	}
	if len(p) == 0 {
		return []byte{}
	}
	if length == 0 {
		return p
	}
	if length < 0 {
		length = utf8.RuneCount(p) + length
	}
	if length <= 0 {
		return []byte{}
	}
	n, _ := RuneIndex(p, length)
	return p[:n]
}

// Get a part from string
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), 0, 3))
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), -3, 3))
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), 0, -3))
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), -6, -3))
// fmt.Printf("\"%s\"\n", str.RuneSub([]byte("ABCDEFGHIJKLMNOPQ"), 100, 10))
// "ABC"
// "OPQ"
// "ABCDEFGHIJKLMN"
// "LMN"
// ""
func RuneSubString(s string, start, length int) string {
	if s == "" {
		return ""
	}
	if start < 0 {
		start = utf8.RuneCountInString(s) + start
	}
	if start < 0 {
		return ""
	}
	if start > 0 {
		n, _ := RuneIndexInString(s, start)
		s = s[n:]
	}
	if s == "" {
		return ""
	}
	if length == 0 {
		return s
	}
	if length < 0 {
		length = utf8.RuneCountInString(s) + length
	}
	if length <= 0 {
		return ""
	}
	n, _ := RuneIndexInString(s, length)
	return s[:n]
}
