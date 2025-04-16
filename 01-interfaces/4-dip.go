package interfaces

// Пример инверсии зависимостей.
//
// SAGA - микросервисный шаблон проектирования, позволяющий
// выполнять распределенные транзакции.
// https://microservices.io/patterns/data/saga.html
type Saga struct {
	state State
	steps []Step
}

type State int

const (
	StateNew State = iota
	StateInProgress
	StateCompleted
	StateFailed
)

// Шаг процесса - интерфейс, являющийся частью бизнес-логики.
// Вся реализация построена на абстрактном "шаге" и не привязана
// к конкретной имплементации.
// Таким образом нигде код бизнес логики не зависит от кода имплементации
// интерфейса `Step`
type Step interface {
	Name() string
	Execute() error
	Undo() error
}

func New() *Saga {
	return &Saga{
		state: StateNew,
		steps: make([]Step, 0),
	}
}

func (s *Saga) AddStep(step Step) {
	s.steps = append(s.steps, step)
}

func (s *Saga) Run() error {
	s.state = StateInProgress
	for _, step := range s.steps {
		if err := step.Execute(); err != nil {
			// В случае, если на каком-либо этапе произошла ошибка,
			// выполняем "компенсирующие" транзакции.
			s.state = StateFailed
			for _, undoStep := range s.steps {
				if undoStep.Name() == step.Name() {
					break
				}
				if err := undoStep.Undo(); err != nil {
					return err
				}
			}
			return err
		}
	}
	s.state = StateCompleted
	return nil
}

func (s *Saga) GetState() State {
	return s.state
}
