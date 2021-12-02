package test

import (
	"reflect"
	"testing"

	"github.com/nbio/xml"
)

// Marshal validates if xml.Marshal(v) produces want or wantErr (if set).
func Marshal(t *testing.T, v interface{}, want string, wantErr bool) {
	x, err := xml.Marshal(v)
	if (err != nil) != wantErr {
		t.Errorf("xml.Marshal() error = %v, wantErr %v", err, wantErr)
		return
	}
	if string(x) != want {
		t.Errorf("xml.Marshal()\nGot:  %v\nWant: %v", string(x), want)
	}

	if v == nil {
		return
	}

	i := reflect.New(reflect.TypeOf(v).Elem()).Interface()
	err = xml.Unmarshal(x, i)
	if err != nil {
		t.Errorf("xml.Unmarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(v, v) {
		t.Errorf("xml.Unmarshal()\nGot:  %#v\nWant: %#v", i, v)
	}
}
