package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

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

// - создать временный каталог и в нем нужные для теста функции/файлы
// - убедиться что collection пустой
// - вызвать makeColl от одного из созданных файлов
// - убедиться, что collection изменился ожидаемым образом
// - вызвать makeColl от следующего из созданных файлов
// - убедиться, что collection изменился ожидаемым образом
// - вызвать makeColl от каталога, вместо файла
// - убедиться, что программа выдаёт ошибку или панику
// - передать makeColl не корректный путь
// - убедиться, что программа выдаёт ошибку или панику
// - записать значение в LastAccess
// - передать путь к более свежему файлу
// - убедиться, что в LastAccess записаны данные модификации самые свежие
// - убедиться, что в collection добавляются новые записи
// - очистить collection и временный каталог

func TestMakeColl(tt *testing.T) {
	t := check.T(tt)

	tmpDir, err := ioutil.TempDir("", "gotest")
	t.Nil(err)
	defer func() { t.Nil(os.Remove(tmpDir)) }()

	origKindleDir := kindleDir
	defer func() { kindleDir = origKindleDir }()
	kindleDir = tmpDir

	dpath := tmpDir + "/documents"
	dirs := []string{
		tmpDir + "/system",
		dpath,
		dpath + "/Author",
		dpath + "/Author/new",
	}
	files := []string{
		dpath + "/Author/new/books.txt",
		dpath + "/Author/new/workpad.txt",
	}

	for _, d := range dirs {
		t.Nil(os.Mkdir(d, 0755))
	}
	defer func() {
		for i := len(dirs) - 1; i >= 0; i-- {
			t.Nil(os.Remove(dirs[i]))
		}
	}()

	for _, f := range files {
		t.Nil(ioutil.WriteFile(f, nil, 0644))
	}
	defer func() {
		for i := len(files) - 1; i >= 0; i-- {
			t.Nil(os.Remove(files[i]))
		}
	}()

	t.Len(collection, 0)

	var collectionTest = make(map[string]*Books)
	defer func() { collection = make(map[string]*Books) }()

	fileInfo, err := os.Stat(dpath + "/Author/new/books.txt")
	t.Nil(err)
	collectionTest["Author-new@en-US"] = &Books{
		Items:      []string{"*e7ecb5dd54c6a7bad47af33c936ed4c1d3deca01"},
		LastAccess: fileInfo.ModTime().Unix() * 1000,
	}
	t.Nil(makeColl(kindleDir + "/documents/Author/new/books.txt"))
	t.DeepEqual(collectionTest, collection)

	fileInfo, err = os.Stat(dpath + "/Author/new/workpad.txt")
	t.Nil(err)
	collectionTest["Author-new@en-US"].Items = append(collectionTest["Author-new@en-US"].Items, "*1c80d93e03067312fd43b17d2a339d375e4bd560")
	collectionTest["Author-new@en-US"].LastAccess = fileInfo.ModTime().Unix() * 1000
	t.Nil(makeColl(kindleDir + "/documents/Author/new/workpad.txt"))
	t.DeepEqual(collection, collectionTest)

	pathFirst := dpath + "/Author/new/timePast.txt"
	timeFirst := time.Now().Add(-10 * time.Second)
	t.Nil(ioutil.WriteFile(pathFirst, nil, 0644))
	t.Nil(os.Chtimes(pathFirst, timeFirst, timeFirst))
	defer func() { os.Remove(pathFirst) }()
	collectionTest["Author-new@en-US"].Items = append(collectionTest["Author-new@en-US"].Items, "*86062d936f58942786e6e9dfb8443ddfedfe0e28")
	t.Log(collectionTest["Author-new@en-US"].LastAccess)
	t.Nil(makeColl(pathFirst))
	t.DeepEqual(collection, collectionTest)

	pathSecond := dpath + "/Author/new/timeNow.txt"
	timeSecond := time.Now()
	t.Nil(ioutil.WriteFile(pathSecond, nil, 0644))
	t.Nil(os.Chtimes(pathSecond, timeSecond, timeSecond))
	defer func() { os.Remove(pathSecond) }()
	fileInfo, err = os.Stat(pathSecond)
	t.Nil(err)
	collectionTest["Author-new@en-US"].Items = append(collectionTest["Author-new@en-US"].Items, "*14a25ed50202fbc2c93fe3acb2dd72386212302c")
	collectionTest["Author-new@en-US"].LastAccess = fileInfo.ModTime().Unix() * 1000
	t.Log(collectionTest["Author-new@en-US"].LastAccess)
	t.Nil(makeColl(pathSecond))
	t.DeepEqual(collection, collectionTest)

	pathThird := dpath + "/Author/new/timeFuture.txt"
	timeThird := time.Now().Add(10 * time.Second)
	t.Nil(ioutil.WriteFile(pathThird, nil, 0644))
	t.Nil(os.Chtimes(pathThird, timeThird, timeThird))
	defer func() { os.Remove(pathThird) }()
	fileInfo, err = os.Stat(pathThird)
	t.Nil(err)
	collectionTest["Author-new@en-US"].Items = append(collectionTest["Author-new@en-US"].Items, "*b3916c8105ceebfa611abf38ce331c4873046de7")
	collectionTest["Author-new@en-US"].LastAccess = fileInfo.ModTime().Unix() * 1000
	t.Log(collectionTest["Author-new@en-US"].LastAccess)
	t.Nil(makeColl(pathThird))
	t.DeepEqual(collection, collectionTest)

	pathNewDir := dpath + "/Author/morenew"
	t.Nil(os.Mkdir(pathNewDir, 0755))
	defer func() { os.Remove(pathNewDir) }()
	pathNewKey := dpath + "/Author/morenew/NewKey.txt"
	t.Nil(ioutil.WriteFile(pathNewKey, nil, 0644))
	defer func() { os.Remove(pathNewKey) }()
	fileInfo, err = os.Stat(pathNewKey)
	t.Nil(err)
	collectionTest["Author-morenew@en-US"] = &Books{
		Items:      []string{"*5a46239d8e667091663fd4792e74bc6bf0eddaea"},
		LastAccess: fileInfo.ModTime().Unix() * 1000,
	}
	t.Nil(makeColl(pathNewKey))
	t.DeepEqual(collection, collectionTest)
}

func TestIsKindle(tt *testing.T) {
	t := check.T(tt)

	tmpDir, err := ioutil.TempDir("", "gotest")
	t.Nil(err)
	defer func() { t.Nil(os.Remove(tmpDir)) }()

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
