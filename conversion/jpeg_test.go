package conversion

import (
	"reflect"
	"testing"
)

func TestConversion_Jpeg_MagicBytesSlice(t *testing.T) {
	j := Jpeg{}
	actual := j.MagicBytesSlice()
	expected := [][]byte{[]byte("\xFF\xD8\xFF")}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal(actual, expected)
	}
}

func TestConversion_Jpeg_HasProcessableExtname(t *testing.T) {
	j := Jpeg{}

	cases := []struct {
		path     string
		expected bool
	}{
		{path: "foo.jpg", expected: true},
		{path: "foo.jpeg", expected: true},
		{path: "foo.png", expected: false},
		{path: "foo.gif", expected: false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			if actual := j.HasProcessableExtname(c.path); actual != c.expected {
				t.Errorf("expected: %t, actual: %t", c.expected, actual)
			}
		})
	}
}
