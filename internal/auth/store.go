package auth

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const serviceName = "gradmotion-cli"

type Store struct{}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Set(profileName, apiKey string) error {
	return keyring.Set(serviceName, keyName(profileName), apiKey)
}

func (s *Store) Get(profileName string) (string, bool, error) {
	v, err := keyring.Get(serviceName, keyName(profileName))
	if errors.Is(err, keyring.ErrNotFound) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("read keyring failed: %w", err)
	}
	return v, true, nil
}

func (s *Store) Delete(profileName string) error {
	err := keyring.Delete(serviceName, keyName(profileName))
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("delete keyring entry failed: %w", err)
	}
	return nil
}

func keyName(profileName string) string {
	return "profile:" + profileName
}
