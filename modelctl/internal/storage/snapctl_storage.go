package storage

import (
	"encoding/json"
	"strings"

	"github.com/canonical/go-snapctl"
)

type SnapctlStorage struct{}

func NewSnapctlStorage() *SnapctlStorage {
	return &SnapctlStorage{}
}

func (s *SnapctlStorage) Set(key, value string) error {
	return snapctl.Set(key, string(value)).Run()
}

func (s *SnapctlStorage) SetDocument(key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return snapctl.Set(key, string(b)).Document().Run()
}

func (s *SnapctlStorage) Get(key string) (map[string]any, error) {
	valJson, err := snapctl.Get(key).Run()
	if err != nil {
		return nil, err
	}
	if valJson == "" {
		return nil, ErrorNotFound
	}

	var valMap map[string]any
	if strings.HasPrefix(valJson, "{") && strings.HasSuffix(valJson, "}") {
		// Object value, parse as JSON
		err = json.Unmarshal([]byte(valJson), &valMap)
		if err != nil {
			return nil, err
		}
	} else {
		// Primitive value, return as-is
		valMap = map[string]any{
			key: valJson,
		}
	}

	return valMap, nil
}

func (s *SnapctlStorage) Unset(key string) error {
	err := snapctl.Unset(key).Run()
	if err != nil {
		return err
	}
	return nil
}
