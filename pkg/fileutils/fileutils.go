package fileutils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tomogoma/go-typed-errors"
)

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, errors.Newf("stat: %v", err)
}

// CopyIfDestNotExists copies the from file into dest file if dest does not exists.
// see CopyFile for Notes.
func CopyIfDestNotExists(from, dest string) error {
	_, err := os.Stat(dest)
	if err == nil {
		fmt.Printf("'%s' ignored, already exists\n", dest)
		return nil
	}
	if !os.IsNotExist(err) {
		return errors.Newf("stat: %v", err)
	}
	return CopyFile(from, dest)
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return errors.Newf("open src: %v", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return errors.Newf("create dst file: %v", err)
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = errors.Newf("%v ...close dst file: %v", err, e)
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return errors.Newf("copy contents: %v", err)
	}

	err = out.Sync()
	if err != nil {
		return errors.Newf("flush buffer after copy: %v", err)
	}

	si, err := os.Stat(src)
	if err != nil {
		return errors.Newf("stat src file: %v", err)
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return errors.Newf("chmod dest to equal src perms: %v", err)
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist.
// Any source content that already exists in destination will be ignored and skipped.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return errors.Newf("stat source: %v", err)
	}
	if !si.IsDir() {
		return errors.Newf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return errors.Newf("stat destination: %v", err)
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return errors.Newf("mkdirall destination: %v", err)
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.Newf("read source: %v", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return errors.Newf("copy child dir: %v", err)
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyIfDestNotExists(srcPath, dstPath)
			if err != nil {
				return errors.Newf("copy child file: %v", err)
			}
		}
	}

	return
}
