package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

func main() {
	// Создаем клиент для подключения к Ollama
	addr, err := url.Parse("http://192.168.0.115:11434")
	if err != nil {
		log.Fatalf("Не удалось создать Ollama клиент: %v", err)
	}
	client := api.NewClient(addr, http.DefaultClient)

	// Указываем модель, которую будем использовать
	model := "qwen2.5:1.5b" // Можно изменить на любую другую установленную модель
	//model := "qwen2.5-coder:7b"

	ctx := context.Background()

	prompt := `
Напиши программу на Go, которая бы возвращала случайное имя 
человека в ответ на HHTP Get-запрос с URL 'http://localhost:8080/random-name'
`

	// Пример генерации (не в режиме чата)
	generateReq := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: &[]bool{true}[0],
	}

	fmt.Printf("\nГенерация с '%s': %s\n", model, generateReq.Prompt)
	generateRaspFunc := func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	}
	err = client.Generate(ctx, generateReq, generateRaspFunc)
	if err != nil {
		log.Fatalf("Ошибка генерации: %v", err)
	}
}
