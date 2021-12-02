package std

import (
	"reflect"
	"testing"

	"github.com/nbio/xml"
)

func TestBool(t *testing.T) {
	type T1 struct {
		XMLName struct{} `xml:"example"`
		Fred    Bool     `xml:"fred"`
		Jane    Bool     `xml:"jane"`
		Susan   Bool     `xml:"susan"`
	}

	type T2 struct {
		XMLName struct{} `xml:"example,selfclosing"`
		Fred    Bool     `xml:"fred,attr"`
		Jane    Bool     `xml:"jane,attr,omitempty"`
		Susan   Bool     `xml:"susan,attr,omitempty"`
	}

	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			`nil`,
			nil,
			``,
			false,
		},
		{
			`no tags`,
			&T1{},
			`<example></example>`,
			false,
		},
		{
			`Fred`,
			&T1{Fred: true},
			`<example><fred/></example>`,
			false,
		},
		{
			`Jane`,
			&T1{Jane: true},
			`<example><jane/></example>`,
			false,
		},
		{
			`Fred and Susan`,
			&T1{Fred: true, Susan: true},
			`<example><fred/><susan/></example>`,
			false,
		},
		{
			`Fred attribute`,
			&T2{Fred: true},
			`<example fred="1"></example>`,
			false,
		},
		{
			`Jane attribute`,
			&T2{Jane: true},
			`<example fred="0" jane="1"></example>`,
			false,
		},
		{
			`Fred and Susan attributes`,
			&T2{Fred: true, Susan: true},
			`<example fred="1" susan="1"></example>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, err := xml.Marshal(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("xml.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(x) != tt.want {
				t.Errorf("xml.Marshal()\nGot:  %v\nWant: %v", string(x), tt.want)
			}

			if tt.v == nil {
				return
			}

			v := reflect.New(reflect.TypeOf(tt.v).Elem()).Interface()
			err = xml.Unmarshal(x, v)
			if err != nil {
				t.Errorf("xml.Unmarshal() error = %v", err)
				return
			}
			if !reflect.DeepEqual(v, tt.v) {
				t.Errorf("xml.Unmarshal()\nGot:  %#v\nWant: %#v", v, v)
			}
		})
	}
}
