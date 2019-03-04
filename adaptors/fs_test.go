package adaptors_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/int128/ghcp/adaptors"
)

func TestFileSystem_FindFiles(t *testing.T) {
	fs := &adaptors.FileSystem{}
	tempDir, err := ioutil.TempDir("", "FindFiles")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("dir1", 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("dir1/a.jpg", []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("dir2", 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("dir2/b.jpg", []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("dir2/c.jpg", []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("InDir", func(t *testing.T) {
		files, err := fs.FindFiles([]string{"."})
		if err != nil {
			t.Fatalf("FindFiles returned error: %+v", err)
		}
		wants := []string{"dir1/a.jpg", "dir2/b.jpg", "dir2/c.jpg"}
		if diff := deep.Equal(wants, files); diff != nil {
			t.Error(diff)
		}
	})
	t.Run("Files", func(t *testing.T) {
		files, err := fs.FindFiles([]string{"dir1/a.jpg", "dir2/c.jpg"})
		if err != nil {
			t.Fatalf("FindFiles returned error: %+v", err)
		}
		wants := []string{"dir1/a.jpg", "dir2/c.jpg"}
		if diff := deep.Equal(wants, files); diff != nil {
			t.Error(diff)
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		files, err := fs.FindFiles([]string{"dir3"})
		if files != nil {
			t.Errorf("files wants nil but %+v", files)
		}
		if err == nil {
			t.Fatalf("err wants non-nil but nil")
		}
	})
}

func TestFileSystem_ReadAsBase64EncodedContent(t *testing.T) {
	fs := &adaptors.FileSystem{}
	tempFile := makeTempFile(t, "hello\nworld")
	defer os.RemoveAll(tempFile)
	content, err := fs.ReadAsBase64EncodedContent(tempFile)
	if err != nil {
		t.Fatalf("ReadAsBase64EncodedContent returned error: %+v", err)
	}
	want := "aGVsbG8Kd29ybGQ="
	if want != content {
		t.Errorf("content wants %s but %s", want, content)
	}
}

func makeTempFile(t *testing.T, content string) string {
	tempFile, err := ioutil.TempFile("", "fs_test")
	if err != nil {
		t.Fatalf("error while creating a temp file: %s", err)
	}
	defer tempFile.Close()
	if _, err := fmt.Fprint(tempFile, content); err != nil {
		t.Fatalf("error while writing content: %s", err)
	}
	return tempFile.Name()
}
