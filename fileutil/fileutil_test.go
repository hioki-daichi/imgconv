package fileutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestFileutil_CopyDirRec(t *testing.T) {
	tempdir, _ := ioutil.TempDir("", "imgconv")
	CopyDirRec("../testdata/", tempdir)
	defer os.RemoveAll(tempdir)
	cases := []struct {
		path string
	}{
		{path: "./jpeg/sample1.jpg"},
		{path: "./jpeg/sample2.jpg"},
		{path: "./jpeg/sample3.jpeg"},
		{path: "./png/sample1.png"},
		{path: "./png/sample2.png"},
		{path: "./gif/sample1.gif"},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			if _, err := os.OpenFile(filepath.Join(tempdir, c.path), os.O_CREATE|os.O_EXCL, 0); !os.IsExist(err) {
				t.FailNow()
			}
		})
	}
}
