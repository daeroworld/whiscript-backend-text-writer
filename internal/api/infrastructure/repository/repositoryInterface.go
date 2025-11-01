package repository

type IRepository[T any] interface {
	Save(e *T) (*T, error)
}
