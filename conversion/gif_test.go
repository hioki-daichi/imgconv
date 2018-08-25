package conversion

import (
	"reflect"
	"testing"
)

func TestConversion_Gif_MagicBytesSlice(t *testing.T) {
	g := Gif{}
	actual := g.MagicBytesSlice()
	expected := [][]byte{[]byte("GIF87a"), []byte("GIF89a")}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal(actual, expected)
	}
}

func TestConversion_Gif_HasProcessableExtname(t *testing.T) {
	g := Gif{}

	cases := []struct {
		path     string
		expected bool
	}{
		{path: "foo.gif", expected: true},
		{path: "foo.jpg", expected: false},
		{path: "foo.jpeg", expected: false},
		{path: "foo.png", expected: false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			if actual := g.HasProcessableExtname(c.path); actual != c.expected {
				t.Errorf("expected: %t, actual: %t", c.expected, actual)
			}
		})
	}
}
