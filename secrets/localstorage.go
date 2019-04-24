package secrets

import (
	"github.com/juju/errors"
)

type LocalStorageConfig struct {
	Values map[string]interface{} `json:"values"`
}

type LocalStorage struct {
	cfg *LocalStorageConfig
}

func NewLocalStorage(cfg *LocalStorageConfig) *LocalStorage {
	return &LocalStorage{cfg: cfg}
}

func (s *LocalStorage) ReadStringValue(path string) (string, error) {
	value, ok := s.cfg.Values[path]
	if !ok {
		return "", errors.NotFoundf("Key '%s' not found", path)
	}

	str, ok := value.(string)
	if !ok {
		return "", errors.NotValidf("Key '%s' is not a string", path)
	}

	return str, nil
}

func (s *LocalStorage) ReadStringArray(path string) ([]string, error) {
	value, ok := s.cfg.Values[path]
	if !ok {
		return nil, errors.NotFoundf("Key '%s' not found", path)
	}

	arr, ok := value.([]interface{})
	if !ok {
		return nil, errors.NotValidf("Key '%s' is not a string array", path)
	}

	retval := make([]string, 0)
	for _, i := range arr {
		tmp, ok := i.(string)
		if !ok {
			return nil, errors.NotValidf("Unable to convert key to string")
		}
		retval = append(retval, tmp)
	}

	return retval, nil
}
