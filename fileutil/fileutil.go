/*
Package fileutil is a collection of convenient functions for manipulating files
*/
package fileutil

import (
	"bytes"
	"io"
)

// StartsContentsWith returns whether file contents start with specified bytes.
func StartsContentsWith(rs io.ReadSeeker, xs []byte) (bool, error) {
	buf := make([]byte, len(xs))

	_, err := rs.Seek(0, 0)
	if err != nil {
		return false, err
	}

	_, err = rs.Read(buf)
	if err != nil {
		return false, err
	}

	_, err = rs.Seek(0, 0)
	if err != nil {
		return false, err
	}

	return bytes.Equal(buf, xs), nil
}
