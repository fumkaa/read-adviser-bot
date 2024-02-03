package files

import (
	"crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"github.com/fumkaa/read-adviser-bot/storage"
)

type Storage struct {
	basePath string
}

// пользователи имеют доступ на чтение и запись
const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return fmt.Errorf("[Save]can't create directories: %w", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("[Save]can't give nameFile: %w", err)
	}
	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("[Save]can't create file: %w", err)
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("[Save]can't encoding page: %w", err)
	}

	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("[PickRandom]can't read directories: %w", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrnoSavedPage
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(files))))
	if err != nil {
		log.Printf("can't generate rand number: %s", err.Error())
	}
	n := nBig.Int64()

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("[Remove]can't get gileName: %w", err)
	}
	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf(fmt.Sprintf("[Remove]can't remove file: %s", filePath), err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("[IsExists]can't get gileName: %w", err)
	}
	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf(fmt.Sprintf("can't check if file %s exists", filePath), err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("[DecodePage]can't open file: %w", err)
	}
	defer file.Close()

	var page storage.Page

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, fmt.Errorf("[DecodePage]can't decode page: %w", err)
	}
	return &page, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
