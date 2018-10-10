package serializer

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"strings"
)

const (
	ID_PREFIX = "@zenaton#"
)

type format struct {
	Data   interface{}   `json:"d,omitempty"`
	Object string        `json:"o,omitempty"`
	Store  []StoreObject `json:"s"`
}

type StoreObject struct {
	Name       string                 `json:"n,omitempty"`
	Keys       []interface{}          `json:"k,omitempty"`
	Values     []interface{}          `json:"v,omitempty"`
	Properties map[string]interface{} `json:"p,omitempty"`
}

// Serializer is just a type that allows you to call Encode and Decode as methods for convenience.
type Serializer struct{}

func (s *Serializer) Encode(data interface{}) (string, error)     { return Encode(data) }
func (s *Serializer) Decode(data string, value interface{}) error { return Decode(data, value) }

type serializer struct {
	encoded  []StoreObject
	decoded  []reflect.Value
	pointers []uintptr
}

type Object struct {
	Name       string                 `json:"n"`
	Properties map[string]interface{} `json:"p"`
}

func Encode(data interface{}) (string, error) {

	bytes, err := json.Marshal(data)
	return string(bytes), err

	//rv := reflect.ValueOf(data)
	//kind := rv.Kind()
	//isValid := validType(kind)
	//if !isValid {
	//	return "", errors.New(fmt.Sprintf("cannot encode data of kind: %s", kind.String()))
	//}
	//
	//s := serializer{}
	//s.encoded = []StoreObject{}
	//s.pointers = []uintptr{}

	//return s.encode(rv, data)
}

func (s *serializer) encode(rv reflect.Value, data interface{}) (string, error) {

	value := format{}

	if reflect.TypeOf(data) == nil {
		//special case, as the json marshal would normally remove the "d". but here we want it to show that the data is nil and not an empty string
		return `{}`, nil
	} else if basicType(rv) {
		value.Data = data
	} else {
		value.Object = s.encodeToStore(rv)
	}

	value.Store = s.encoded
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(encoded), nil

}

func basicType(rv reflect.Value) bool {
	v := rv
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	kind := v.Kind()
	return kind == reflect.Bool ||
		kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 ||
		kind == reflect.Uintptr ||
		kind == reflect.Float32 ||
		kind == reflect.Float64 ||
		kind == reflect.String
}

func (s *serializer) encodeToStore(object reflect.Value) string {
	for object.Kind() == reflect.Interface {
		object = object.Elem()
	}
	//fmt.Println("kind in encodeToStore: ", object.Kind())
	if object.Kind() == reflect.Ptr {
		if object.IsNil() {
			return ""
		}
		id := indexOf(s.pointers, object.Pointer())
		//fmt.Println("id: ", id)
		if id != -1 {
			return storeID(id)
		}
	}
	return s.storeAndEncode(object)
}

func (s *serializer) storeAndEncode(object reflect.Value) string {
	id := len(s.pointers)
	if object.Kind() == reflect.Ptr {
		s.pointers = insertPtr(s.pointers, object.Pointer(), id)
	} else {
		// this pointer is never actually used. It is only added so that the length of pointers is correct
		s.pointers = insertPtr(s.pointers, reflect.ValueOf(&object).Pointer(), id)
	}

	s.encoded = insert(s.encoded, s.encodedObjectByType(object), id)
	return storeID(id)
}

func (s *serializer) encodedObjectByType(object reflect.Value) StoreObject {
	object = reflect.Indirect(object)
	kind := object.Kind()
	//fmt.Println("kind in encodedObjectByType: ", kind)
	switch kind {
	case reflect.Struct:
		return s.encodeStruct(object)
	case reflect.Array, reflect.Slice:
		return s.encodeArray(object)
	case reflect.Map:
		return s.encodeMap(object)
	}

	return StoreObject{}
}

func storeID(id int) string {
	return ID_PREFIX + strconv.Itoa(id)
}

func validType(kind reflect.Kind) bool {

	switch kind {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Array,
		reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.String,
		reflect.Struct, reflect.Invalid:
		return true
	case reflect.UnsafePointer, reflect.Complex64, reflect.Complex128, reflect.Chan,
		reflect.Func:
		return false
	}
	return false
}

func (s *serializer) encodeArray(a reflect.Value) StoreObject {

	var array []interface{}
	for i := 0; i < a.Len(); i++ {
		rv := a.Index(i)
		//fmt.Println("rv: ", rv)
		for rv.Kind() == reflect.Interface && rv.Elem().Kind() != reflect.Invalid {
			rv = rv.Elem()
		}
		kind := rv.Kind()
		//fmt.Println("kind in encodeArray2::::: ", kind)
		if basicType(rv) || kind == reflect.Interface {
			fmt.Println("basic kind or interface: ", rv.Interface(), rv)
			array = append(array, rv.Interface())
			continue
		}
		array = append(array, s.encodeToStore(rv))
	}

	return StoreObject{
		Values: array,
	}
}

func (s *serializer) encodeMap(m reflect.Value) StoreObject {

	var keys []interface{}
	var values []interface{}

	keyValues := m.MapKeys()
	for _, kv := range keyValues {
		for kv.Kind() == reflect.Interface {
			kv = kv.Elem()
		}
		//fmt.Println("key::::: ", kind)
		if basicType(kv) {
			keys = append(keys, kv.Interface())
		} else {
			keys = append(keys, s.encodeToStore(kv))
		}

		//todo: abstract this out into another function, as I keep doing this
		valueValue := m.MapIndex(kv)
		for valueValue.Kind() == reflect.Interface {
			valueValue = valueValue.Elem()
		}
		//fmt.Println("value::::: ", valueValue.Kind())
		if basicType(valueValue) {
			values = append(values, valueValue.Interface())
		} else {
			values = append(values, s.encodeToStore(valueValue))
		}
	}

	return StoreObject{
		Keys:   keys,
		Values: values,
	}
}

func (s *serializer) encodeStruct(object reflect.Value) StoreObject {

	return StoreObject{
		Name:       object.Type().Name(),
		Properties: s.encodeProperties(object),
	}
}

func (s *serializer) encodeProperties(o reflect.Value) map[string]interface{} {
	dataT := o.Type()
	propMap := make(map[string]interface{})
	for i := 0; i < o.NumField(); i++ {
		key := dataT.Field(i).Name
		//fmt.Println("1key: ", key, "kind: ", o.Field(i).Kind())
		if basicType(o.Field(i)) {
			if o.Field(i).CanInterface() {
				propMap[key] = o.Field(i).Interface()
			}
			continue
		}
		//fmt.Println("")
		propMap[key] = s.encodeToStore(o.Field(i))
	}
	return propMap
}

func indexOf(slice []uintptr, item uintptr) int {

	//fmt.Println("pointers: ", slice, "item: ", item)
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func Decode(data string, value interface{}) error {

	return json.Unmarshal([]byte(data), value)

	//var parsedJSON format
	//err := json.Unmarshal([]byte(data), &parsedJSON)
	//if err != nil {
	//	return err
	//}
	//
	//rv := reflect.ValueOf(value)
	//
	////value must be a pointer
	//if rv.Kind() != reflect.Ptr || rv.IsNil() {
	//	return errors.New("serializer.Decode: must use a pointer value")
	//}
	//
	//s := serializer{}
	//return s.decode(rv, parsedJSON)
}

func (s *serializer) decode(rv reflect.Value, parsedJSON format) error {
	s.encoded = parsedJSON.Store

	simpleValue := parsedJSON.Data

	if simpleValue != nil {
		switch rv.Elem().Kind() {
		case reflect.Bool:
			rv.Elem().SetBool(simpleValue.(bool))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rv.Elem().SetInt(int64(simpleValue.(float64)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			rv.Elem().SetUint(uint64(simpleValue.(float64)))
		case reflect.Float32, reflect.Float64:
			rv.Elem().SetFloat(simpleValue.(float64))
		case reflect.String:
			rv.Elem().SetString(simpleValue.(string))
		case reflect.Interface:
			if simpleValue != nil {
				rv.Elem().Set(reflect.ValueOf(simpleValue))
			}
		default:
			panic(fmt.Sprintf("unknown kind: %s", rv.Elem().Kind()))
		}
		return nil
	}

	//fmt.Println(4)
	id := parsedJSON.Object
	if id != "" {
		idInt, err := strconv.Atoi(strings.TrimLeft(id, ID_PREFIX))
		if err != nil {
			return err
		}
		return s.decodeFromStore(idInt, s.encoded[idInt], rv)
	}

	return nil
}

func (s *serializer) decodeFromStore(id int, encoded StoreObject, rv reflect.Value) error {

	//fmt.Println("id, decoded: ", id, s.decoded)

	if len(s.decoded) > id {
		decoded := s.decoded[id]
		//fmt.Println("********", decoded.Kind(), rv.Kind())
		//rv.Set(indirect(decoded))
		rv.Set(decoded)
		return nil
	}

	if encoded.Properties != nil {
		//fmt.Println("in the thing2", rv, rv.Kind())
		s.decodeStruct(id, encoded.Properties, rv)
		return nil
	}

	if encoded.Keys != nil {
		s.decodeMap(id, encoded.Keys, encoded.Values, rv)
		return nil
	}

	if encoded.Values != nil {
		s.decodeArray(id, encoded.Values, rv)
		return nil
	}

	return errors.New("serializer.Decode: wrong format of data")
}

func (s *serializer) decodeArray(id int, array interface{}, rv reflect.Value) {

	arr := array.([]interface{})

	var newRV reflect.Value
	switch rv.Kind() {
	case reflect.Interface:
		var newSlice []interface{}
		newRV = reflect.ValueOf(&newSlice)
		//fmt.Println("interface(((((((((((((((((((((")
	default:
		//fmt.Println("default(((((((((((((((((((((")
		newRV = rv
	}

	//fmt.Println("in the thing", rv, rv.Kind())
	s.decoded = insertRV(s.decoded, newRV, id)
	newRV = indirect(newRV)

	for i, arrV := range arr {

		// Get element of array, growing if necessary.
		if newRV.Kind() == reflect.Slice {

			// Grow slice if necessary
			if i >= newRV.Cap() {
				newcap := newRV.Cap() + newRV.Cap()/2
				if newcap < 4 {
					newcap = 4
				}
				newv := reflect.MakeSlice(newRV.Type(), newRV.Len(), newcap)
				reflect.Copy(newv, newRV)
				newRV.Set(newv)
			}
			if i >= newRV.Len() {
				newRV.SetLen(i + 1)
			}
		}

		if i < newRV.Len() {
			// Decode into element.
			s.decodeElement(newRV.Index(i), arrV)
		} else {
			panic("shouldn't get here")
		}
	}

	if rv.CanAddr() {
		rv.Set(s.decoded[id])
	} else {
		rv.Elem().Set(indirect(s.decoded[id]))
	}
}

//todo: should I handle the case in which the KEY_OBJECT_PROPERTIES don't match the struct passed in? seems yes/. actually just do like the json package does
func (s *serializer) decodeStruct(id int, encodedObject interface{}, v reflect.Value) {

	object := encodedObject.(map[string]interface{})

	newV := v
	s.decoded = insertRV(s.decoded, newV, id)
	newV = indirect(newV)

	for key, value := range object {
		field := indirect(newV).FieldByName(key)
		//fmt.Println("type of field: ", field)
		//fmt.Println("in the thing", field, field.Kind())
		s.decodeElement(field, value)
	}

	if v.CanAddr() {
		v.Set(s.decoded[id])
	} else {
		v.Elem().Set(indirect(s.decoded[id]))
	}
}

func (s *serializer) decodeMap(id int, keys interface{}, values interface{}, v reflect.Value) {

	ks := keys.([]interface{})
	vs := values.([]interface{})

	var newV reflect.Value
	switch indirect(v).Kind() {
	case reflect.Interface:
		newMap := make(map[interface{}]interface{})
		newV = reflect.ValueOf(&newMap)
		//fmt.Println("interface(((((((((((((((((((((")
	default:
		//fmt.Println("default(((((((((((((((((((((")
		newV = v
	}

	//fmt.Println("the thing:::::::::: ", newV, newV.Kind())
	s.decoded = append(s.decoded, newV)
	newV = indirect(newV)
	if newV.IsNil() {
		newV.Set(reflect.MakeMap(newV.Type()))
	}

	for i, k := range ks {
		v := vs[i]

		//fmt.Println("newV.Type()", newV.Type())

		newKey := reflect.New(newV.Type().Key()).Elem()
		newValue := reflect.New(newV.Type().Elem()).Elem()
		s.decodeElement(newKey, k)
		s.decodeElement(newValue, v)
		newV.SetMapIndex(newKey, newValue)
	}

	if v.CanAddr() {
		v.Set(indirect(newV))
	} else {
		v.Elem().Set(indirect(newV))
	}
}

func (s *serializer) decodeElement(rv reflect.Value, value interface{}) {

	potentialID, ok := value.(string)
	if ok {
		id, isStoreID := s.storeID(potentialID)
		if isStoreID {
			encoded := s.encoded[id]
			s.decodeFromStore(id, encoded, rv)
			return
		}
	}

	rv = indirect(rv)
	switch rv.Kind() {
	case reflect.Bool:
		rv.SetBool(value.(bool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv.SetInt(int64(value.(float64)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv.SetUint(uint64(value.(float64)))
	case reflect.Float32, reflect.Float64:
		rv.Set(reflect.ValueOf(value).Convert(rv.Type()))
	case reflect.String:
		rv.SetString(value.(string))
	case reflect.Uintptr:
		rv.Set(reflect.ValueOf(uintptr(value.(float64))))
	case reflect.Interface:
		rv.Set(reflect.ValueOf(value))
	case reflect.Ptr:
		rv.Set(reflect.ValueOf(value))
	case reflect.Invalid:
		panic("this should never be invalid")
	case reflect.Array, reflect.Slice:
		//todo? s.decodeLegacyArray(value, rv)
	case reflect.Struct:
		//todo? why do I have to do nothing here?
	//case reflect.Complex64, reflect.Complex128: 	//todo: this is not supported by json i guess?
	//	var c complex128
	//	str := fmt.Sprintf(`"%s"`, value.(string))
	//	err := json.Unmarshal([]byte(str), &c)
	//	if err != nil {
	//		//todo: panic?
	//		panic(err)
	//	}
	//	fld.SetComplex(c)
	//case reflect.UnsafePointer:

	default:
		panic(fmt.Sprintf("unknown kind: %s", rv.Kind()))
	}
}

func (s *serializer) storeID(str string) (int, bool) {
	if !strings.HasPrefix(str, ID_PREFIX) {
		return 0, false
	}
	id := strings.TrimLeft(str, ID_PREFIX)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return idInt, false
	}
	return idInt, idInt <= len(s.encoded)
}

func setPtr(field reflect.Value, value interface{}) {
	switch field.Type().Elem().Kind() {
	case reflect.Ptr:
		//todo: recurse
	case reflect.Bool:
		v := value.(bool)
		field.Set(reflect.ValueOf(&v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := int64(value.(float64))
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := uint64(value.(float64))
		field.Set(reflect.ValueOf(&v))
	case reflect.Float32, reflect.Float64:
		v := value.(float64)
		field.Set(reflect.ValueOf(&v))
	case reflect.String:
		v := value.(string)
		field.Set(reflect.ValueOf(&v))
	}
	//todo: other possible kinds
}

func insert(arr []StoreObject, value StoreObject, i int) []StoreObject {
	if len(arr) > i {
		arr[i] = value
		return arr
	}
	newArr := make([]StoreObject, i+1)
	copy(newArr, arr)
	newArr[i] = value
	return newArr
}

func insertRV(arr []reflect.Value, value reflect.Value, i int) []reflect.Value {
	if len(arr) > i {
		arr[i] = value
		return arr
	}
	newArr := make([]reflect.Value, i+1)
	copy(newArr, arr)
	newArr[i] = value
	return newArr
}

func insertPtr(arr []uintptr, value uintptr, i int) []uintptr {
	if len(arr) > i {
		arr[i] = value
		return arr
	}
	newArr := make([]uintptr, i+1)
	copy(newArr, arr)
	newArr[i] = value
	return newArr
}

func indirect(v reflect.Value) reflect.Value {

	// makes indirect more safe to call on values that are not Ptr or Interface
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Interface {
		return v
	}

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			//if e.Kind() == reflect.Ptr && !e.IsNil() {
			v = e
			continue
			//}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}
	return v
}

//todo: could be useful
//// unpackValue returns values inside of non-nil interfaces when possible.
//// This is useful for data types like structs, arrays, slices, and maps which
//// can contain varying types packed inside an interface.
//func (d *dumpState) unpackValue(v reflect.Value) reflect.Value {
//	if v.Kind() == reflect.Interface && !v.IsNil() {
//		v = v.Elem()
//	}
//	return v
//}
