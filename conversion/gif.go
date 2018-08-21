package conversion

import (
	"image"
	"image/gif"
	"io"
	"os"
	"path/filepath"

	"github.com/hioki-daichi/imgconv/fileutil"
)

// Gif https://en.wikipedia.org/wiki/GIF
type Gif struct {
	Options *gif.Options
}

// Encode encodes the specified file to GIF
func (g *Gif) Encode(w io.Writer, img image.Image) error {
	return gif.Encode(w, img, g.Options)
}

// Decode decodes the specified GIF file
func (g *Gif) Decode(r io.Reader) (image.Image, error) {
	return gif.Decode(r)
}

// Extname returns "gif"
func (g *Gif) Extname() string {
	return "gif"
}

// IsDecodable returns whether the file content is GIF
func (g *Gif) IsDecodable(fp *os.File) bool {
	return fileutil.StartsContentsWith(fp, []uint8{71, 73, 70, 56, 55, 97}) || fileutil.StartsContentsWith(fp, []uint8{71, 73, 70, 56, 57, 97})
}

// HasProcessableExtname returns whether the specified path has ".gif"
func (g *Gif) HasProcessableExtname(path string) bool {
	return filepath.Ext(path) == ".gif"
}
