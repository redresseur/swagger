package wrapper

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSmartPoint(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Ok bool `json:"ok"`
	}

	test := person{Name: "author", Ok: true,}

	data, _ := json.Marshal(&test)

	//pprint := func( p *person){
	//	t.Logf("%s %v", p.Name, p.Ok)
	//}

	SmartPoint(data, reflect.TypeOf(&test))
	//p := &person{}
	//
	//unsafe.Pointer(SmartPoint(data, reflect.TypeOf(test)))
	//pprint()
}