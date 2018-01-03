package main

import (
	"testing"

	"github.com/powerman/check"
	_ "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) { check.TestMain(m) }

// - путь path/documents/cat
// - путь path/documents

func TestAbs2KindlePath(tt *testing.T) {
	t := check.T(tt)
	cases := []struct {
		path string
		want string
	}{
		{"/home/tanya/documents/Author/book.txt", "Author/book.txt"},
		{"/home/tanya/documents", ""},
	}
	for _, v := range cases {
		got := Abs2KindlePath(v.path)
		t.Equal(got, v.want)
	}
}

/*
func TestIsKindle(tt *testing.T) {
	t := check.T(tt)
// - !IsDir()
// - путь IsDir, содержит documents, system
// - путь IsDir, содержит documents, не содержит system
// - путь IsDir, содержит system, не содержит documents
// - путь IsDir, не содержит documents, не содержит system
// - путь !IsDir, содержит documents, system
// - путь !IsDir, содержит documents, не содержит system
// - путь !IsDir, содержит system, не содержит documents
	cases := []struct{
		dir string
		want bool
	}{
		{"/gocode/src/library",true},
		{"",true},
		{"",true},
		{"",true},
		{"",true},
	}

	}
}
*/
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
