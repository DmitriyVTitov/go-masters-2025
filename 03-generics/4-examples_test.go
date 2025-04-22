package generics

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

type TestType struct {
	ID   int
	Name string
}

func (*TestType) Validate() error {
	return nil
}

var v = TestType{
	ID:   100,
	Name: "Name",
}

func Test_decodeRequest(t *testing.T) {

	b, _ := json.Marshal(v)
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	got, err := decodeAndValidateRequest[*TestType](req)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &v) {
		t.Fatalf("decodeRequest() = %v, want %v", got, v)
	}
	req.Body.Close()
}

func TestQueue(t *testing.T) {
	q := New[TestType]()

	q.Enqueue(&v)
	got := q.Dequeue()
	if !reflect.DeepEqual(got, &v) {
		t.Fatalf("Dequeue() = %v, want %v", got, v)
	}
	if q.Len() != 0 {
		t.Fatalf("Len() = %v, want %v", q.Len(), 0)
	}
}

func TestFanIn(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	out := FanIn(ch1, ch2)

	go func() {
		ch1 <- 10
		ch2 <- 20
		close(ch1)
		close(ch2)
	}()

	for v := range out {
		t.Log(v)
	}
}
