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

type token struct {
	code [32]byte
	createTimeStamp int64
	validTime int64
	Code string `json:"code"`
}

func New(u string, vt int64, endian int) (*token, error) {
	var tErr error
	if singleMode {
		tErr = &tokenErr{"已经在单令牌模式运行"}
		return nil, tErr
	}

	tokenArray := [4]*token{}
	tokenArray[endian] = createToken(u, vt)
	tokens[u] = [4]*token{}
	tokens[u] = tokenArray

	return tokenArray[endian], nil
}

func NewSingle(u string, vt int64) (*token, error) {
	var tErr error
	if !singleMode {
		if len(tokens) != 0 {
			tErr = &tokenErr{"已经在多端令牌模式运行"}
			return nil, tErr
		}

		singleMode = true
	}

	tokens[u] = createToken(u, vt)

	return tokens[u].(*token), nil
}

func GetCurrentToken(u string, endian int) (*token, error) {
	ts, ok := tokens[u]

	if !ok {
		return nil, nil
	}

	var tErr error

	switch ts.(type) {
	case [4]*token:
		if endian == -1 {
			tErr = &tokenErr{"已经在多端令牌模式运行, 未传入端号"}
			return nil, tErr
		}
		return ts.([4]*token)[endian], nil
	case *token:
		if endian != -1 {
			tErr = &tokenErr{"已经在单令牌模式运行, 传入多余的端号"}
			return nil, tErr
		}
		return ts.(*token), nil
	default:
		if endian != -1 {
			tErr = &tokenErr{"令牌类型不正确, 需要[4]*token, 得到" + fmt.Sprintf("%v", reflect.TypeOf(ts))}
		}

		tErr = &tokenErr{"令牌类型不正确, 需要*token, 得到" + fmt.Sprintf("%v", reflect.TypeOf(ts))}

		return nil, tErr
	}
}

func createToken(u string, vt int64) *token {
	currentTimeStamp := time.Now().UnixNano()
	sum := sha256.Sum256([]byte(u + strconv.FormatInt(currentTimeStamp, 16)))
	t := &token{sum, currentTimeStamp, vt * 1000000, hex.EncodeToString(sum[:])}

	return t
}

func (t *token) GetCreateTimeStamp() int64 {
	return t.createTimeStamp
}

func (t *token)validation(u string) bool {
	if(sha256.Sum256([]byte(u + strconv.FormatInt(t.createTimeStamp, 16))) != t.code) {
		return false
	}

	if(time.Now().UnixNano() - t.createTimeStamp <= t.validTime) {
		return false
	}

	t.createTimeStamp = time.Now().UnixNano()
	return true
}

func (t *token)update(u string) {
	t.code = sha256.Sum256([]byte(u + strconv.FormatInt(t.createTimeStamp, 16)))
	t.Code = hex.EncodeToString(t.code[:])
}

func GetAll() map[string]interface{} {
	return tokens
}
