package dip

// Принцип инверсии зависимостей (Dependency Inversion Principle - DIP)
// Зависимости должны строиться относительно абстракций, а не деталей.

// Database - интерфейс для работы с базой данных
type Database interface {
	GetData() string
}

// MySQL - конкретная реализация Database
type MySQL struct{}

func (m MySQL) GetData() string {
	return "Данные из MySQL"
}

// PostgreSQL - другая реализация Database
type PostgreSQL struct{}

func (p PostgreSQL) GetData() string {
	return "Данные из PostgreSQL"
}

// Service зависит от абстракции Database, а не от конкретной реализации
type Service struct {
	db Database
}

func (s Service) GetData() string {
	return s.db.GetData()
}
