package wrapper

import (
	"encoding/json"
	"reflect"
)

func SmartPoint(data []byte, t reflect.Type)interface{}{
	v := reflect.New(t)

	if err := json.Unmarshal(data, &v); err != nil{
		panic(err)
	}
	return v
}