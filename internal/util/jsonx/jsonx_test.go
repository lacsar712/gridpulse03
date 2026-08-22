package jsonx

import "testing"

func TestDecodeExportAge(t *testing.T) {
	got, err := DecodeExport([]byte(`{"name":"alpha","age":9}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "alpha" || got.Age != 9 {
		t.Fatalf("got=%+v", got)
	}
}

func TestDecodeMapRejectsInvalid(t *testing.T) {
	_, err := DecodeMap([]byte(`{not-json`))
	if err == nil {
		t.Fatal("invalid json must return error")
	}
}

func TestDecodeMapOK(t *testing.T) {
	m, err := DecodeMap([]byte(`{"n":3}`))
	if err != nil {
		t.Fatal(err)
	}
	if m["n"] != 3 {
		t.Fatalf("%v", m)
	}
}
