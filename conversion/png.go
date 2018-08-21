package conversion

import (
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/hioki-daichi/imgconv/fileutil"
)

// Png https://en.wikipedia.org/wiki/Portable_Network_Graphics
type Png struct {
	Encoder *png.Encoder
}

// Encode encodes the specified file to PNG
func (p *Png) Encode(w io.Writer, img image.Image) error {
	return p.Encoder.Encode(w, img)
}

// Decode decodes the specified PNG file
func (p *Png) Decode(r io.Reader) (image.Image, error) {
	return png.Decode(r)
}

// Extname returns "png"
func (p *Png) Extname() string {
	return "png"
}

// IsDecodable returns whether the file content is PNG
func (p *Png) IsDecodable(fp *os.File) bool {
	return fileutil.StartsContentsWith(fp, []uint8{137, 80, 78, 71, 13, 10, 26, 10})
}

// HasProcessableExtname returns whether the specified path has ".png"
func (p *Png) HasProcessableExtname(path string) bool {
	return filepath.Ext(path) == ".png"
}
