package domain

import "encoding/json"

// Object представляет JSON-объект как набор исходных JSON-значений.
type Object = map[string]json.RawMessage
