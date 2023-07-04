package jsonx

import (
	"encoding/json"
	"reflect"
)

func InterfaceToSchema(value interface{}) (schema string) {
	var construct func(v reflect.Value)

	construct = func(v reflect.Value) {
		vt := v.Type()

		if vt.Kind() == reflect.Ptr {
			vt = vt.Elem()
			v = v.Elem()
		}

		for i := 0; i < vt.NumField(); i++ {
			if f := v.Field(i); f.IsValid() {
				ft := vt.Field(i).Type

				switch ft.Kind() {
				case reflect.Ptr:
					fv := reflect.New(ft.Elem())
					construct(fv.Elem())
					f.Set(fv)

				case reflect.Slice:
					switch ft.Elem().Kind() {
					case reflect.Ptr:
						ot := ft.Elem().Elem()
						fv := reflect.New(ot)

						if ot.Kind() == reflect.Struct {
							construct(fv.Elem())
						}

						f.Set(reflect.Append(f, fv))

					default:
						ot := ft.Elem()

						switch ot.Kind() {
						case reflect.String:
							f.Set(reflect.Append(f, reflect.ValueOf(ot.String())))

						case reflect.Ptr:
							fv := reflect.New(ot)

							if ot.Elem().Kind() == reflect.Struct {
								construct(fv.Elem())
							}

							f.Set(reflect.Append(f, fv.Elem()))

						case reflect.Struct:
							fv := reflect.New(ot)
							construct(fv.Elem())
							f.Set(reflect.Append(f, fv.Elem()))

						default:
							f.Set(reflect.Append(f, reflect.New(ot).Elem()))
						}
					}

				case reflect.String:
					f.Set(reflect.ValueOf("string"))

				case reflect.Struct:
					construct(f)

				case reflect.Int32:
					f.Set(reflect.ValueOf(int32(1)))
				}
			}
		}
	}

	construct(reflect.ValueOf(value))

	if b, err := json.MarshalIndent(value, "", "    "); err == nil {
		schema = string(b)
	}

	return
}
