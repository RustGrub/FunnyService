package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func IntValueFromRequest(r *http.Request, key string) (*int, error) {
	rVal := r.URL.Query().Get(key)

	if rVal == "" {
		return nil, errors.New("empty key: " + key)
	}
	intVal, err := strconv.Atoi(rVal)
	if err != nil {
		return nil, err
	}
	return &intVal, nil
}

func IsEmptyOrOnlySpacesString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
