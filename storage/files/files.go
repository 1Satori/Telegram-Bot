package files

import (
	"Telegram_Bot/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath}
}

func (s Storage) Save(page *storage.Page) error {
	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, 0774); err != nil {
		return fmt.Errorf("mkdir all: %v", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("encode page: %v", err)
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	filePath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, fmt.Errorf("read dir: %v", err)
	}

	if len(files) == 0 {
		return nil, errors.New("no files found")
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(filePath, file.Name()))
}

func (s Storage) Romeve(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("file name: %v", err)
	}

	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("remove file: %v", err)
	}

	return nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("file name: %v", err)
	}

	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("stat file: %v", err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("decode page: %v", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
