package gathering

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestGathering_Gather_Nonexistence(t *testing.T) {
	expected := "lstat nonexistent_path: no such file or directory"
	decoder := &conversion.Jpeg{}
	g := Gatherer{Decoder: decoder}
	_, err := g.Gather("nonexistent_path")
	if actual := err.Error(); actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestGathering_Gather_Unopenable(t *testing.T) {
	tempdir, _ := ioutil.TempDir("", "imgconv")
	defer os.RemoveAll(tempdir)

	path := filepath.Join(tempdir, "unopenable.jpg")
	if _, err := os.OpenFile(path, os.O_CREATE, 000); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	expected := "open " + path + ": permission denied"

	decoder := &conversion.Jpeg{}
	g := Gatherer{Decoder: decoder}
	_, err := g.Gather(tempdir)
	if actual := err.Error(); actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}
