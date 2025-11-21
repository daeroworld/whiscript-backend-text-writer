package repository

type IRepository[T any] interface {
	Upsert(e *T) (*T, error)
}
