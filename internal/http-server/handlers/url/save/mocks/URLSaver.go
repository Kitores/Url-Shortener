package mocks

import (
	"github.com/stretchr/testify/mock"
)

type URLSaver struct {
	mock.Mock
}

func (m *URLSaver) SaveURL(urlToSave string, alias string) error {
	var resultErr error
	ret := m.Called(urlToSave, alias)
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		resultErr = rf(urlToSave, alias)
	} else {
		resultErr = ret.Error(0)
	}
	return resultErr
}

type mockConstructorTestingTNewURLSaver interface {
	mock.TestingT
	Cleanup(func())
}

func NewURLSaver(t mockConstructorTestingTNewURLSaver) *URLSaver {
	mok := &URLSaver{}
	mok.Mock.Test(t)

	//t.Cleanup(func() { mok.AssertExpectations(t) })

	return mok
}
