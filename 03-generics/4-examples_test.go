package generics

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

type T struct {
	ID   int
	Name string
}

var v = T{
	ID:   100,
	Name: "Name",
}

func Test_decodeRequest(t *testing.T) {

	b, _ := json.Marshal(v)
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	got, err := decodeRequest[T](req)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &v) {
		t.Fatalf("decodeRequest() = %v, want %v", got, v)
	}
}

func TestQueue(t *testing.T) {
	q := New[T]()

	q.Enqueue(&v)
	got := q.Dequeue()
	if !reflect.DeepEqual(got, &v) {
		t.Fatalf("Dequeue() = %v, want %v", got, v)
	}
	if q.Len() != 0 {
		t.Fatalf("Len() = %v, want %v", q.Len(), 0)
	}
}
