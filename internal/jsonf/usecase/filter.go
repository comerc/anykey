package usecase

import (
	"anykey/internal/jsonf/domain"
)

// DedupeKeepOrder удаляет дубликаты и сохраняет исходный порядок полей
func DedupeKeepOrder(keepFields []string) []string {
	dedup := make([]string, 0, len(keepFields))
	seen := make(map[string]struct{}, len(keepFields))
	for _, f := range keepFields {
		if _, ok := seen[f]; ok {
			continue
		}
		seen[f] = struct{}{}
		dedup = append(dedup, f)
	}
	return dedup
}

// FilterObjects возвращает новые объекты, содержащие только поля из keepOrder
// Порядок полей сохраняется на этапе кодирования; здесь просто удаляются лишние поля.
func FilterObjects(objects []domain.Object, keepOrder []string) []domain.Object {
	filtered := make([]domain.Object, 0, len(objects))
	keep := make(map[string]struct{}, len(keepOrder))
	for _, k := range keepOrder {
		keep[k] = struct{}{}
	}
	for _, obj := range objects {
		out := make(domain.Object, len(keepOrder))
		for k := range obj {
			if _, ok := keep[k]; ok {
				out[k] = obj[k]
			}
		}
		filtered = append(filtered, out)
	}
	return filtered
}
