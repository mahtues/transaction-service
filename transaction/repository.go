package transaction

type Repository interface {
	SaveTransaction() error
	LoadTransaction() error
}
