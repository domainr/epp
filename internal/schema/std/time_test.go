package std

import (
	"testing"
	"time"

	"github.com/domainr/epp/internal/schema/test"
)

func TestTime(t *testing.T) {
	may19, err := time.Parse(time.RFC3339, "2015-05-19T06:34:21.1Z")
	if err != nil {
		t.Fatal(err)
	}

	type T struct {
		XMLName struct{} `xml:"example"`
		Value   *Time    `xml:"when"`
		Attr    *Time    `xml:"when,attr,omitempty"`
	}

	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			`no tags`,
			&T{},
			`<example></example>`,
			false,
		},
		{
			`zero value chardata`,
			&T{Value: &Time{}},
			`<example><when>0001-01-01T00:00:00Z</when></example>`,
			false,
		},
		{
			`zero value attr`,
			&T{Attr: &Time{}},
			`<example when="0001-01-01T00:00:00Z"></example>`,
			false,
		},
		{
			`chardata`,
			&T{Value: &Time{may19}},
			`<example><when>2015-05-19T06:34:21.1Z</when></example>`,
			false,
		},
		{
			`attr`,
			&T{Attr: &Time{may19}},
			`<example when="2015-05-19T06:34:21.1Z"></example>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Marshal(t, tt.v, tt.want, tt.wantErr)
		})
	}
}
