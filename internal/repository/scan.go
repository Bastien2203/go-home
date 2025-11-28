package repository

type Scanner interface {
	Scan(dest ...any) error
}
