package store

import (
	"encoding/json"
	"errors"
	"os"
)

type Store interface {
	Read(data interface{}) (err error)
	Write(data interface{}) (err error)
}

type Type string

const (
	FileType Type = "filestorage"
)

type FileStore struct {
	FileName string
}

func (fs *FileStore) Read(data interface{}) (err error) {
	file, err := os.ReadFile(fs.FileName)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		file = []byte("[]")
	}

	err = json.Unmarshal(file, data)

	return err
}

func (fs *FileStore) Write(data interface{}) (err error) {
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fs.FileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	_, err = f.Write(fileData)
	if err != nil {
		return err
	}
	return nil
}

func NewStorage(storageType Type, fileName string) Store {
	switch storageType {
	case FileType:
		return &FileStore{FileName: fileName}
	}

	return nil
}
