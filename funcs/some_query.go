package funcs

type WorkerQuery interface {
	WorkerNameToPort(name string) (int, error)
	WorkerNameToUID(name string) (string, error)
}
