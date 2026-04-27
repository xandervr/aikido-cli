package output

import (
	"reflect"
	"testing"
)

func reflectTypeOf[T any]() reflect.Type {
	var zero T
	return reflect.TypeOf(zero)
}

func TestCollectColumns_PicksTaggedFieldsInOrder(t *testing.T) {
	type s struct {
		A string `aikido:"column,header=Alpha"`
		B int
		C string `aikido:"column,header=Charlie"`
	}
	cols := collectColumns(reflectTypeOf[s]())
	if len(cols) != 2 {
		t.Fatalf("expected 2 cols, got %d", len(cols))
	}
	if cols[0].Header != "Alpha" || cols[1].Header != "Charlie" {
		t.Fatalf("bad headers: %+v", cols)
	}
}
