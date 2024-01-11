package services

import (
	"context"
	"log"

	"github.com/Yaxhveer/go-url-shortner/util"
)

type DBStore interface {
	GetURL(ctx context.Context, shortURL string) (string, error)
	GetShortURL(ctx context.Context, url string) (string, error)
	AddURL(ctx context.Context, url string, shortURL string) error
	HaveURL(ctx context.Context, url string) (bool, error)
}

type service struct {
	storage DBStore
}

func (s *service) GetURL(ctx context.Context, shortURL string) (string, error) {
	url, err := s.storage.GetURL(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *service) ShortenURL(ctx context.Context, url string) (string, error) {
	contains, err := s.storage.HaveURL(ctx, url)
	if err != nil {
		return "", nil
	}

	if contains {
		log.Println(url, "url with shorten url already exists")

		shortURL, err := s.storage.GetShortURL(ctx, url)
		if err != nil {
			return "", err
		}
		return shortURL, nil
	}

	shortURL := util.GenerateShortenURL();
	
	err = s.storage.AddURL(ctx, url, shortURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func NewService(s DBStore) *service {
	return &service{storage: s}
}
