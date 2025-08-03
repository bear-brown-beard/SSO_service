package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ParseAndValidate(r *http.Request, data interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("Ошибка при чтении тела запроса")
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, data); err != nil {
		return err
	}

	if err := validate.Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New("Ошибка в поле '" + err.Field() + "': " + err.Tag())
		}
	}

	return nil
}
func ParseIntOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func ParseIntSlice(value string) []int {
	if value == "" {
		return nil
	}
	strSlice := strings.Split(value, ",")
	result := make([]int, 0, len(strSlice))

	for _, str := range strSlice {
		if num, err := strconv.Atoi(str); err == nil {
			result = append(result, num)
		}
	}
	return result
}
