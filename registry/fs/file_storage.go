package fs

import (
	"errors"
	"os"
	"path/filepath"
)

type FileStorage struct {
	folder string
}

func (fs *FileStorage) Read(filename string) (string, error) {
	full_path := filepath.Join(fs.folder, filename)

	file, err := os.ReadFile(full_path)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func (fs *FileStorage) Write(filename string, data []byte) error {
	full_path := filepath.Join(fs.folder, filename)

	file, err := os.Create(full_path)
	if err != nil {
		return err
	}

	cnt_read := 0
	for cnt_read < len(data) {
		n, err := file.Write(data[cnt_read:])
		if err != nil {
			return err
		}
		cnt_read += n
	}

	return nil
}

func (fs *FileStorage) Remove(filename string) error {
	full_path := filepath.Join(fs.folder, filename)
	return os.Remove(full_path)
}

func createIfNotExists(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func CreateFileStorage(folder string) (*FileStorage, error) {
	abs_folder, err := filepath.Abs(folder)
	if err != nil {
		return nil, err
	}

	if err := createIfNotExists(abs_folder); err != nil {
		return nil, err
	}

	return &FileStorage{
		folder: abs_folder,
	}, nil
}
