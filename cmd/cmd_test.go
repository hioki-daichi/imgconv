package cmd

import (
	"bytes"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/hioki-daichi/imgconv/conversion"
	"github.com/hioki-daichi/imgconv/fileutil"
)

func TestCmd_Run(t *testing.T) {
	t.Parallel()
	jpegDecoder := &conversion.Jpeg{}
	pngDecoder := &conversion.Png{}
	gifDecoder := &conversion.Gif{}

	jpegEncoder := &conversion.Jpeg{Options: &jpeg.Options{Quality: 1}}
	pngEncoder := &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.NoCompression}}
	gifEncoder := &conversion.Gif{Options: &gif.Options{NumColors: 1}}

	cases := []struct {
		decoder  conversion.Decoder
		encoder  conversion.Encoder
		force    bool
		expected func(string) string
	}{
		{decoder: jpegDecoder, encoder: pngEncoder, force: true, expected: func(tempdir string) string {
			return `Converted: "` + tempdir + `/jpeg/sample1.png"\nConverted: "` + tempdir + `/jpeg/sample2.png"\nConverted: "` + tempdir + `/jpeg/sample3.png"\n`
		}},
		{decoder: jpegDecoder, encoder: gifEncoder, force: true, expected: func(tempdir string) string {
			return `Converted: "` + tempdir + `/jpeg/sample1.gif"\nConverted: "` + tempdir + `/jpeg/sample2.gif"\nConverted: "` + tempdir + `/jpeg/sample3.gif"\n`
		}},
		{decoder: pngDecoder, encoder: jpegEncoder, force: true, expected: func(tempdir string) string {
			return `Converted: "` + tempdir + `/png/sample1.jpg"\nConverted: "` + tempdir + `/png/sample2.jpg"\n`
		}},
		{decoder: pngDecoder, encoder: gifEncoder, force: true, expected: func(tempdir string) string {
			return `Converted: "` + tempdir + `/png/sample1.gif"\nConverted: "` + tempdir + `/png/sample2.gif"\n`
		}},
		{decoder: gifDecoder, encoder: jpegEncoder, force: true, expected: func(tempdir string) string {
			return `Converted: "` + tempdir + `/gif/sample1.jpg"\n`
		}},
		{decoder: gifDecoder, encoder: pngEncoder, force: true, expected: func(tempdir string) string {
			return `Converted: "` + tempdir + `/gif/sample1.png"\n`
		}},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := &bytes.Buffer{}
			runner := Runner{OutStream: buf, Decoder: c.decoder, Encoder: c.encoder, Force: c.force}
			withTempDir(t, func(t *testing.T, tempdir string) {
				expectedStr := c.expected(tempdir)

				err := runner.Run(tempdir)
				if err != nil {
					t.FailNow()
				}

				if bufStr := buf.String(); !regexp.MustCompile(expectedStr).MatchString(bufStr) {
					t.Errorf("expected: %s, actual: %s", expectedStr, bufStr)
				}
			})
		})
	}
}

func TestCmd_Run_Nonexistence(t *testing.T) {
	t.Parallel()
	expected := "lstat nonexistent_path: no such file or directory"
	decoder := &conversion.Jpeg{}
	encoder := &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.NoCompression}}
	w := ioutil.Discard

	runner := Runner{OutStream: w, Decoder: decoder, Encoder: encoder, Force: true}
	err := runner.Run("nonexistent_path")
	actual := err.Error()
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestCmd_Run_Conflict(t *testing.T) {
	t.Parallel()
	decoder := &conversion.Jpeg{}
	encoder := &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.NoCompression}}
	w := ioutil.Discard

	withTempDir(t, func(t *testing.T, tempdir string) {
		expected := "File already exists: " + tempdir + "/jpeg/sample1.png"

		var runner Runner
		var err error

		runner = Runner{OutStream: w, Decoder: decoder, Encoder: encoder, Force: true}
		err = runner.Run(tempdir)
		if err != nil {
			t.Fatal(err)
		}

		runner = Runner{OutStream: w, Decoder: decoder, Encoder: encoder, Force: false}
		err = runner.Run(tempdir)

		actual := err.Error()
		if actual != expected {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	})
}

func withTempDir(t *testing.T, f func(t *testing.T, tempdir string)) {
	t.Helper()
	tempdir, _ := ioutil.TempDir("", "imgconv")
	fileutil.CopyDirRec("../testdata/", tempdir)
	defer os.RemoveAll(tempdir)
	f(t, tempdir)
}
