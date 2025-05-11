package isp

// Принцип разделения интерфейса (Interface Segregation Principle - ISP)
// Клиенты не должны зависеть от интерфейсов, которые они не используют.

// Юольшой интерфейс, включающий в себя все возможные методы,
// которые используются в пакете.
type Artist interface {
	Dance()
	Sing()
	PlayGuitar()
}

// Минимальный интерфейс, требуемый для функции Dance.
type Dancer interface {
	Dance()
}

// Нарушение ISP, поскольку от аргумента требуется
// реализация методов, которые не используются.
func BadDance(dancer Artist) {
	dancer.Dance()
}

// Соблюдение ISP за счёт разделения интерфейсов
// для разных задач.
func GoodDance(dancer Dancer) {
	dancer.Dance()
}
