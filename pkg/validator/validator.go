package validator

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var v10 *validator.Validate
var schemaDecoder = schema.NewDecoder()

var customValidators = map[string]func(fl validator.FieldLevel) bool{
	"snake_case": validateSnakeCase,
	"not_empty":  validateNotEmpty,
	"date":       validateDate,
	"tag":        validateTag,
}

func init() {
	v10 = validator.New(validator.WithRequiredStructEnabled())
	for tag, fn := range customValidators {
		_ = v10.RegisterValidation(tag, fn)
	}
}

func BindJSON(object any, r *http.Request) error {
	if isFormRequest(r) {
		if err := r.ParseForm(); err != nil {
			return err
		}
		schemaDecoder.IgnoreUnknownKeys(true)
		if err := schemaDecoder.Decode(object, r.Form); err != nil {
			return err
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(object); err != nil {
			return err
		}
	}

	return Validate(object)
}

func isFormRequest(r *http.Request) bool {
	return r.Method == http.MethodGet || strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data")
}

func Validate(object any) error {
	err := v10.Struct(object)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		return Error{
			Msg: err.Error(),
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	fieldErrors := map[string]string{}
	for _, fieldErr := range validationErrors {
		key := buildPath(reflect.TypeOf(object).Elem(), prepareNamespace(fieldErr.Namespace()))
		fieldErrors[key] = fieldErr.Tag()
		if fieldErr.Param() != "" {
			fieldErrors[key] += "=" + fieldErr.Param()
		}
	}

	return Error{
		Fields: fieldErrors,
	}
}

func buildPath(objectType reflect.Type, namespace []string) string {
	field := namespace[0]
	if _, err := strconv.Atoi(field); err == nil {
		if len(namespace) > 1 {
			return field + "." + buildPath(objectType.Elem(), namespace[1:])
		}
		return field
	}

	var f reflect.StructField
	if objectType.Kind() == reflect.Ptr {
		f, _ = objectType.Elem().FieldByName(field)
	} else {
		f, _ = objectType.FieldByName(field)
	}

	tag := getJSONTag(f.Tag)
	path := tag

	if len(namespace) > 1 {
		path += "." + buildPath(f.Type, namespace[1:])
	}

	return path
}

func prepareNamespace(namespace string) []string {
	namespace = strings.SplitN(namespace, ".", 2)[1]
	namespace = strings.ReplaceAll(strings.ReplaceAll(namespace, "[", "."), "]", "")

	return strings.Split(namespace, ".")
}

func getJSONTag(tag reflect.StructTag) string {
	if val, ok := tag.Lookup("schema"); ok {
		return val
	}
	return strings.Split(tag.Get("json"), ",")[0]
}
