package cli

import (
	"bufio"
	"encoding/json"
	"io"

	domain "anykey/internal/jsonf/domain"
)

// WriteObjectsOrdered пишет массив объектов, оставляя только указанные поля и сохраняя их порядок
func WriteObjectsOrdered(w io.Writer, objects []domain.Object, keepOrder []string) error {
	bw := bufio.NewWriter(w)
	bw.WriteByte('[')
	for i, obj := range objects {
		if i > 0 {
			bw.WriteByte(',')
		}
		bw.WriteByte('{')
		wrote := false
		for _, field := range keepOrder {
			if val, ok := obj[field]; ok {
				if wrote {
					bw.WriteByte(',')
				}
				keyBytes, _ := json.Marshal(field)
				bw.Write(keyBytes)
				bw.WriteByte(':')
				bw.Write(val)
				wrote = true
			}
		}
		bw.WriteByte('}')
	}
	bw.WriteByte(']')
	return bw.Flush()
}
