package cli

import (
	"encoding/json"
	"fmt"

	domain "anykey/internal/jsonf/domain"
)

func ParseJSONArray(input []byte) ([]domain.Object, error) {
	var objects []domain.Object
	if err := json.Unmarshal(input, &objects); err != nil {
		return nil, fmt.Errorf("input must be a JSON array of objects: %w", err)
	}
	return objects, nil
}
