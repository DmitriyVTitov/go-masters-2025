package main

import (
	"log"
	"math"
	"net/http"
	_ "net/http/pprof" // автоматически регистрирует обработчики в DefaultServeMux
	"sync"
)

// Вызов профиля:
// go tool pprof -web http://localhost:6060/debug/pprof/profile?seconds=5
// go tool pprof -web http://localhost:6060/debug/pprof/heap

func main() {
	// http://localhost:6060/cpu
	// http://localhost:6060/mem
	http.Handle("/cpu", http.HandlerFunc(cpuHandler))
	http.Handle("/mem", http.HandlerFunc(memHandler))

	// Запускаем HTTP сервер с pprof
	// http://localhost:6060/debug/pprof/
	log.Fatal(http.ListenAndServe("localhost:6060", nil))
}

func cpuHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	var wg sync.WaitGroup
	wg.Add(4)

	f := func() {
		for i := range 10_000_000 {
			_ = math.Sin(float64(i))
		}

		wg.Done()
	}

	for range 4 {
		go f()
	}

	wg.Wait()
	w.Write([]byte(r.URL.Path))
}

func memHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	buf := make([]byte, 1_000_000_000)
	for i := range len(buf) {
		buf[i] = byte(i % 256)
	}

	for i := range len(buf) {
		_ = buf[i]
	}

	w.Write([]byte(r.URL.Path))
}
