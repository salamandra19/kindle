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
func TestAbs2KindlePath(tt *testing.T) {
	t := check.T(tt)

	origKindleDir := kindleDir
	cases := []struct {
		path string
		want string
	}{
		{"/mnt/kindle/documents/Author/book.txt", "Author/book.txt"},
		//{"book.txt", "book.txt"},
		//{"", ""},
	}
	kindleDir = "/mnt/kindle"
	for _, v := range cases {
		got := Abs2KindlePath(v.path)
		t.Equal(got, v.want)
	}
	kindleDir = origKindleDir
}

func TestS(t *testing.T) {
	t.Run("TestIsKindle_NotADir", func(tt *testing.T) {
		t := check.T(tt)

		tmpDir, err := ioutil.TempDir("", "gotest")
		t.Nil(err)
		defer func() { t.Nil(os.Remove(tmpDir)) }()

		path := tmpDir + "/kindle"
		t.Nil(ioutil.WriteFile(path, nil, 0644))
		defer func() { t.Nil(os.Remove(path)) }()

		t.False(isKindle(path))
	})

	t.Run("TestIsKindle", func(tt *testing.T) {
		t := check.T(tt)

		tmpDir, err := ioutil.TempDir("", "TestIsKindle")
		t.Nil(err)
		defer func() { os.Remove(tmpDir) }()

		dpath := tmpDir + "/documents"
		t.Nil(os.Mkdir(dpath, 0755))
		defer func() { os.Remove(dpath) }()

		spath := tmpDir + "/system"
		t.Nil(os.Mkdir(spath, 0755))
		defer func() { os.Remove(spath) }()

		t.True(isKindle(tmpDir))
	})

	t.Run("TestIsKindle_NoSystem", func(tt *testing.T) {
		t := check.T(tt)

		tmpDir, err := ioutil.TempDir("", "TestIsKindle")
		t.Nil(err)
		defer func() { os.Remove(tmpDir) }()

		path := tmpDir + "/system"
		t.Nil(os.Mkdir(path, 0755))
		defer func() { os.Remove(path) }()

		t.False(isKindle(tmpDir))
		t.Nil(err)
	})

	t.Run("TestIsKindle_NoDocuments", func(tt *testing.T) {
		t := check.T(tt)

		tmpDir, err := ioutil.TempDir("", "TestIsKindle")
		t.Nil(err)
		defer func() { os.Remove(tmpDir) }()

		path := tmpDir + "/system"
		t.Nil(os.Mkdir(path, 0755))
		defer func() { os.Remove(path) }()

		t.False(isKindle(tmpDir))
	})

	t.Run("TestIsKindle_Empty", func(tt *testing.T) {
		t := check.T(tt)

		tmpDir, err := ioutil.TempDir("", "TestIsKindle")
		t.Nil(err)
		defer func() { t.Nil(os.Remove(tmpDir)) }()

		t.False(isKindle(tmpDir))

	})

	t.Run("TestIsKindle_DocumentsNotADir", func(tt *testing.T) {
		t := check.T(tt)

		tmpDir, err := ioutil.TempDir("", "TestIsKindle")
		t.Nil(err)
		defer func() { t.Nil(os.Remove(tmpDir)) }()

		fpath := tmpDir + "/documents"
		t.Nil(ioutil.WriteFile(fpath, nil, 0644))
		defer func() { t.Nil(os.Remove(fpath)) }()

		path := tmpDir + "/system"
		t.Nil(os.Mkdir(path, 0755))
		defer func() { t.Nil(os.Remove(path)) }()

		t.False(isKindle(tmpDir))
	})

	t.Run("TestIsKindle_SystemNotADir", func(tt *testing.T) {
		t := check.T(tt)

		tmpDir, err := ioutil.TempDir("", "TestIsKindle")
		t.Nil(err)
		defer func() { t.Nil(os.Remove(tmpDir)) }()

		fpath := tmpDir + "/system"
		t.Nil(ioutil.WriteFile(fpath, nil, 0644))
		defer func() { t.Nil(os.Remove(fpath)) }()

		path := tmpDir + "/documents"
		t.Nil(os.Mkdir(path, 0755))
		defer func() { t.Nil(os.Remove(path)) }()

		t.False(isKindle(tmpDir))
	})

	// - пустая строка   = !nil
	// - не файл/каталог = !nil
	// - файл            = ErrNotADir
	// - каталог         = nil
	t.Run("TestDirExists", func(tt *testing.T) {
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
	})

	t.Run("TestMatch", func(tt *testing.T) {
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
	})
}
