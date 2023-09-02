package libwx

import (
	"testing"
)

func Test_RelHumidity_type(t *testing.T) {
	var sut interface{} = RelHumidity(50)
	switch sut.(type) {
	case int:
		// ok
	case string:
		t.Fatal("RelHumidity should not be a string")
	default:
		t.Fatalf("RelHumidity is %T", sut)
	}
}
