package utils

import (
	"errors"
	"fmt"
	"logger"
	"reflect"
	"runtime/debug"
	"unicode"
)

func IsExportedName(name string) bool {
	return name != "" && unicode.IsUpper(rune(name[0]))
}
func IsExportedField(structField reflect.StructField) bool {
	return structField.PkgPath == ""
}

// CopyExportedStructFields copies all exported struct fields from src
// that are assignable to their name siblings at dstPtr to dstPtr.
// src can be a struct or a pointer to a struct, dstPtr must be
// a pointer to a struct.
func CopyExportedStructFields(src, dstPtr interface{}) (copied int) {
	vsrc := reflect.ValueOf(src)
	if vsrc.Kind() == reflect.Ptr {
		vsrc = vsrc.Elem()
	}
	vdst := reflect.ValueOf(dstPtr).Elem()
	return CopyExportedStructFieldsVal(vsrc, vdst)
}

func ExportedStructFields(v reflect.Value) map[string]reflect.Value {
	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("Expected a struct, got %s", t))
	}
	result := make(map[string]reflect.Value)
	exportedStructFields(v, t, result)
	return result
}

func exportedStructFields(v reflect.Value, t reflect.Type, result map[string]reflect.Value) {
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		if IsExportedField(structField) {
			if structField.Anonymous && structField.Type.Kind() == reflect.Struct {
				exportedStructFields(v.Field(i), structField.Type, result)
			} else {
				result[structField.Name] = v.Field(i)
			}
		}
	}
}

func CopyExportedStructFieldsVal(src, dst reflect.Value) (copied int) {
	if src.Kind() != reflect.Struct {
		panic(fmt.Errorf("CopyExportedStructFieldsVal: src must be struct, got %s", src.Type()))
	}
	if dst.Kind() != reflect.Struct {
		panic(fmt.Errorf("CopyExportedStructFieldsVal: dst must be struct, got %s", dst.Type()))
	}
	if !dst.CanSet() {
		panic(fmt.Errorf("CopyExportedStructFieldsVal: dst (%s) is not set-able", dst.Type()))
	}
	srcFields := ExportedStructFields(src)
	dstFields := ExportedStructFields(dst)
	for name, srcV := range srcFields {
		if dstV, ok := dstFields[name]; ok {
			if srcV.Type().AssignableTo(dstV.Type()) {
				dstV.Set(srcV)
				copied++
			}
		}
	}
	return copied
}

func IsEmptyValue(v interface{}) bool {
	reflectValue := reflect.ValueOf(v)

	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return isEmptyValue(reflectValue)
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func Invoke(any interface{}, name string, args ...interface{}) (result []reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Painc error : %v", err)
			trackBack := string(debug.Stack())
			logger.Error(trackBack)
			result = append(result, reflect.ValueOf(err))
		}
	}()
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	result = reflect.ValueOf(any).MethodByName(name).Call(inputs)
	return
}

func GetMethods(v interface{}) map[string][]string {
	funcMap := make(map[string][]string)
	reflectType := reflect.TypeOf(v)
	//reflectkind := reflectType.Kind()
	// if reflectkind == reflect.Ptr {
	// 	reflectType = reflectType.Elem()
	// }
	for i := 0; i < reflectType.NumMethod(); i++ {
		funcList := make([]string, 0)
		method := reflectType.Method(i)
		methodType := method.Type
		methodName := method.Name
		if !IsExportedName(methodName) {
			continue
		}
		for j := 1; j < methodType.NumIn(); j++ {
			params := methodType.In(j)
			funcList = append(funcList, params.Name())
		}
		funcMap[methodName] = funcList
	}
	logger.Debug("GetMethods:%v", funcMap)
	return funcMap
}

// GetField returns the value of the provided obj field. obj can whether
// be a structure or pointer to structure.
func GetField(obj interface{}, name string) (interface{}, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	field := objValue.FieldByName(name)
	if !field.IsValid() {
		return nil, fmt.Errorf("No such field: %s in obj", name)
	}

	return field.Interface(), nil
}

// GetFieldKind returns the kind of the provided obj field. obj can whether
// be a structure or pointer to structure.
func GetFieldKind(obj interface{}, name string) (reflect.Kind, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return reflect.Invalid, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	field := objValue.FieldByName(name)

	if !field.IsValid() {
		return reflect.Invalid, fmt.Errorf("No such field: %s in obj", name)
	}

	return field.Type().Kind(), nil
}

// GetFieldTag returns the provided obj field tag value. obj can whether
// be a structure or pointer to structure.
func GetFieldTag(obj interface{}, fieldName, tagKey string) (string, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return "", errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()

	field, ok := objType.FieldByName(fieldName)
	if !ok {
		return "", fmt.Errorf("No such field: %s in obj", fieldName)
	}

	if !IsExportedField(field) {
		return "", errors.New("Cannot GetFieldTag on a non-exported struct field")
	}

	return field.Tag.Get(tagKey), nil
}

// SetField sets the provided obj field with provided value. obj param has
// to be a pointer to a struct, otherwise it will soundly fail. Provided
// value type should match with the struct field you're trying to set.
func SetField(obj interface{}, name string, value interface{}) error {
	// Fetch the field reflect.Value
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	// If obj field value is not settable an error is thrown
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		invalidTypeError := errors.New("Provided value type didn't match obj field type")
		return invalidTypeError
	}

	structFieldValue.Set(val)
	return nil
}

// HasField checks if the provided field name is part of a struct. obj can whether
// be a structure or pointer to structure.
func HasField(obj interface{}, name string) (bool, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return false, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	field, ok := objType.FieldByName(name)
	if !ok || !IsExportedField(field) {
		return false, nil
	}

	return true, nil
}

// Fields returns the struct fields names list. obj can whether
// be a structure or pointer to structure.
func Fields(obj interface{}) ([]string, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	var fields []string
	for i := 0; i < fieldsCount; i++ {
		field := objType.Field(i)
		if IsExportedField(field) {
			fields = append(fields, field.Name)
		}
	}

	return fields, nil
}

// Items returns the field - value struct pairs as a map. obj can whether
// be a structure or pointer to structure.
func Items(obj interface{}) (map[string]interface{}, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	items := make(map[string]interface{})

	for i := 0; i < fieldsCount; i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)

		// Make sure only exportable and addressable fields are
		// returned by Items
		if IsExportedField(field) {
			items[field.Name] = fieldValue.Interface()
		}
	}

	return items, nil
}

// Tags lists the struct tag fields. obj can whether
// be a structure or pointer to structure.
func Tags(obj interface{}, key string) (map[string]string, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	tags := make(map[string]string)

	for i := 0; i < fieldsCount; i++ {
		structField := objType.Field(i)

		if IsExportedField(structField) {
			tags[structField.Name] = structField.Tag.Get(key)
		}
	}

	return tags, nil
}

func reflectValue(obj interface{}) reflect.Value {
	var val reflect.Value

	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		val = reflect.ValueOf(obj).Elem()
	} else {
		val = reflect.ValueOf(obj)
	}

	return val
}

func hasValidType(obj interface{}, types []reflect.Kind) bool {
	for _, t := range types {
		if reflect.TypeOf(obj).Kind() == t {
			return true
		}
	}

	return false
}

func isStruct(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Struct
}

func isPointer(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Ptr
}
