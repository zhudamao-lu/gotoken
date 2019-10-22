package gotoken

type tokenErr struct {
	message string
}

func (tE *tokenErr) Error() string {
	return tE.message
}

func (tE *tokenErr) setMessage(msg string) {
	tE.message = msg
}
