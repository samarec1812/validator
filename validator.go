package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Validate(v any) error {
	typeV := reflect.TypeOf(v).Kind()
	switch typeV {
	case reflect.Struct:
		errorsGroup := make(ValidationErrors, 0)
	TODO:
		for i := 0; i < reflect.TypeOf(v).NumField(); i++ {
			field := reflect.TypeOf(v).Field(i)
			value := reflect.ValueOf(v).Field(i)

			if !field.IsExported() && field.Tag.Get("validate") != "" {
				errorsGroup = append(errorsGroup, ValidationError{ErrValidateForUnexportedFields})
				continue
			} else if field.Tag.Get("validate") == "" {
				continue
			}

			tagsValue := strings.Split(field.Tag.Get("validate"), ";")
			for _, val := range tagsValue {
				tagPair := strings.Split(strings.TrimSpace(val), ":")
				if len(tagPair) != 2 {
					errorsGroup = append(errorsGroup, ValidationError{fmt.Errorf("error with field %s", field.Name)})
					continue TODO
				}
				if err := ValidTag(tagPair[0], tagPair[1]); err != nil {
					errorsGroup = append(errorsGroup, ValidationError{err})
					continue TODO
				}
				// тип поля значение поля тег поля значение тега
				if err := ValidValue(field, value, tagPair[0], tagPair[1]); err != nil {
					errorsGroup = append(errorsGroup, ValidationError{err})
				}
			}
		}
		if len(errorsGroup) > 0 {
			return errorsGroup
		}
		return nil
	default:
		return ErrNotStruct
	}
}

func ValidTag(key, value string) error {
	switch key {
	case "len":
		num, err := strconv.Atoi(value)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if num < 0 {
			return ErrInvalidValidatorSyntax
		}
	case "in":
		if value == "" {
			return ErrInvalidValidatorSyntax
		}
		values := strings.Split(value, ",")
		if len(values) == 0 {
			return ErrInvalidValidatorSyntax

		}
	case "min":
		_, err := strconv.Atoi(value)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
	case "max":
		_, err := strconv.Atoi(value)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
	}

	return nil
}

func ValidValue(field reflect.StructField, value reflect.Value, tag, tagValue string) error {
	switch field.Type.String() {
	case "string":
		err := ValidString(value.String(), tag, tagValue)
		if err != nil {
			return fmt.Errorf("error with field %s: %s", field.Name, err.Error())
		}
	case "int":
		err := ValidInt(value.Int(), tag, tagValue)
		if err != nil {
			return fmt.Errorf("error with field %s: %s", field.Name, err.Error())
		}
	case "[]int":
		for i := 0; i < value.Len(); i++ {
			err := ValidInt(value.Index(i).Int(), tag, tagValue)
			if err != nil {
				return fmt.Errorf("error with field %s: %s", field.Name, err.Error())
			}
		}
	case "[]string":
		for i := 0; i < value.Len(); i++ {
			err := ValidString(value.Index(i).String(), tag, tagValue)
			if err != nil {
				return fmt.Errorf("error with field %s: %s", field.Name, err.Error())
			}
		}
	}
	return nil
}

func ValidString(value, tag, tagValue string) error {
	switch tag {
	case "len":
		num, _ := strconv.Atoi(tagValue)
		if len(value) != num {
			return fmt.Errorf("length does not matching")
		}
	case "in":
		values := strings.Split(tagValue, ",")
		if !Contains(values, value) {
			return fmt.Errorf("value - %s not contain in tag value", value)
		}
	case "min":
		num, _ := strconv.Atoi(tagValue)
		if len(value) < num {
			return fmt.Errorf("length less than required")
		}
	case "max":
		num, _ := strconv.Atoi(tagValue)
		if len(value) > num {
			return fmt.Errorf("length is longer than required")
		}
	}
	return nil
}

func ValidInt(value int64, tag string, tagValue string) error {
	switch tag {
	case "in":
		values, err := PerformToIntSlice(tagValue)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if !Contains(values, value) {
			return fmt.Errorf("value - %d not contain in tag value", value)
		}
	case "min":
		num, _ := strconv.ParseInt(tagValue, 10, 64)
		if value < num {
			return fmt.Errorf("number less than required")
		}
	case "max":
		num, _ := strconv.ParseInt(tagValue, 10, 64)
		if value > num {
			return fmt.Errorf("number is greater than required")
		}
	}
	return nil
}

func PerformToIntSlice(tagValue string) ([]int64, error) {
	arr := make([]int64, 0)

	for _, val := range strings.Split(tagValue, ",") {
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		arr = append(arr, num)
	}
	return arr, nil
}

func Contains[T comparable](t []T, needle T) bool {
	for _, v := range t {
		if v == needle {
			return true
		}
	}
	return false
}
