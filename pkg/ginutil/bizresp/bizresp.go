package bizresp

import (
	"reflect"
	"strconv"
)

type errcode struct {
	code int
	msg  string
}

func (e errcode) String() string {
	return strconv.Itoa(e.code) + " : " + e.msg
}

type ErrCode struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e errcode) Reply() (int, ErrCode) {
	return e.code / 100, ErrCode{
		Code: e.code,
		Msg:  e.msg,
	}
}

func (e errcode) WithMsg(msg string) (int, ErrCode) {
	return e.code / 100, ErrCode{
		Code: e.code,
		Msg:  msg,
	}
}

func WithErr(err error) (int, ErrCode) {
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
	return ec.Reply()
}
