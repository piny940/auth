package domain

type ErrRecordNotFound struct{}

func (e ErrRecordNotFound) Error() string {
	return "record not found"
}
