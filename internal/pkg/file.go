package pkg

import (
	"archive/tar"
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v6"
	"github.com/moby/go-archive"
	"github.com/moby/patternmatcher"
)

// WriteFile writes data to the specified path, creating parent directories if necessary.
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func WriteFileToBillyFs(fs billy.Filesystem, path string, content []byte) (err error) {
	f, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = f.Write(content)
	return err
}

func ReadDockerignore(fs billy.Filesystem) (excludes []string, err error) {
	f, err := fs.Open(".dockerignore")
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		excludes = append(excludes, line)
	}

	return excludes, scanner.Err()
}

func DeleteDir(path string) error {
	return os.RemoveAll(path)
}

func TarDirectory(dir string, excludes []string) (io.ReadCloser, error) {
	return archive.TarWithOptions(dir, &archive.TarOptions{
		ExcludePatterns: excludes,
	})
}

func TarFilesystem(
	fs billy.Filesystem,
	excludes []string,
) (io.ReadCloser, error) {
	pr, pw := io.Pipe()

	pm, err := patternmatcher.New(excludes)
	if err != nil {
		return nil, err
	}

	go func() {
		tw := tar.NewWriter(pw)

		err := walkFilesystem(fs, ".", func(path string) error {
			if shouldExclude(path, pm) {
				return nil
			}

			info, err := fs.Stat(path)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := fs.Open(path)
			if err != nil {
				return err
			}

			hdr, err := tar.FileInfoHeader(info, "")
			if err != nil {
				if err := file.Close(); err != nil {
					return err
				}
				return err
			}

			hdr.Name = filepath.ToSlash(path)

			if err := tw.WriteHeader(hdr); err != nil {
				if err := file.Close(); err != nil {
					return err
				}
				return err
			}

			_, err = io.Copy(tw, file)
			if closeErr := file.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
			return err
		})

		if err != nil {
			if err := tw.Close(); err != nil {
				return
			}
			_ = pw.CloseWithError(err)
			return
		}

		if err := tw.Close(); err != nil {
			_ = pw.CloseWithError(err)
			return
		}

		_ = pw.Close()
	}()

	return pr, nil
}

func walkFilesystem(
	fs billy.Filesystem,
	dir string,
	fn func(path string) error,
) error {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			if err := walkFilesystem(fs, path, fn); err != nil {
				return err
			}
			continue
		}

		if err := fn(path); err != nil {
			return err
		}
	}

	return nil
}

func shouldExclude(
	path string,
	pm *patternmatcher.PatternMatcher,
) bool {
	path = filepath.ToSlash(path)

	matched, err := pm.MatchesOrParentMatches(path)
	if err != nil {
		return false
	}

	return matched
}
