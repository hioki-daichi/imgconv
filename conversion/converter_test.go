package conversion

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func withTempDir(t *testing.T, f func(t *testing.T, tempdir string)) {
	tempdir, _ := ioutil.TempDir("", "imgconv")
	copyDirRec("../testdata/", tempdir)
	defer os.RemoveAll(tempdir)
	f(t, tempdir)
}

func copyDirRec(src string, dest string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		sf, err := os.Open(path)
		if err != nil {
			return err
		}

		destDir := filepath.Join(dest, strings.TrimLeft(filepath.Dir(path), src))

		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, filepath.Base(path))

		df, err := os.Create(destPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(df, sf)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
