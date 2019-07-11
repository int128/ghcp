package env

import (
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/int128/ghcp/adaptors"
	"golang.org/x/xerrors"
)

// FileSystem provides manipulation of file system.
type FileSystem struct{}

// FindFiles returns a list of files in the paths.
func (fs *FileSystem) FindFiles(paths []string) ([]adaptors.File, error) {
	var files []adaptors.File
	for _, path := range paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return xerrors.Errorf("error while walk: %w", err)
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			files = append(files, adaptors.File{
				Path:       path,
				Executable: info.Mode()&0100 != 0, // mask the executable bit of owner
			})
			return nil
		}); err != nil {
			return nil, xerrors.Errorf("error while finding files in %s: %w", path, err)
		}
	}
	return files, nil
}

// ReadAsBase64EncodedContent returns content of the file as base64 encoded string.
func (fs *FileSystem) ReadAsBase64EncodedContent(filename string) (string, error) {
	r, err := os.Open(filename)
	if err != nil {
		return "", xerrors.Errorf("error while opening file %s: %w", filename, err)
	}
	defer r.Close()
	var s strings.Builder
	e := base64.NewEncoder(base64.StdEncoding, &s)
	if _, err := io.Copy(e, r); err != nil {
		return "", xerrors.Errorf("error while encoding file %s: %w", filename, err)
	}
	e.Close()
	return s.String(), nil
}
