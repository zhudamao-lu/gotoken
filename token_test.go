package gotoken

import (
	"testing"
)

func TestToken(t *testing.T) {
	token, err := New("mosalut", 20, TOKEN_WEB)

//	token, err := NewSingle("mosalut", 20)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(token)

	token, err = GetCurrentToken("mosalut", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(token)
}
