package errcode

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type ErrCode struct {
	code int
	msg  string
}

type export struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e ErrCode) Code() int { return e.code }

func (e ErrCode) String() string { return strconv.Itoa(e.code) + " : " + e.msg }

func (e ErrCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(export{
		Code: e.code,
		Msg:  e.msg,
	})
}

func (e ErrCode) Response() (int, ErrCode) {
	return e.code / 100, e
}

func (e ErrCode) WithMsg(msg string) (int, ErrCode) {
	return e.code / 100, ErrCode{
		code: e.code,
		msg:  msg,
	}
}

func FromError(err error) (int, ErrCode) {
	s := reflect.TypeOf(err).String()
	ec := ServerCommonError
	switch s {
	case "validator.ValidationErrors":
		ec = InvalidParam
	case "proto.RedisError":
		ec = ServerRedisError
	case "nsq.ErrProtocol":
		ec = ServerNsqError
	case "*url.Error":
		ec = ResponseTimeout
	case "*json.SyntaxError", "*json.UnmarshalTypeError":
		ec = ResponseWrong
	}
	return ec.Response()
}
