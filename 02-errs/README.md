# Ошибки в Go: понятие и работа с ними

В Go ошибки (errors) - это стандартный способ обработки исключительных ситуаций в программе. В отличие от многих других языков, в Go нет механизма исключений (try-catch). Вместо этого функции, которые могут завершиться с ошибкой, возвращают её как последнее значение.

## Основные концепции ошибок в Go

1. **Ошибки - это значения**: В Go ошибки представлены интерфейсом `error`:
   ```go
   type error interface {
       Error() string
   }
   ```

2. **Проверка ошибок**: После вызова функции, которая может вернуть ошибку, нужно явно проверить её:
   ```go
   result, err := someFunction()
   if err != nil {
       // Обработка ошибки
   }
   ```

## Создание ошибок

1. **Использование `errors.New`**:
   ```go
   import "errors"
   
   func divide(a, b float64) (float64, error) {
       if b == 0 {
           return 0, errors.New("division by zero")
       }
       return a / b, nil
   }
   ```

2. **Использование `fmt.Errorf`** (с форматированием):
   ```go
   if age < 0 {
       return fmt.Errorf("invalid age %d: age cannot be negative", age)
   }
   ```

## Проверка ошибок

Основной способ проверки:
```go
file, err := os.Open("filename.txt")
if err != nil {
    log.Fatal(err) // или другая обработка
}
// Продолжаем работу с file
```

## Расширенная работа с ошибками

1. **Кастомные типы ошибок**:
   ```go
   type MyError struct {
       Code    int
       Message string
   }
   
   func (e *MyError) Error() string {
       return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
   }
   
   func someFunc() error {
       return &MyError{Code: 404, Message: "Not Found"}
   }
   ```

2. **Проверка типа ошибки**:
   ```go
   err := someFunc()
   if myErr, ok := err.(*MyError); ok {
       // Обработка MyError
       fmt.Println("Code:", myErr.Code)
   }
   ```

3. **Обертывание ошибок (Go 1.13+)**:
   ```go
   if err != nil {
       return fmt.Errorf("context: %w", err)
   }
   
   // Проверка обернутой ошибки
   if errors.Is(err, os.ErrNotExist) {
       // Обработка случая, когда файл не существует
   }
   
   var pathError *os.PathError
   if errors.As(err, &pathError) {
       // Обработка PathError
   }
   ```

## Лучшие практики

1. Всегда проверяйте ошибки, не игнорируйте их
2. Добавляйте контекст к ошибкам при передаче вверх по стеку вызовов
3. Используйте `errors.Is` и `errors.As` вместо прямых сравнений
4. Для простых ошибок достаточно `errors.New` или `fmt.Errorf`
5. Для сложных сценариев создавайте собственные типы ошибок

## Пример комплексной обработки

```go
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()
    
    // Другие операции с файлом
    
    return nil
}

func main() {
    err := processFile("data.txt")
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("File does not exist")
        } else {
            fmt.Printf("Unexpected error: %v\n", err)
        }
        os.Exit(1)
    }
}
```

Ошибки в Go - это мощный и гибкий механизм, который при правильном использовании делает код более надежным и понятным.