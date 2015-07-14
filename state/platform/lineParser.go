// Copyright 2015 MICHII Shunsuke. All rights reserved.

package platform

import (
	"bytes"
	"errors"
	"unicode"
)

// LineParser reads line-delimited key-value data.
//
// LineParser aims the file which includes key-value attributes in the each
// lines such as following strings.
//
//   Key1="value1"
//   Key2="value2"
//   ...
//
// LineParser has following options.
// Delimiter specifies the rune of key-value delimiter.
//
type LineParser struct {
	Delimiter  rune
	TrimSpaces bool
	TrimQuotes bool
}

// Parse parses the given array of bytes and returns key-value map.
func (p *LineParser) Parse(b []byte) (map[string][]byte, error) {
	m := make(map[string][]byte)
	for i := 0; i < len(b); {
		nextDelim := bytes.IndexRune(b[i:], p.Delimiter)
		if nextDelim < 0 {
			return nil, errors.New("failed to parse, the delimiter is not found")
		}
		key := b[i : i+nextDelim]
		i += nextDelim + 1

		nextLn := bytes.IndexByte(b[i:], '\n')
		if nextLn < 0 {
			return nil, errors.New("failed to parse, the end of line is not found")
		}
		val := b[i : i+nextLn]
		if p.TrimSpaces {
			val = bytes.TrimSpace(val)
		}
		if p.TrimQuotes {
			val = bytes.TrimFunc(val, isQuotationMark)
		}
		m[string(key)] = val
		i += nextLn + 1
	}
	return m, nil
}

func isQuotationMark(r rune) bool {
	return unicode.Is(unicode.Quotation_Mark, r)
}
