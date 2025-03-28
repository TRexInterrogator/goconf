package goconf

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// load environment variables from file or os
// | c: generic config object to map to
// | ep: relative path to env file, nil == '.env'
func Load[T any](c *T, ep *string) error {
	envPath, err := makeEnvPath(ep)
	if err != nil {
		return err
	}

	// load envs from file or os
	if envFileExists(envPath) {
		vm, err := readEnvFromFile(envPath)
		if err != nil {
			return err
		}
		return assignConfFields(c, vm)
	} else {
		envNames, err := getConfFields(c)
		if err != nil {
			return err
		}

		vm := readEnvFromOS(envNames)
		return assignConfFields(c, vm)
	}
}

// returns type of generic object
func getTypeOfPtr(ptr any) (reflect.Type, error) {
	if ptr == nil {
		return nil, errors.New("getTypeOfPtr, generic was null")
	}

	val := reflect.ValueOf(ptr)

	if val.Kind() != reflect.Ptr {
		return nil, errors.New("getTypeOfPtr, generic not pointer")
	}

	return val.Elem().Type(), nil
}

// translates struct to array of field names
func getConfFields[T any](s T) ([]string, error) {
	typ, err := getTypeOfPtr(s)
	if err != nil {
		return nil, err
	}

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("config can only be mapped to struct")
	}

	numFields := typ.NumField()
	fieldNames := make([]string, numFields)

	for i := range numFields {
		fieldNames[i] = typ.Field(i).Name
	}

	return fieldNames, nil
}

// assigns value map back to struct (strings only)
func assignConfFields[T any](s *T, values map[string]string) error {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return errors.New("input must be a struct pointer")
	}

	for i := range typ.NumField() {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValueStr, ok := values[fieldName]
		if !ok {
			continue // Skip if the field name is not in the map
		}

		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			return fmt.Errorf("cannot set field %s", fieldName)
		}

		switch field.Type.Kind() {
		case reflect.String:
			fieldVal.SetString(fieldValueStr)
		default:
			return fmt.Errorf("unsupported field type for field %s", fieldName)
		}
	}
	return nil
}

// reads local env file and generates value map
func readEnvFromFile(envPath string) (map[string]string, error) {
	file, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return makeValueMap(lines), nil
}

// creates a value map based on env lines
func makeValueMap(lines []string) map[string]string {
	vm := map[string]string{}

	for _, line := range lines {
		ls := strings.Split(line, "=")
		if len(ls) == 2 {
			vm[ls[0]] = ls[1]
			continue
		}
		vm[ls[0]] = ""
	}

	return vm
}

// reads environment variables from OS
func readEnvFromOS(varNames []string) map[string]string {
	vm := map[string]string{}

	for _, name := range varNames {
		vm[name] = os.Getenv(name)
	}

	return vm
}

// checks if local env file exists
func envFileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// constructs path to env file
func makeEnvPath(ep *string) (string, error) {
	env := ".env"
	if ep != nil {
		env = *ep
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, env), nil
}
