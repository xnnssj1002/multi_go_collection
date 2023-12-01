package multi_go_collection

type Tasker interface {
	exec() interface{}
	timeout() int
	getName() string
}

type HandleFunc func() interface{}

type Task struct {
	name string
	f    HandleFunc
	t    int
}

type OptionFunc func(task *Task)

func WithTaskName(name string) OptionFunc {
	return func(task *Task) {
		task.name = name
	}
}
func WithTimeout(t int) OptionFunc {
	return func(task *Task) {
		task.t = t
	}
}

func (t *Task) exec() interface{} {
	return t.f()
}

func (t *Task) timeout() int {
	return t.t
}

func (t *Task) getName() string {
	return t.name
}

func NewTask(f HandleFunc, ofs ...OptionFunc) *Task {
	task := &Task{
		f: f,
	}
	if len(ofs) > 0 {
		for _, of := range ofs {
			of(task)
		}
	}

	return task
}
