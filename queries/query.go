package queries

type Query interface {
	HashCode() (uint64, error)
}
