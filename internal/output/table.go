package output

import (
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
)

type colDef struct {
	Header string
	Index  []int
}

func (r *Renderer) renderTable(v any) (handled bool, err error) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	var elemType reflect.Type
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		elemType = rv.Type().Elem()
	case reflect.Struct:
		elemType = rv.Type()
	default:
		return false, nil
	}

	cols := collectColumns(elemType)
	if len(cols) == 0 {
		return false, nil
	}

	tw := tabwriter.NewWriter(r.Out, 0, 0, 2, ' ', 0)
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.Header
	}
	fmt.Fprintln(tw, strings.Join(headers, "\t"))

	writeRow := func(elem reflect.Value) {
		fields := make([]string, len(cols))
		for i, c := range cols {
			fv := elem.FieldByIndex(c.Index)
			fields[i] = formatField(fv)
		}
		fmt.Fprintln(tw, strings.Join(fields, "\t"))
	}

	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		for i := 0; i < rv.Len(); i++ {
			writeRow(rv.Index(i))
		}
	} else {
		writeRow(rv)
	}

	return true, tw.Flush()
}

func collectColumns(t reflect.Type) []colDef {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	var cols []colDef
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("aikido")
		if tag == "" {
			continue
		}
		parts := strings.Split(tag, ",")
		if parts[0] != "column" {
			continue
		}
		header := f.Name
		for _, p := range parts[1:] {
			if strings.HasPrefix(p, "header=") {
				header = strings.TrimPrefix(p, "header=")
			}
		}
		cols = append(cols, colDef{Header: header, Index: f.Index})
	}
	return cols
}

func formatField(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
