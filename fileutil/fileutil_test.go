package fileutil

import (
	"bytes"
	"testing"
)

func TestFileutil_StartsContentsWith(t *testing.T) {
	cases := []struct {
		a        []byte
		b        []byte
		expected bool
	}{
		{a: []byte("\x01\x02"), b: []byte("\x01"), expected: true},
		{a: []byte("\x01\x02"), b: []byte("\x01\x02"), expected: true},
		{a: []byte("\x01\x02"), b: []byte("\x01\x02\x03"), expected: false},
		{a: []byte("\x01\x02"), b: []byte("\x02"), expected: false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			if actual, _ := StartsContentsWith(bytes.NewReader(c.a), c.b); c.expected != actual {
				t.Errorf("expected: %t, actual: %t", c.expected, actual)
			}
		})
	}

}
