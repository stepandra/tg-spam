// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"io"
	"sync"
)

// SampleUpdaterMock is a mock implementation of lib.SampleUpdater.
//
//	func TestSomethingThatUsesSampleUpdater(t *testing.T) {
//
//		// make and configure a mocked lib.SampleUpdater
//		mockedSampleUpdater := &SampleUpdaterMock{
//			AppendFunc: func(msg string) error {
//				panic("mock out the Append method")
//			},
//			ReaderFunc: func() (io.ReadCloser, error) {
//				panic("mock out the Reader method")
//			},
//		}
//
//		// use mockedSampleUpdater in code that requires lib.SampleUpdater
//		// and then make assertions.
//
//	}
type SampleUpdaterMock struct {
	// AppendFunc mocks the Append method.
	AppendFunc func(msg string) error

	// ReaderFunc mocks the Reader method.
	ReaderFunc func() (io.ReadCloser, error)

	// calls tracks calls to the methods.
	calls struct {
		// Append holds details about calls to the Append method.
		Append []struct {
			// Msg is the msg argument value.
			Msg string
		}
		// Reader holds details about calls to the Reader method.
		Reader []struct {
		}
	}
	lockAppend sync.RWMutex
	lockReader sync.RWMutex
}

// Append calls AppendFunc.
func (mock *SampleUpdaterMock) Append(msg string) error {
	if mock.AppendFunc == nil {
		panic("SampleUpdaterMock.AppendFunc: method is nil but SampleUpdater.Append was just called")
	}
	callInfo := struct {
		Msg string
	}{
		Msg: msg,
	}
	mock.lockAppend.Lock()
	mock.calls.Append = append(mock.calls.Append, callInfo)
	mock.lockAppend.Unlock()
	return mock.AppendFunc(msg)
}

// AppendCalls gets all the calls that were made to Append.
// Check the length with:
//
//	len(mockedSampleUpdater.AppendCalls())
func (mock *SampleUpdaterMock) AppendCalls() []struct {
	Msg string
} {
	var calls []struct {
		Msg string
	}
	mock.lockAppend.RLock()
	calls = mock.calls.Append
	mock.lockAppend.RUnlock()
	return calls
}

// Reader calls ReaderFunc.
func (mock *SampleUpdaterMock) Reader() (io.ReadCloser, error) {
	if mock.ReaderFunc == nil {
		panic("SampleUpdaterMock.ReaderFunc: method is nil but SampleUpdater.Reader was just called")
	}
	callInfo := struct {
	}{}
	mock.lockReader.Lock()
	mock.calls.Reader = append(mock.calls.Reader, callInfo)
	mock.lockReader.Unlock()
	return mock.ReaderFunc()
}

// ReaderCalls gets all the calls that were made to Reader.
// Check the length with:
//
//	len(mockedSampleUpdater.ReaderCalls())
func (mock *SampleUpdaterMock) ReaderCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockReader.RLock()
	calls = mock.calls.Reader
	mock.lockReader.RUnlock()
	return calls
}
