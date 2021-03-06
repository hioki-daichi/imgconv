package fileutil

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileutil_StartsContentsWith(t *testing.T) {
	cases := map[string]struct {
		a        []byte
		b        []byte
		expected bool
	}{
		"\x01\x02 starts with \x01":                 {a: []byte("\x01\x02"), b: []byte("\x01"), expected: true},
		"\x01\x02 starts with \x01\x02":             {a: []byte("\x01\x02"), b: []byte("\x01\x02"), expected: true},
		"\x01\x02 does not start with \x01\x02\x03": {a: []byte("\x01\x02"), b: []byte("\x01\x02\x03"), expected: false},
		"\x01\x02 does not start with \x02":         {a: []byte("\x01\x02"), b: []byte("\x02"), expected: false},
	}

	for n, c := range cases {
		c := c
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			actual, err := StartsContentsWith(bytes.NewReader(c.a), c.b)
			if err != nil {
				t.Fatalf("err %s", err)
			}
			if actual != c.expected {
				t.Errorf(`expected="%t" actual="%t"`, c.expected, actual)
			}
		})
	}
}

func TestFileutil_StartsContentsWith_Unreadable(t *testing.T) {
	t.Parallel()

	expected := "EOF"

	fp, err := os.Open("./testdata/empty.txt")
	if err != nil {
		t.Fatalf("err %s", err)
	}

	_, err = StartsContentsWith(fp, []byte("\x01"))

	actual := err.Error()
	if actual != expected {
		t.Errorf(`expected="%s" actual="%s"`, expected, actual)
	}
}

func TestFileutil_StartsContentsWith_Unseekable(t *testing.T) {
	t.Parallel()

	expected := "unseekable"

	b := []byte("\x01")

	var err error
	var actual string

	// Seek() is called twice in fileutil.StartsContentsWith
	// In case of failure with the 1st Seek()
	_, err = StartsContentsWith(&readSeekerMock{seekableUntil: 0}, b)
	actual = err.Error()
	if actual != expected {
		t.Errorf(`expected="%s" actual="%s"`, expected, actual)
	}

	// In case of failure with the 2nd Seek()
	_, err = StartsContentsWith(&readSeekerMock{seekableUntil: 1}, b)
	actual = err.Error()
	if actual != expected {
		t.Errorf(`expected="%s" actual="%s"`, expected, actual)
	}
}

func TestFileutil_CopyDirRec(t *testing.T) {
	cases := map[string]struct {
		path string
	}{
		"./jpeg/sample1.jpg":  {path: "./jpeg/sample1.jpg"},
		"./jpeg/sample2.jpg":  {path: "./jpeg/sample2.jpg"},
		"./jpeg/sample3.jpeg": {path: "./jpeg/sample3.jpeg"},
		"./png/sample1.png":   {path: "./png/sample1.png"},
		"./png/sample2.png":   {path: "./png/sample2.png"},
		"./gif/sample1.gif":   {path: "./gif/sample1.gif"},
	}

	for n, c := range cases {
		c := c
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			tempdir, err := ioutil.TempDir("", "imgconv")
			if err != nil {
				t.Fatalf("err %s", err)
			}

			CopyDirRec("../testdata/", tempdir)
			defer os.RemoveAll(tempdir)

			_, err = os.OpenFile(filepath.Join(tempdir, c.path), os.O_CREATE|os.O_EXCL, 0)
			if !os.IsExist(err) {
				t.Fatalf("err %s", err)
			}
		})
	}
}

func TestFileutil_CopyDirRec_Nonexistence(t *testing.T) {
	t.Parallel()

	expected := "lstat ./nonexistent_src: no such file or directory"

	err := CopyDirRec("./nonexistent_src", "./nonexistent_dst")

	actual := err.Error()
	if actual != expected {
		t.Errorf(`expected="%s" actual="%s"`, expected, actual)
	}
}

func TestFileutil_CopyDirRec_Unopenable(t *testing.T) {
	t.Parallel()

	srcDir, err := ioutil.TempDir("", "imgconv")
	if err != nil {
		t.Fatalf("err %s", err)
	}

	srcPath := filepath.Join(srcDir, "unopenable.txt")

	expected := "open " + srcPath + ": permission denied"

	_, err = os.OpenFile(srcPath, os.O_CREATE, 000)
	if err != nil {
		t.Fatalf("err %s", err)
	}
	defer os.Remove(srcPath)

	dstDir, err := ioutil.TempDir("", "imgconv")
	if err != nil {
		t.Fatalf("err %s", err)
	}

	err = CopyDirRec(srcDir, dstDir)

	actual := err.Error()
	if actual != expected {
		t.Errorf(`expected="%s" actual="%s"`, expected, actual)
	}
}

func TestFileutil_CopyDirRec_MkdirFailure(t *testing.T) {
	t.Parallel()

	tempDir, err := ioutil.TempDir("", "imgconv")
	if err != nil {
		t.Fatalf("err %s", err)
	}

	dstPath := filepath.Join(tempDir, "foo")
	expected := "mkdir " + dstPath + "/gif: permission denied"

	err = os.Mkdir(dstPath, 0000)
	if err != nil {
		t.Fatalf("err %s", err)
	}
	defer os.Remove(dstPath)

	err = CopyDirRec("../testdata/", dstPath)

	actual := err.Error()
	if actual != expected {
		t.Errorf(`expected="%s" actual="%s"`, expected, actual)
	}
}

func TestFileutil_Copy_CreateError(t *testing.T) {
	expected := "error on create"
	tempDir, err := ioutil.TempDir("", "imgconv")
	if err != nil {
		t.Fatalf("err %s", err)
	}
	m := &createCopierMock{errOnCreate: true, errOnCopy: false}
	err = Copy(m, "../testdata/", tempDir)
	actual := err.Error()
	if actual != expected {
		t.Errorf(`expected="%s" actual="%s"`, expected, actual)
	}
}

func TestFileutil_Copy_CopyError(t *testing.T) {
	expected := "error on copy"
	tempDir, err := ioutil.TempDir("", "imgconv")
	if err != nil {
		t.Fatalf("err %s", err)
	}
	m := &createCopierMock{errOnCreate: false, errOnCopy: true}
	err = Copy(m, "../testdata/", tempDir)
	actual := err.Error()
	if actual != expected {
		t.Errorf(`expected="%s" actual="%s"`, expected, actual)
	}
}

type readSeekerMock struct {
	times int

	// Example:
	//
	//     m := &readSeekerMock{seekableUntil: 2}
	//     var err error
	//     _, err = m.Seek(0, 0); fmt.Println(err) // <nil>
	//     _, err = m.Seek(0, 0); fmt.Println(err) // <nil>
	//     _, err = m.Seek(0, 0); fmt.Println(err) // unseekable
	seekableUntil int
}

func (m *readSeekerMock) Read(_ []byte) (int, error) {
	return 0, nil
}

func (m *readSeekerMock) Seek(_ int64, _ int) (int64, error) {
	m.times++

	var err error
	if m.times > m.seekableUntil {
		err = errors.New("unseekable")
	}

	return 0, err
}

type createCopierMock struct {
	errOnCreate bool
	errOnCopy   bool
}

func (c *createCopierMock) Create(name string) (*os.File, error) {
	if c.errOnCreate {
		return nil, errors.New("error on create")
	}
	return nil, nil
}

func (c *createCopierMock) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	if c.errOnCopy {
		return 0, errors.New("error on copy")
	}
	return 0, nil
}
