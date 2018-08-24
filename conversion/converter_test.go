package conversion

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hioki-daichi/imgconv/fileutil"
)

func TestConversion_Convert(t *testing.T) {
	jpegDecoder := &Jpeg{}
	pngDecoder := &Png{}
	gifDecoder := &Gif{}

	jpegEncoder := &Jpeg{Options: &jpeg.Options{Quality: 1}}
	pngEncoder := &Png{Encoder: &png.Encoder{CompressionLevel: png.NoCompression}}
	gifEncoder := &Gif{Options: &gif.Options{NumColors: 1}}

	cases := []struct {
		decoder  Decoder
		encoder  Encoder
		path     string
		force    bool
		expected error
	}{
		// JPEG to PNG
		{decoder: jpegDecoder, encoder: pngEncoder, path: "./jpeg/sample1.jpg", force: true, expected: nil},
		// JPEG to GIF
		{decoder: jpegDecoder, encoder: gifEncoder, path: "./jpeg/sample1.jpg", force: true, expected: nil},

		// PNG to JPEG
		{decoder: pngDecoder, encoder: jpegEncoder, path: "./png/sample1.png", force: true, expected: nil},
		// PNG to GIF
		{decoder: pngDecoder, encoder: gifEncoder, path: "./png/sample1.png", force: true, expected: nil},

		// GIF to JPEG
		{decoder: gifDecoder, encoder: jpegEncoder, path: "./gif/sample1.gif", force: true, expected: nil},
		// GIF to PNG
		{decoder: gifDecoder, encoder: pngEncoder, path: "./gif/sample1.gif", force: true, expected: nil},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			converter := &Converter{Decoder: c.decoder, Encoder: c.encoder}

			withTempDir(t, func(t *testing.T, tempdir string) {
				if _, actual := converter.Convert(filepath.Join(tempdir, c.path), c.force); c.expected != actual {
					t.Errorf("expected: %s, actual: %s", c.expected, actual)
				}
			})
		})
	}
}

func TestConversion_Convert_Conflict(t *testing.T) {
	converter := &Converter{Decoder: &Jpeg{}, Encoder: &Png{Encoder: &png.Encoder{CompressionLevel: png.NoCompression}}}

	withTempDir(t, func(t *testing.T, tempdir string) {
		expected := "File already exists: " + tempdir + "/jpeg/sample1.png"

		path := filepath.Join(tempdir, "./jpeg/sample1.jpg")

		_, _ = converter.Convert(path, false)
		_, err := converter.Convert(path, false)

		actual := err.Error()

		if actual != expected {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	})
}

func TestConversion_Convert_Nonexistence(t *testing.T) {
	expected := "open ./nonexistent_path: no such file or directory"
	converter := &Converter{Decoder: &Jpeg{}, Encoder: &Png{Encoder: &png.Encoder{CompressionLevel: png.NoCompression}}}
	_, err := converter.Convert("./nonexistent_path", true)
	actual := err.Error()
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func withTempDir(t *testing.T, f func(t *testing.T, tempdir string)) {
	tempdir, _ := ioutil.TempDir("", "imgconv")
	fileutil.CopyDirRec("../testdata/", tempdir)
	defer os.RemoveAll(tempdir)
	f(t, tempdir)
}
