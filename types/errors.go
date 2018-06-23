package types

import (
	"encoding/json"

	"github.com/go-errors/errors"
)

var MaxStackDepth = 50

type GameOverException struct {
	*BaseException
}

type BusinessException struct {
	*BaseException
}

type BaseException struct {
	message string
	code    string
	base    *errors.Error
}

type ExceptionMessage struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Stack   string `json:"stack"`
}

func NewGameOverException(message string, code string) *GameOverException {

	a := &GameOverException{BaseException: &BaseException{}}
	x := errors.New(a)

	e := &GameOverException{BaseException: &BaseException{}}
	e.message = message
	e.code = code
	e.base = x

	return e
}

func NewBusinessException(message string, code string) *BusinessException {
	a := &BusinessException{BaseException: &BaseException{}}
	x := errors.New(a)

	e := &BusinessException{BaseException: &BaseException{}}
	e.message = message
	e.code = code
	e.base = x

	return e
}

func (e *GameOverException) Error() string {
	exp := &ExceptionMessage{Code: e.code, Message: e.message, Stack: string(e.base.Stack())}
	expBytes, _ := json.Marshal(exp)
	return string(expBytes)
}

func (e *BusinessException) Error() string {
	exp := &ExceptionMessage{Code: e.code, Message: e.message, Stack: string(e.base.Stack())}
	expBytes, _ := json.Marshal(exp)
	return string(expBytes)
}
