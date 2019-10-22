package gotoken

import (
	"testing"
)

func TestToken(t *testing.T) {
//	token, err := New("mosalut", 20, TOKEN_PC) // 多端令牌模式
	token, err := New("mosalut", 20, TOKEN_SINGLE) // 单端令牌模式
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(token)

	/**
	 * 根据用户唯一标志（此处用户名，使用中是唯一字符串皆可）和端号，对应获取token
	 */
	/*
	token, err = GetCurrentToken("mosalut", TOKEN_PC) // 多端令牌模式
	if err != nil {
		t.Error(err.Error())
		return
	}

	if token == nil {
		t.Log("无此令牌")
		return
	}
	*/

	/**
	 * 根据用户唯一标志（此处用户名，使用中是唯一字符串皆可）对应获取token
	 */
	token, err = GetCurrentToken("mosalut", TOKEN_SINGLE) // 单端令牌模式
	if err != nil {
		t.Error(err.Error())
		return
	}

	if token == nil {
		t.Log("无此令牌")
		return
	}

	t.Log(token)

	ok := token.Validation(token.Code) // 模拟登陆后调用API时校验令牌
	t.Log(ok)

	token.Update("mosalut") // 如果需要，可以更新令牌

	t.Log(token)

	ok = token.Validation(token.Code)
	t.Log(ok)
}
