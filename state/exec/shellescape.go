package exec

var escapeCharTableOffset = int(',')
var escapeCharTable = []bool{
	false, // , = 44
	false, // -
	false, // .
	false, // /
	false, // 0
	false, // 1
	false, // 2
	false, // 3
	false, // 4
	false, // 5
	false, // 6
	false, // 7
	false, // 8
	false, // 9
	false, // :
	false, // ;
	true,  // <
	true,  // =
	true,  // >
	true,  // ?
	false, // @
	false, // A
	false, // B
	false, // C
	false, // D
	false, // E
	false, // F
	false, // G
	false, // H
	false, // I
	false, // J
	false, // K
	false, // L
	false, // M
	false, // N
	false, // O
	false, // P
	false, // Q
	false, // R
	false, // S
	false, // T
	false, // U
	false, // V
	false, // W
	false, // X
	false, // Y
	false, // Z
	true,  // [
	true,  // \
	true,  // ]
	true,  // ^
	false, // _
	true,  // `
	false, // a
	false, // b
	false, // c
	false, // d
	false, // e
	false, // f
	false, // g
	false, // h
	false, // i
	false, // j
	false, // k
	false, // l
	false, // m
	false, // n
	false, // o
	false, // p
	false, // q
	false, // r
	false, // s
	false, // t
	false, // u
	false, // v
	false, // w
	false, // x
	false, // y
	false, // z = 122
}

func needEscape(r rune) bool {
	if ',' <= r && r <= 'z' {
		return escapeCharTable[int(r)-escapeCharTableOffset]
	}
	return true
}

func ShellEscape(s string) string {
	ext := 0
	for _, r := range s {
		if needEscape(r) {
			ext++
		}
	}

	escaped := make([]byte, len(s)+ext)
	i := 0
	for _, r := range s {
		if needEscape(r) {
			escaped[i] = '\\'
			i++
		}
		escaped[i] = byte(r)
		i++
	}
	return string(escaped)
}
