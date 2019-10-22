package gotoken

import (
	"testing"
)

func TestToken(t *testing.T) {
	token, err := New("", 20, TOKEN_WEB)
//	token, err := NewSingle("", 20)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(token)
	t.Log(token.GetCodeString())
}
