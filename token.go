package gotoken

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"reflect"
	"time"
	"fmt"
)

type hashCode [32]byte

func (hC hashCode) String() string {
	return string(hC[:])
}

// 令牌结构
type token struct {
	code [32]byte // 通过唯一标志和当前时间生成的哈希码
	createTimestamp int64 // 创建时间和更新时间戳
	validTime int64 // 有效时间戳
	Code string `json:"code"` // 哈希码16进制字符串的表示
}

// 创建一个令牌，如果参数endian = -1，则表示单端令牌模式，否则为多端
func New(u string, vt int64, endian int) (*token, error) {
	if endian < -1 || endian > 3 {
		return nil, &tokenErr{"端号超出范围, [-1, 3]"}
	}

	if endian == -1 {
		token, err := newSingle(u, vt)
		if err != nil {
			return nil, err
		}

		return token, nil
	}

	if singleMode {
		return nil, &tokenErr{"已经在单令牌模式运行"}
	}

	tokenArray := [4]*token{}
	tokenArray[endian] = createToken(u, vt)
	tokens[u] = [4]*token{}
	tokens[u] = tokenArray

	return tokenArray[endian], nil
}

// 创建一个单端模式令牌
func newSingle(u string, vt int64) (*token, error) {
	if !singleMode {
		if len(tokens) != 0 {
			return nil, &tokenErr{"已经在多端令牌模式运行"}
		}

		singleMode = true
	}

	tokens[u] = createToken(u, vt)

	return tokens[u].(*token), nil
}

// 获取当前令牌，如果参数endian = -1，则表示单端令牌模式，否则为多端
func GetCurrentToken(u string, endian int) (*token, error) {
	ts, ok := tokens[u]

	if !ok {
		return nil, nil
	}

	switch ts.(type) {
	case [4]*token:
		if endian == -1 {
			return nil, &tokenErr{"已经在多端令牌模式运行"}
		}
		return ts.([4]*token)[endian], nil
	case *token:
		if endian != -1 {
			return nil, &tokenErr{"已经在单令牌模式运行"}
		}
		return ts.(*token), nil
	default:
		var tErr error
		if endian != -1 {
			tErr = &tokenErr{"令牌类型不正确, 需要[4]*token, 得到" + fmt.Sprintf("%v", reflect.TypeOf(ts))}
		}

		tErr = &tokenErr{"令牌类型不正确, 需要*token, 得到" + fmt.Sprintf("%v", reflect.TypeOf(ts))}

		return nil, tErr
	}
}

// 创建令牌
func createToken(u string, vt int64) *token {
	currentTimeStamp := time.Now().UnixNano()
	sum := sha256.Sum256([]byte(u + strconv.FormatInt(currentTimeStamp, 16)))
	t := &token{sum, currentTimeStamp, vt * 1000000000, hex.EncodeToString(sum[:])}

	return t
}

// 获取令牌创建时间 暂时没啥用
func (t *token) GetCreateTimeStamp() int64 {
	return t.createTimestamp
}

// 令牌校验
func (t *token)Validation(code string) bool {
	// 判断令牌的哈希码的16进制字符串表示和传来的code是否相等
	if code != t.Code {
		return false
	}

	// 判断令牌是否超时
	if(time.Now().UnixNano() - t.createTimestamp > t.validTime) {
		return false
	}

	t.createTimestamp = time.Now().UnixNano()
	return true
}

func (t *token)Update(u string) {
	t.code = sha256.Sum256([]byte(u + strconv.FormatInt(t.createTimestamp, 16)))
	t.Code = hex.EncodeToString(t.code[:])
}

func GetAll() map[string]interface{} {
	return tokens
}
