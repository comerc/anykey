package cli

import (
	"encoding/json"
	"fmt"
	"io"
)

// StreamFilterAndWrite читает из r JSON-массив объектов и записывает
// отфильтрованный JSON-массив в w, сохраняя исходный порядок полей в каждом объекте.
// Остаются только поля, перечисленные в keepFields. Дубликаты имён полей
// в keepFields игнорируются (берётся первое вхождение).
func StreamFilterAndWrite(r io.Reader, keepFields []string, w io.Writer) error {
	if len(keepFields) == 0 {
		return fmt.Errorf("no fields specified; pass field names as arguments")
	}

	keepSet := make(map[string]struct{}, len(keepFields))
	for _, f := range keepFields {
		keepSet[f] = struct{}{}
	}

	// Кэшируем сериализованные имена ключей, чтобы не маршалить на каждом объекте
	keyBytes := make(map[string][]byte, len(keepFields))
	for _, f := range keepFields {
		if _, ok := keyBytes[f]; ok {
			continue
		}
		b, _ := json.Marshal(f)
		keyBytes[f] = b
	}

	dec := json.NewDecoder(r)

	// Ожидаем начало массива
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '[' {
		return fmt.Errorf("input must be a JSON array of objects: got %v", tok)
	}

	// Пишем начало массива
	if _, err := w.Write([]byte{'['}); err != nil {
		return err
	}

	firstElem := true
	for dec.More() {
		if !firstElem {
			if _, err := w.Write([]byte{','}); err != nil {
				return err
			}
		}
		firstElem = false

		// Ожидаем начало объекта
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		d, ok := tok.(json.Delim)
		if !ok || d != '{' {
			return fmt.Errorf("expected object start, got %v", tok)
		}

		if _, err := w.Write([]byte{'{'}); err != nil {
			return err
		}
		wrote := false
		for dec.More() {
			// Читаем ключ
			ktok, err := dec.Token()
			if err != nil {
				return err
			}
			key, ok := ktok.(string)
			if !ok {
				return fmt.Errorf("expected string key, got %v", ktok)
			}

			// Читаем значение как "сырой" JSON
			var raw json.RawMessage
			if err := dec.Decode(&raw); err != nil {
				return err
			}

			if _, keep := keepSet[key]; keep {
				if wrote {
					if _, err := w.Write([]byte{','}); err != nil {
						return err
					}
				}
				if _, err := w.Write(keyBytes[key]); err != nil {
					return err
				}
				if _, err := w.Write([]byte{':'}); err != nil {
					return err
				}
				if _, err := w.Write(raw); err != nil {
					return err
				}
				wrote = true
			}
		}

		// Читаем конец объекта
		tok, err = dec.Token()
		if err != nil {
			return err
		}
		d, ok = tok.(json.Delim)
		if !ok || d != '}' {
			return fmt.Errorf("expected object end, got %v", tok)
		}
		if _, err := w.Write([]byte{'}'}); err != nil {
			return err
		}
	}

	// Ожидаем конец массива
	tok, err = dec.Token()
	if err != nil {
		return err
	}
	delim, ok = tok.(json.Delim)
	if !ok || delim != ']' {
		return fmt.Errorf("expected array end, got %v", tok)
	}
	if _, err := w.Write([]byte{']'}); err != nil {
		return err
	}
	return nil
}
