package gathering

import (
	"reflect"
	"testing"

	"github.com/hioki-daichi/imgconv/conversion"
)

func TestGathering_Gather(t *testing.T) {
	cases := []struct {
		decoder  conversion.Decoder
		expected []string
	}{
		{decoder: &conversion.Jpeg{}, expected: []string{"../testdata/jpeg/sample1.jpg", "../testdata/jpeg/sample2.jpg", "../testdata/jpeg/sample3.jpeg"}},
		{decoder: &conversion.Png{}, expected: []string{"../testdata/png/sample1.png", "../testdata/png/sample2.png"}},
		{decoder: &conversion.Gif{}, expected: []string{"../testdata/gif/sample1.gif"}},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			g := Gatherer{Decoder: c.decoder}
			if actual, _ := g.Gather("../testdata/"); !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("expected: %s, actual: %s", c.expected, actual)
			}
		})
	}
}
