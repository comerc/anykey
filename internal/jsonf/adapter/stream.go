package cli

import (
	"bufio"
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

	// Кэш сериализованных ключей; также служит множеством для проверки наличия
	keyBytes := make(map[string][]byte, len(keepFields))
	for _, f := range keepFields {
		if _, ok := keyBytes[f]; ok {
			continue
		}
		b, _ := json.Marshal(f)
		keyBytes[f] = b
	}

	dec := json.NewDecoder(r)
	bw := bufio.NewWriter(w)

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
	if err := bw.WriteByte('['); err != nil {
		return err
	}

	firstElem := true
	for dec.More() {
		if !firstElem {
			if err := bw.WriteByte(','); err != nil {
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

		if err := bw.WriteByte('{'); err != nil {
			return err
		}
		wrote := false
		var raw json.RawMessage
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
			raw = raw[:0]
			if err := dec.Decode(&raw); err != nil {
				return err
			}

			if kb, keep := keyBytes[key]; keep {
				if wrote {
					if err := bw.WriteByte(','); err != nil {
						return err
					}
				}
				if _, err := bw.Write(kb); err != nil {
					return err
				}
				if err := bw.WriteByte(':'); err != nil {
					return err
				}
				if _, err := bw.Write(raw); err != nil {
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
		if err := bw.WriteByte('}'); err != nil {
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
	if err := bw.WriteByte(']'); err != nil {
		return err
	}
	return bw.Flush()
}
