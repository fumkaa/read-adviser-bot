package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(username string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

var ErrnoSavedPage = fmt.Errorf("no saved page")

func (p Page) Hash() (string, error) {
	h := sha1.New()

	_, err := io.WriteString(h, p.URL)
	if err != nil {
		return "", fmt.Errorf("can't write p.URL in h: %w", err)
	}
	_, err = io.WriteString(h, p.UserName)
	if err != nil {
		return "", fmt.Errorf("can't write p.UserName in h: %w", err)
	}

	return fmt.Sprintf("%x", h), nil
}
