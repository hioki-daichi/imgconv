package conversion

import (
	"reflect"
	"testing"
)

func TestConversion_Png_MagicBytesSlice(t *testing.T) {
	p := Png{}
	actual := p.MagicBytesSlice()
	expected := [][]byte{[]byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A")}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal(actual, expected)
	}
}

func TestConversion_Png_HasProcessableExtname(t *testing.T) {
	p := Png{}

	cases := []struct {
		path     string
		expected bool
	}{
		{path: "foo.png", expected: true},
		{path: "foo.jpg", expected: false},
		{path: "foo.jpeg", expected: false},
		{path: "foo.gif", expected: false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			if actual := p.HasProcessableExtname(c.path); actual != c.expected {
				t.Errorf("expected: %t, actual: %t", c.expected, actual)
			}
		})
	}
}
