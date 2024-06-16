package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type singleNameFilter struct {
	t    *testing.T
	dir  string // name of directory to skip (if empty, do nothing)
	file string // name of file to exclude (if empty, do nothing)
}

func (f *singleNameFilter) SkipDir(path string) bool {
	f.t.Logf("visiting dir %s", path)
	base := filepath.Base(path)
	return f.dir == base
}

func (f *singleNameFilter) ExcludeFile(path string) bool {
	f.t.Logf("visiting file %s", path)
	base := filepath.Base(path)
	return f.file == base
}

func TestFileSystem_FindFiles(t *testing.T) {
	fs := &FileSystem{}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("dir1", 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("dir1/a.jpg", []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("dir2", 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("dir2/b.jpg", []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("dir2/c.jpg", []byte{}, 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("FindDirectory", func(t *testing.T) {
		got, err := fs.FindFiles([]string{"."}, &singleNameFilter{t: t})
		if err != nil {
			t.Fatalf("FindFiles returned error: %+v", err)
		}
		want := []File{
			{Path: "dir1/a.jpg"},
			{Path: "dir2/b.jpg"},
			{Path: "dir2/c.jpg", Executable: true},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("FilterIsNil", func(t *testing.T) {
		got, err := fs.FindFiles([]string{"."}, nil)
		if err != nil {
			t.Fatalf("FindFiles returned error: %+v", err)
		}
		want := []File{
			{Path: "dir1/a.jpg"},
			{Path: "dir2/b.jpg"},
			{Path: "dir2/c.jpg", Executable: true},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("FindFiles", func(t *testing.T) {
		got, err := fs.FindFiles([]string{"dir1/a.jpg", "dir2/c.jpg"}, &singleNameFilter{t: t})
		if err != nil {
			t.Fatalf("FindFiles returned error: %+v", err)
		}
		want := []File{
			{Path: "dir1/a.jpg"},
			{Path: "dir2/c.jpg", Executable: true},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("NoSuchFile", func(t *testing.T) {
		files, err := fs.FindFiles([]string{"dir3"}, &singleNameFilter{t: t})
		if files != nil {
			t.Errorf("files wants nil but %+v", files)
		}
		if err == nil {
			t.Fatalf("err wants non-nil but nil")
		}
	})
	t.Run("ExcludeDirectory", func(t *testing.T) {
		got, err := fs.FindFiles([]string{"."}, &singleNameFilter{t: t, dir: "dir2"})
		if err != nil {
			t.Fatalf("FindFiles returned error: %+v", err)
		}
		want := []File{
			{Path: "dir1/a.jpg"},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("SkipFile", func(t *testing.T) {
		got, err := fs.FindFiles([]string{"."}, &singleNameFilter{t: t, file: "b.jpg"})
		if err != nil {
			t.Fatalf("FindFiles returned error: %+v", err)
		}
		want := []File{
			{Path: "dir1/a.jpg"},
			{Path: "dir2/c.jpg", Executable: true},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestFileSystem_ReadAsBase64EncodedContent(t *testing.T) {
	fs := &FileSystem{}
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "fs_test")
	if err := os.WriteFile(tempFile, []byte("hello\nworld"), 0644); err != nil {
		t.Fatal(err)
	}
	content, err := fs.ReadAsBase64EncodedContent(tempFile)
	if err != nil {
		t.Fatalf("ReadAsBase64EncodedContent returned error: %+v", err)
	}
	want := "aGVsbG8Kd29ybGQ="
	if want != content {
		t.Errorf("content wants %s but %s", want, content)
	}
}
