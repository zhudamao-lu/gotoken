package gotoken

const (
	TOKEN_SINGLE = -1
	TOKEN_WEB = 0
	TOKEN_APP = 1
	TOKEN_PC = 2
	TOKEN_OTHERS = 3
)

var (
	singleMode = false
	tokens = make(map[string]interface{})
)
