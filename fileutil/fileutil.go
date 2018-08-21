/*
Package fileutil is a collection of convenient functions for manipulating files
*/
package fileutil

import (
	"bytes"
	"io"
)

// StartsContentsWith returns whether file contents start with specified bytes.
func StartsContentsWith(rs io.ReadSeeker, xs []uint8) bool {
	buf := make([]byte, len(xs))
	rs.Seek(0, 0)
	rs.Read(buf)
	rs.Seek(0, 0)
	return bytes.Equal(buf, xs)
}
