package opt

import (
	"errors"
	"image/gif"
	"image/jpeg"
	"image/png"
	"reflect"
	"testing"

	"github.com/hioki-daichi/imgconv/conversion"
)

func TestOpt_Parse(t *testing.T) {
	cases := []struct {
		args    []string
		dirname string
		options *Options
		err     error
	}{
		// no argument
		{args: []string{}, dirname: "", options: nil, err: errors.New("you must specify a directory")},

		// dirname noly
		{args: []string{"./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.DefaultCompression}}, Force: false}, err: nil},

		// force
		{args: []string{"-f", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.DefaultCompression}}, Force: true}, err: nil},

		// in/out format specified
		{args: []string{"-J", "-p", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.DefaultCompression}}, Force: false}, err: nil},
		{args: []string{"-J", "-p", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.DefaultCompression}}, Force: false}, err: nil},
		{args: []string{"-J", "-g", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Gif{Options: &gif.Options{NumColors: 256}}, Force: false}, err: nil},
		{args: []string{"-P", "-j", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Png{}, Encoder: &conversion.Jpeg{Options: &jpeg.Options{Quality: 100}}, Force: false}, err: nil},
		{args: []string{"-P", "-g", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Png{}, Encoder: &conversion.Gif{Options: &gif.Options{NumColors: 256}}, Force: false}, err: nil},
		{args: []string{"-G", "-j", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Gif{}, Encoder: &conversion.Jpeg{Options: &jpeg.Options{Quality: 100}}, Force: false}, err: nil},
		{args: []string{"-G", "-p", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Gif{}, Encoder: &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.DefaultCompression}}, Force: false}, err: nil},

		// quality
		{args: []string{"-P", "-j", "--quality=0", "./testdata/"}, dirname: "", options: nil, err: errors.New("--quality must be greater than or equal to 1")},
		{args: []string{"-P", "-j", "--quality=1", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Png{}, Encoder: &conversion.Jpeg{Options: &jpeg.Options{Quality: 1}}, Force: false}, err: nil},
		{args: []string{"-P", "-j", "--quality=100", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Png{}, Encoder: &conversion.Jpeg{Options: &jpeg.Options{Quality: 100}}, Force: false}, err: nil},
		{args: []string{"-P", "-j", "--quality=101", "./testdata/"}, dirname: "", options: nil, err: errors.New("--quality must be less than or equal to 100")},

		// num-colors
		{args: []string{"-J", "-g", "--num-colors=0", "./testdata/"}, dirname: "", options: nil, err: errors.New("--num-colors must be greater than or equal to 1")},
		{args: []string{"-J", "-g", "--num-colors=1", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Gif{Options: &gif.Options{NumColors: 1}}, Force: false}, err: nil},
		{args: []string{"-J", "-g", "--num-colors=256", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Gif{Options: &gif.Options{NumColors: 256}}, Force: false}, err: nil},
		{args: []string{"-J", "-g", "--num-colors=257", "./testdata/"}, dirname: "", options: nil, err: errors.New("--num-colors must be less than or equal to 256")},

		// compression-level
		{args: []string{"-J", "-p", "--compression-level=default", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.DefaultCompression}}, Force: false}, err: nil},
		{args: []string{"-J", "-p", "--compression-level=no", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.NoCompression}}, Force: false}, err: nil},
		{args: []string{"-J", "-p", "--compression-level=best-speed", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.BestSpeed}}, Force: false}, err: nil},
		{args: []string{"-J", "-p", "--compression-level=best-compression", "./testdata/"}, dirname: "./testdata/", options: &Options{Decoder: &conversion.Jpeg{}, Encoder: &conversion.Png{Encoder: &png.Encoder{CompressionLevel: png.BestCompression}}, Force: false}, err: nil},
		{args: []string{"-J", "-p", "--compression-level=foo", "./testdata/"}, dirname: "", options: nil, err: errors.New("--compression-level is not included in the list: \"default\", \"no\", \"best-speed\", \"best-compression\"")},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			dirname, options, err := Parse(c.args...)
			if err != c.err && err.Error() != c.err.Error() {
				t.Errorf("expected: %s, actual: %s", c.err.Error(), err.Error())
			}

			if dirname != c.dirname {
				t.FailNow()
			}

			if options != c.options {
				if !reflect.DeepEqual(options.Decoder, c.options.Decoder) {
					t.FailNow()
				}

				if !reflect.DeepEqual(options.Encoder, c.options.Encoder) {
					t.FailNow()
				}

				if options.Force != c.options.Force {
					t.FailNow()
				}
			}
		})
	}
}
