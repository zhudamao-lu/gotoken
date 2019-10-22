package gotoken

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type hashCode [32]byte

func (hC hashCode) String() string {
	return string(hC[:])
}

type token struct {
	Code [32]byte `json:"code"`
	createTimeStamp int64
	validTime int64
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

func createToken(u string, vt int64) *token {
	currentTimeStamp := time.Now().UnixNano()
	sum := sha256.Sum256([]byte(u + strconv.FormatInt(currentTimeStamp, 16)))
	t := &token{sum, currentTimeStamp, vt * 1000000}

	return t
}

func (t *token) GetCodeString() string {
	return hex.EncodeToString(t.Code[:])
}

func (t *token) GetCreateTimeStamp() int64 {
	return t.createTimeStamp
}

func (t *token)validation(u string) bool {
	if(sha256.Sum256([]byte(u + strconv.FormatInt(t.createTimeStamp, 16))) != t.Code) {
		return false
	}

	if(time.Now().UnixNano() - t.createTimeStamp <= t.validTime) {
		return false
	}

	return true
}

func GetAll() map[string]interface{} {
	return tokens
}

