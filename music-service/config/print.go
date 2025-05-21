package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func PrintConfig(cfg any) {
	fmt.Println("Configuration:")
	fmt.Println("--------------")
	printReflected(reflect.ValueOf(cfg), 0)
}

func printReflected(v reflect.Value, depth int) {
	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			fmt.Printf("%s<nil>\n", strings.Repeat("  ", depth))
			return
		}
		v = v.Elem()
	}

	// Only process structs
	if v.Kind() != reflect.Struct {
		fmt.Printf("%s%v\n", strings.Repeat("  ", depth), v.Interface())
		return
	}

	t := v.Type()
	fmt.Printf("%s%s:\n", strings.Repeat("  ", depth), t.Name())

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		fmt.Printf("%s%s: ", strings.Repeat("  ", depth+1), fieldType.Name)

		// Handle time.Duration specially
		if field.Type().String() == "time.Duration" {
			fmt.Printf("%v\n", field.Interface().(time.Duration))
			continue
		}

		// Recursively handle nested structs
		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct) {
			fmt.Println()
			printReflected(field, depth+2)
			continue
		}

		// Mask sensitive fields
		if strings.Contains(strings.ToLower(fieldType.Name), "password") ||
			strings.Contains(strings.ToLower(fieldType.Name), "secret") ||
			strings.Contains(strings.ToLower(fieldType.Name), "key") {
			fmt.Println("******")
			continue
		}

		fmt.Printf("%v\n", field.Interface())
	}
}
