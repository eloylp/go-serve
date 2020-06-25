package mock

import (
	"github.com/stretchr/testify/mock"
)

type fakeLogger struct {
	mock.Mock
}

func NewFakeLogger() *fakeLogger { //nolint:golint
	return &fakeLogger{}
}

func (f *fakeLogger) Info(msg string) {
	f.Called(msg)
}

func (f *fakeLogger) Infof(msg string, args ...interface{}) {
	var called []interface{}
	called = append(called, msg)
	called = append(called, args...)
	f.Called(called...)
}

func (f *fakeLogger) Error(msg string) {
	f.Called(msg)
}

func (f *fakeLogger) Errorf(msg string, args ...interface{}) {
	var called []interface{}
	called = append(called, msg)
	called = append(called, args...)
	f.Called(called...)
}
