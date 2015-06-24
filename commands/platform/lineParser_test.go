package platform_test

import (
	"bytes"
	"testing"

	"github.com/harukasan/orchestra-pit/commands/platform"
)

func TestParse(t *testing.T) {
	input := []byte("key1: yes\nkey2: no\n")
	p := &platform.LineParser{
		Delimiter:  ':',
		TrimSpaces: true,
		TrimQuotes: true,
	}

	m, err := p.Parse(input)
	if err != nil {
		t.Errorf("got error, %v", err)
	}

	k1 := "key1"
	v1 := []byte("yes")
	if !bytes.Equal(m[k1], v1) {
		t.Errorf("%s: got %s, expected: %s", k1, m[k1], v1)
	}

	k2 := "key2"
	v2 := []byte("no")
	if !bytes.Equal(m[k2], v2) {
		t.Errorf("%s: got %s, expected: %s", k2, m[k2], v2)
	}
}

func BenchmarkParse(b *testing.B) {
	input := []byte("key1: yes\nkey2: no\n")
	p := &platform.LineParser{
		Delimiter:  ':',
		TrimSpaces: true,
		TrimQuotes: true,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := p.Parse(input)
		if err != nil {
			b.Errorf("got error, %v", err)
		}
	}
}
