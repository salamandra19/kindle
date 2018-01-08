package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/powerman/check"
	_ "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) { check.TestMain(m) }

// - путь path/documents/dir
func TestAbs2KindlePath(tt *testing.T) {
	t := check.T(tt)

	origKindleDir := kindleDir
	defer func() { kindleDir = origKindleDir }()
	kindleDir = "/mnt/kindle"

	t.Equal(Abs2KindlePath("/mnt/kindle/documents/Author/book.txt"), "Author/book.txt")
	t.PanicMatch(func() { Abs2KindlePath("book.txt") }, `not a kindle path`)
	t.PanicMatch(func() { Abs2KindlePath("") }, `not a kindle path`)
}

func TestIsKindle(tt *testing.T) {
	t := check.T(tt)

	tmpDir, err := ioutil.TempDir("", "gotest")
	t.Nil(err)
	defer func() { os.Remove(tmpDir) }()

	dpath := tmpDir + "/documents"
	spath := tmpDir + "/system"

	t.Run("Correct", func(tt *testing.T) {
		t := check.T(tt)

		t.Nil(os.Mkdir(dpath, 0755))
		t.Nil(os.Mkdir(spath, 0755))
		t.True(isKindle(tmpDir))
		t.Nil(os.Remove(dpath))
		t.Nil(os.Remove(spath))
	})

	t.Run("NotADir", func(tt *testing.T) {
		t := check.T(tt)

		path := tmpDir + "/kindle"
		t.Nil(ioutil.WriteFile(path, nil, 0644))
		t.False(isKindle(path))
		t.Nil(os.Remove(path))
	})

	t.Run("NoDocuments", func(tt *testing.T) {
		t := check.T(tt)

		t.Nil(os.Mkdir(spath, 0755))
		t.False(isKindle(tmpDir))
		t.Nil(os.Remove(spath))
	})

	t.Run("NoSystem", func(tt *testing.T) {
		t := check.T(tt)

		t.Nil(os.Mkdir(dpath, 0755))
		t.False(isKindle(tmpDir))
		t.Nil(os.Remove(dpath))
	})

	t.Run("Empty", func(tt *testing.T) {
		t := check.T(tt)

		t.False(isKindle(tmpDir))

	})

	t.Run("DocumentsNotADir", func(tt *testing.T) {
		t := check.T(tt)

		t.Nil(ioutil.WriteFile(dpath, nil, 0644))
		t.Nil(os.Mkdir(spath, 0755))
		t.False(isKindle(tmpDir))
		t.Nil(os.Remove(dpath))
		t.Nil(os.Remove(spath))

	})

	t.Run("SystemNotADir", func(tt *testing.T) {
		t := check.T(tt)

		t.Nil(os.Mkdir(dpath, 0755))
		t.Nil(ioutil.WriteFile(spath, nil, 0644))
		t.False(isKindle(tmpDir))
		t.Nil(os.Remove(dpath))
		t.Nil(os.Remove(spath))
	})
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
