package gotoken

import (
	"testing"
)

func TestToken(t *testing.T) {
	token, err := New("mosalut", 20, TOKEN_WEB) // 多端令牌模式

//	token, err := New("mosalut", 20) // 单端令牌模式
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(token)

	// 根据用户唯一标志（此处用户名，使用中是唯一字符串皆可）和端号，对应获取token
	token, err = GetCurrentToken("mosalut", TOKEN_PC) // 多端令牌模式
	if err != nil {
		t.Error(err.Error())
		return
	}

	if token == nil {
		t.Log("无此令牌")
		return
	}

	// 根据用户唯一标志（此处用户名，使用中是唯一字符串皆可）对应获取token
//	token, err = GetCurrentToken("mosalut" ) // 单端令牌模式

	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(token)

	ok := token.Validation(token.Code)
	t.Log(ok)
}
