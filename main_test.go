package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/powerman/check"
	_ "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) { check.TestMain(m) }

// - путь path/documents/cat
func Tests(t *testing.T) {
	t.Run("Abs2KindlePath", func(tt *testing.T) {
		t := check.T(tt)
		origKindleDir := kindleDir
		cases := []struct {
			path string
			want string
		}{
			{"/mnt/kindle/documents/Author/book.txt", "Author/book.txt"},
		}
		kindleDir = "/mnt/kindle"
		for _, v := range cases {
			got := Abs2KindlePath(v.path)
			t.Equal(got, v.want)
		}
		kindleDir = origKindleDir
	})
}

// - !IsDir()
func TestIsKindle0(tt *testing.T) {
	t := check.T(tt)
	tmpDir, err := ioutil.TempDir("", "TestIsKindle")
	t.Nil(err)
	fpath := tmpDir + "/documents"
	f, err := os.Create(fpath)
	t.Nil(err)
	defer f.Close()
	cases := []struct {
		dir  string
		want bool
	}{
		{fpath, false},
	}
	for _, v := range cases {
		got := isKindle(v.dir)
		t.Equal(got, v.want, kindleDir)
	}
	err = os.Remove(tmpDir + "/documents")
	t.Nil(err)
	err = os.Remove(tmpDir)
	t.Nil(err)
}

// - путь IsDir, содержит documents, system
func TestIsKindle1(tt *testing.T) {
	t := check.T(tt)
	tmpDir, err := ioutil.TempDir("", "TestIsKindle")
	t.Nil(err)
	path := tmpDir + "/documents"
	err = os.Mkdir(path, 0755)
	t.Nil(err)
	path = tmpDir + "/system"
	err = os.Mkdir(path, 0755)
	t.Nil(err)
	cases := []struct {
		dir  string
		want bool
	}{
		{tmpDir, true},
	}
	for _, v := range cases {
		got := isKindle(v.dir)
		t.Equal(got, v.want, kindleDir)
	}
	err = os.Remove(tmpDir + "/documents")
	t.Nil(err)
	err = os.Remove(tmpDir + "/system")
	t.Nil(err)
	err = os.Remove(tmpDir)
	t.Nil(err)
}

// - путь IsDir, содержит documents, не содержит system
func TestIsKindle2(tt *testing.T) {
	t := check.T(tt)
	tmpDir, err := ioutil.TempDir("", "TestIsKindle")
	t.Nil(err)
	path := tmpDir + "/documents"
	err = os.Mkdir(path, 0755)
	t.Nil(err)
	cases := []struct {
		dir  string
		want bool
	}{
		{tmpDir, false},
	}
	for _, v := range cases {
		got := isKindle(v.dir)
		t.Equal(got, v.want, kindleDir)
	}
	err = os.Remove(tmpDir + "/documents")
	t.Nil(err)
	err = os.Remove(tmpDir)
	t.Nil(err)
}

// - путь IsDir, содержит system, не содержит documents
func TestIsKindle3(tt *testing.T) {
	t := check.T(tt)
	tmpDir, err := ioutil.TempDir("", "TestIsKindle")
	t.Nil(err)
	path := tmpDir + "/system"
	err = os.Mkdir(path, 0755)
	t.Nil(err)
	cases := []struct {
		dir  string
		want bool
	}{
		{tmpDir, false},
	}
	for _, v := range cases {
		got := isKindle(v.dir)
		t.Equal(got, v.want, kindleDir)
	}
	err = os.Remove(tmpDir + "/system")
	t.Nil(err)
	err = os.Remove(tmpDir)
	t.Nil(err)
}

// - путь IsDir, не содержит documents, не содержит system
func TestIsKindle4(tt *testing.T) {
	t := check.T(tt)
	tmpDir, err := ioutil.TempDir("", "TestIsKindle")
	t.Nil(err)
	cases := []struct {
		dir  string
		want bool
	}{
		{tmpDir, false},
	}
	for _, v := range cases {
		got := isKindle(v.dir)
		t.Equal(got, v.want, kindleDir)
	}
	err = os.Remove(tmpDir)
	t.Nil(err)
}

// - путь IsDir, documents не каталог, system каталог
func TestIsKindle5(tt *testing.T) {
	t := check.T(tt)
	tmpDir, err := ioutil.TempDir("", "TestIsKindle")
	t.Nil(err)
	fpath := tmpDir + "/documents"
	f, err := os.Create(fpath)
	t.Nil(err)
	defer f.Close()
	path := tmpDir + "/system"
	err = os.Mkdir(path, 0755)
	t.Nil(err)
	cases := []struct {
		dir  string
		want bool
	}{
		{tmpDir, false},
	}
	for _, v := range cases {
		got := isKindle(v.dir)
		t.Equal(got, v.want, kindleDir)
	}
	err = os.Remove(tmpDir + "/documents")
	t.Nil(err)
	err = os.Remove(tmpDir + "/system")
	t.Nil(err)
	err = os.Remove(tmpDir)
	t.Nil(err)
}

// - путь IsDir, system не каталог, documents каталог
func TestIsKindle6(tt *testing.T) {
	t := check.T(tt)
	tmpDir, err := ioutil.TempDir("", "TestIsKindle")
	t.Nil(err)
	fpath := tmpDir + "/system"
	f, err := os.Create(fpath)
	t.Nil(err)
	defer f.Close()
	path := tmpDir + "/documents"
	err = os.Mkdir(path, 0755)
	t.Nil(err)
	cases := []struct {
		dir  string
		want bool
	}{
		{tmpDir, false},
	}
	for _, v := range cases {
		got := isKindle(v.dir)
		t.Equal(got, v.want, kindleDir)
	}
	err = os.Remove(tmpDir + "/system")
	t.Nil(err)
	err = os.Remove(tmpDir + "/documents")
	t.Nil(err)
	err = os.Remove(tmpDir)
	t.Nil(err)
}

// - пустая строка   = !nil
// - не файл/каталог = !nil
// - файл            = ErrNotADir
// - каталог         = nil
func TestDirExists(tt *testing.T) {
	t := check.T(tt)
	cases := []struct {
		dir     string
		wanterr bool
		errtext string
	}{
		{"", true, "no such file or directory"},
		{"/no/such/junk", true, "no such file or directory"},
		{"main.go", true, ErrNotADir.Error()},
		{"/", false, ""},
	}
	for _, v := range cases {
		err := dirExists(v.dir)
		if v.wanterr {
			t.NotNil(err)
			t.Match(err, v.errtext)
		} else {
			t.Nil(err)
		}
	}
}

func TestMatch(tt *testing.T) {
	t := check.T(tt)
	cases := []struct {
		want bool
		path string
	}{
		{false, "a.pdf"},
		{true, "a.azw"},
		{true, "a.azw123"},
	}
	for _, v := range cases {
		got := match(v.path)
		t.Equal(got, v.want, v.path)
	}
}
