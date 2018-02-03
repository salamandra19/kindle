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

// - передать путь к существующим файлам
//   - убедиться, что программа произвела корректные записи в collection
func addbooks(t *check.C, path string) int64 {
	fileInfo, err := os.Stat(path)
	t.Nil(err)
	lastAccess := fileInfo.ModTime().Unix() * 1000
	t.Nil(filePath(path, fileInfo, err))
	return lastAccess
}

func TestFilePath(tt *testing.T) {
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
		dpath + "/Author/file.txt",
		dpath + "/Author/new/file.pdf",
		dpath + "/Author/new/file.mobi",
		dpath + "/Author/new/file.prc",
		dpath + "/Author/new/file.azw",
		dpath + "/Author/new/file.azw3",
		dpath + "/Author/new/file.jpg",
		dpath + "/Author/new/file",
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

	var collectionTest = make(map[string]*Books)
	defer func() { collection = make(map[string]*Books) }()

	fileInfo, err := os.Stat(dpath + "/Author/file.txt")
	t.Nil(err)
	collectionTest["Author@en-US"] = &Books{
		Items:      []string{"*7c491ea1e58939598afbbcb48a5c5174d47b4650"},
		LastAccess: fileInfo.ModTime().Unix() * 1000,
	}
	t.Nil(filePath(dpath+"/Author/file.txt", fileInfo, err))
	t.DeepEqual(collection, collectionTest)

	fileInfo, err = os.Stat(dpath + "/Author/new/file.pdf")
	t.Nil(err)
	collectionTest["Author-new@en-US"] = &Books{
		Items:      []string{"*f827769a187f768d4145dfaceddfc4525b8e2e49"},
		LastAccess: fileInfo.ModTime().Unix() * 1000,
	}
	t.Nil(filePath(dpath+"/Author/new/file.pdf", fileInfo, err))
	t.DeepEqual(collection, collectionTest)

	books := collectionTest["Author-new@en-US"]

	lastAccess := addbooks(t, dpath+"/Author/new/file.mobi")
	books.Items = append(books.Items, "*6e57bc1c50ce96b679b7eb2b153f610e179d6609")
	books.LastAccess = lastAccess
	t.DeepEqual(collection, collectionTest)

	lastAccess = addbooks(t, dpath+"/Author/new/file.prc")
	books.Items = append(books.Items, "*8c2e2e36abee81f32ebf61d40f038a8c80b0893b")
	books.LastAccess = lastAccess
	t.DeepEqual(collection, collectionTest)

	lastAccess = addbooks(t, dpath+"/Author/new/file.azw")
	books.Items = append(books.Items, "*fecdd63969b387326772b2ff5586b83ede5682cf")
	books.LastAccess = lastAccess
	t.DeepEqual(collection, collectionTest)

	lastAccess = addbooks(t, dpath+"/Author/new/file.azw3")
	books.Items = append(books.Items, "*e465ef6ecde1994cdfc397df6841fd82baccc285")
	books.LastAccess = lastAccess
	t.DeepEqual(collection, collectionTest)

	// - передать путь с не соответствующим расширением
	//   - убедиться, что запись в collection не произведена
	fileInfo, err = os.Stat(dpath + "/Author/new/file.jpg")
	t.Nil(err)
	t.Nil(filePath(dpath+"/Author/new/file.jpg", fileInfo, err))
	t.DeepEqual(collection, collectionTest)

	// - передать путь без расширения
	//   - убедиться,что запись в collection не произведена
	fileInfo, err = os.Stat(dpath + "/Author/new/file")
	t.Nil(err)
	t.Nil(filePath(dpath+"/Author/new/file", fileInfo, err))
	t.DeepEqual(collection, collectionTest)

	// - передать путь к каталогу
	//   - убедиться, что запись в collection не произведена
	fileInfo, err = os.Stat(dpath + "/Author/new")
	t.Nil(err)
	t.Nil(filePath(dpath+"/Author/new", fileInfo, err))
	t.DeepEqual(collection, collectionTest)
}
func addBook(t *check.C, path string, tm time.Duration) (int64, func()) {
	t.Helper()
	t.Nil(ioutil.WriteFile(path, nil, 0644))
	t.Nil(os.Chtimes(path, time.Now().Add(tm), time.Now().Add(tm)))
	t.Nil(makeColl(path))
	fileInfo, err := os.Stat(path)
	t.Nil(err)
	lastAccess := fileInfo.ModTime().Unix() * 1000
	cleanup := func() { t.Nil(os.Remove(path)) }
	return lastAccess, cleanup
}

func TestMakeColl(tt *testing.T) {
	t := check.T(tt)

	// - создать временный каталог и в нем нужные для теста функции/файлы
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

	//   - убедиться что collection пустой
	t.Len(collection, 0)

	// - вызвать makeColl от одного из созданных файлов
	//   - убедиться, что collection изменился ожидаемым образом
	var collectionTest = make(map[string]*Books)

	// - очистить collection и временный каталог
	defer func() { collection = make(map[string]*Books) }()

	fileInfo, err := os.Stat(dpath + "/Author/new/books.txt")
	t.Nil(err)
	collectionTest["Author-new@en-US"] = &Books{
		Items:      []string{"*e7ecb5dd54c6a7bad47af33c936ed4c1d3deca01"},
		LastAccess: fileInfo.ModTime().Unix() * 1000,
	}
	t.Nil(makeColl(kindleDir + "/documents/Author/new/books.txt"))
	t.DeepEqual(collectionTest, collection)

	// - вызвать makeColl от следующего из созданных файлов
	//   - убедиться, что collection изменился ожидаемым образом
	fileInfo, err = os.Stat(dpath + "/Author/new/workpad.txt")
	t.Nil(err)
	collectionTest["Author-new@en-US"].Items = append(collectionTest["Author-new@en-US"].Items, "*1c80d93e03067312fd43b17d2a339d375e4bd560")
	collectionTest["Author-new@en-US"].LastAccess = fileInfo.ModTime().Unix() * 1000
	t.Nil(makeColl(kindleDir + "/documents/Author/new/workpad.txt"))
	t.DeepEqual(collection, collectionTest)

	// - записать значение в LastAccess
	// - передать путь к более свежему файлу
	//   - убедиться, что в LastAccess записаны данные модификации самые свежие
	books := collectionTest["Author-new@en-US"]

	lastAccess, cleanup := addBook(t, dpath+"/Author/new/timePast.txt", -10*time.Second)
	defer cleanup()
	books.Items = append(books.Items, "*86062d936f58942786e6e9dfb8443ddfedfe0e28")
	t.DeepEqual(collection, collectionTest)

	lastAccess, cleanup = addBook(t, dpath+"/Author/new/timeNow.txt", 1*time.Second)
	defer cleanup()
	books.Items = append(books.Items, "*14a25ed50202fbc2c93fe3acb2dd72386212302c")
	books.LastAccess = lastAccess
	t.DeepEqual(collection, collectionTest)

	lastAccess, cleanup = addBook(t, dpath+"/Author/new/timeFuture.txt", 10*time.Second)
	defer cleanup()
	books.Items = append(books.Items, "*b3916c8105ceebfa611abf38ce331c4873046de7")
	books.LastAccess = lastAccess
	t.DeepEqual(collection, collectionTest)

	// - убедиться, что в collection добавляются новые записи
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

	// - вызвать makeColl от каталога, вместо файла
	//  - убедиться, что программа выдаёт ошибку или панику
	t.PanicMatch(func() { makeColl(kindleDir + "/documents/Author/new") }, `not a path to file`)

	// - передать makeColl не корректный путь
	//  - убедиться, что программа выдаёт ошибку или панику
	t.PanicMatch(func() { makeColl("/catalog/name.txt") }, `not a kindle path`)
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
