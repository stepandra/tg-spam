// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"github.com/umputun/tg-spam/lib"
	"io"
	"sync"
)

// DetectorMock is a mock implementation of bot.Detector.
//
//	func TestSomethingThatUsesDetector(t *testing.T) {
//
//		// make and configure a mocked bot.Detector
//		mockedDetector := &DetectorMock{
//			AddApprovedUsersFunc: func(ids ...string)  {
//				panic("mock out the AddApprovedUsers method")
//			},
//			ApprovedUsersFunc: func() []string {
//				panic("mock out the ApprovedUsers method")
//			},
//			CheckFunc: func(msg string, userID string) (bool, []lib.CheckResult) {
//				panic("mock out the Check method")
//			},
//			LoadSamplesFunc: func(exclReader io.Reader, spamReaders []io.Reader, hamReaders []io.Reader) (lib.LoadResult, error) {
//				panic("mock out the LoadSamples method")
//			},
//			LoadStopWordsFunc: func(readers ...io.Reader) (lib.LoadResult, error) {
//				panic("mock out the LoadStopWords method")
//			},
//			RemoveApprovedUsersFunc: func(ids ...string)  {
//				panic("mock out the RemoveApprovedUsers method")
//			},
//			UpdateHamFunc: func(msg string) error {
//				panic("mock out the UpdateHam method")
//			},
//			UpdateSpamFunc: func(msg string) error {
//				panic("mock out the UpdateSpam method")
//			},
//		}
//
//		// use mockedDetector in code that requires bot.Detector
//		// and then make assertions.
//
//	}
type DetectorMock struct {
	// AddApprovedUsersFunc mocks the AddApprovedUsers method.
	AddApprovedUsersFunc func(ids ...string)

	// ApprovedUsersFunc mocks the ApprovedUsers method.
	ApprovedUsersFunc func() []string

	// CheckFunc mocks the Check method.
	CheckFunc func(msg string, userID string) (bool, []lib.CheckResult)

	// LoadSamplesFunc mocks the LoadSamples method.
	LoadSamplesFunc func(exclReader io.Reader, spamReaders []io.Reader, hamReaders []io.Reader) (lib.LoadResult, error)

	// LoadStopWordsFunc mocks the LoadStopWords method.
	LoadStopWordsFunc func(readers ...io.Reader) (lib.LoadResult, error)

	// RemoveApprovedUsersFunc mocks the RemoveApprovedUsers method.
	RemoveApprovedUsersFunc func(ids ...string)

	// UpdateHamFunc mocks the UpdateHam method.
	UpdateHamFunc func(msg string) error

	// UpdateSpamFunc mocks the UpdateSpam method.
	UpdateSpamFunc func(msg string) error

	// calls tracks calls to the methods.
	calls struct {
		// AddApprovedUsers holds details about calls to the AddApprovedUsers method.
		AddApprovedUsers []struct {
			// Ids is the ids argument value.
			Ids []string
		}
		// ApprovedUsers holds details about calls to the ApprovedUsers method.
		ApprovedUsers []struct {
		}
		// Check holds details about calls to the Check method.
		Check []struct {
			// Msg is the msg argument value.
			Msg string
			// UserID is the userID argument value.
			UserID string
		}
		// LoadSamples holds details about calls to the LoadSamples method.
		LoadSamples []struct {
			// ExclReader is the exclReader argument value.
			ExclReader io.Reader
			// SpamReaders is the spamReaders argument value.
			SpamReaders []io.Reader
			// HamReaders is the hamReaders argument value.
			HamReaders []io.Reader
		}
		// LoadStopWords holds details about calls to the LoadStopWords method.
		LoadStopWords []struct {
			// Readers is the readers argument value.
			Readers []io.Reader
		}
		// RemoveApprovedUsers holds details about calls to the RemoveApprovedUsers method.
		RemoveApprovedUsers []struct {
			// Ids is the ids argument value.
			Ids []string
		}
		// UpdateHam holds details about calls to the UpdateHam method.
		UpdateHam []struct {
			// Msg is the msg argument value.
			Msg string
		}
		// UpdateSpam holds details about calls to the UpdateSpam method.
		UpdateSpam []struct {
			// Msg is the msg argument value.
			Msg string
		}
	}
	lockAddApprovedUsers    sync.RWMutex
	lockApprovedUsers       sync.RWMutex
	lockCheck               sync.RWMutex
	lockLoadSamples         sync.RWMutex
	lockLoadStopWords       sync.RWMutex
	lockRemoveApprovedUsers sync.RWMutex
	lockUpdateHam           sync.RWMutex
	lockUpdateSpam          sync.RWMutex
}

// AddApprovedUsers calls AddApprovedUsersFunc.
func (mock *DetectorMock) AddApprovedUsers(ids ...string) {
	if mock.AddApprovedUsersFunc == nil {
		panic("DetectorMock.AddApprovedUsersFunc: method is nil but Detector.AddApprovedUsers was just called")
	}
	callInfo := struct {
		Ids []string
	}{
		Ids: ids,
	}
	mock.lockAddApprovedUsers.Lock()
	mock.calls.AddApprovedUsers = append(mock.calls.AddApprovedUsers, callInfo)
	mock.lockAddApprovedUsers.Unlock()
	mock.AddApprovedUsersFunc(ids...)
}

// AddApprovedUsersCalls gets all the calls that were made to AddApprovedUsers.
// Check the length with:
//
//	len(mockedDetector.AddApprovedUsersCalls())
func (mock *DetectorMock) AddApprovedUsersCalls() []struct {
	Ids []string
} {
	var calls []struct {
		Ids []string
	}
	mock.lockAddApprovedUsers.RLock()
	calls = mock.calls.AddApprovedUsers
	mock.lockAddApprovedUsers.RUnlock()
	return calls
}

// ResetAddApprovedUsersCalls reset all the calls that were made to AddApprovedUsers.
func (mock *DetectorMock) ResetAddApprovedUsersCalls() {
	mock.lockAddApprovedUsers.Lock()
	mock.calls.AddApprovedUsers = nil
	mock.lockAddApprovedUsers.Unlock()
}

// ApprovedUsers calls ApprovedUsersFunc.
func (mock *DetectorMock) ApprovedUsers() []string {
	if mock.ApprovedUsersFunc == nil {
		panic("DetectorMock.ApprovedUsersFunc: method is nil but Detector.ApprovedUsers was just called")
	}
	callInfo := struct {
	}{}
	mock.lockApprovedUsers.Lock()
	mock.calls.ApprovedUsers = append(mock.calls.ApprovedUsers, callInfo)
	mock.lockApprovedUsers.Unlock()
	return mock.ApprovedUsersFunc()
}

// ApprovedUsersCalls gets all the calls that were made to ApprovedUsers.
// Check the length with:
//
//	len(mockedDetector.ApprovedUsersCalls())
func (mock *DetectorMock) ApprovedUsersCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockApprovedUsers.RLock()
	calls = mock.calls.ApprovedUsers
	mock.lockApprovedUsers.RUnlock()
	return calls
}

// ResetApprovedUsersCalls reset all the calls that were made to ApprovedUsers.
func (mock *DetectorMock) ResetApprovedUsersCalls() {
	mock.lockApprovedUsers.Lock()
	mock.calls.ApprovedUsers = nil
	mock.lockApprovedUsers.Unlock()
}

// Check calls CheckFunc.
func (mock *DetectorMock) Check(msg string, userID string) (bool, []lib.CheckResult) {
	if mock.CheckFunc == nil {
		panic("DetectorMock.CheckFunc: method is nil but Detector.Check was just called")
	}
	callInfo := struct {
		Msg    string
		UserID string
	}{
		Msg:    msg,
		UserID: userID,
	}
	mock.lockCheck.Lock()
	mock.calls.Check = append(mock.calls.Check, callInfo)
	mock.lockCheck.Unlock()
	return mock.CheckFunc(msg, userID)
}

// CheckCalls gets all the calls that were made to Check.
// Check the length with:
//
//	len(mockedDetector.CheckCalls())
func (mock *DetectorMock) CheckCalls() []struct {
	Msg    string
	UserID string
} {
	var calls []struct {
		Msg    string
		UserID string
	}
	mock.lockCheck.RLock()
	calls = mock.calls.Check
	mock.lockCheck.RUnlock()
	return calls
}

// ResetCheckCalls reset all the calls that were made to Check.
func (mock *DetectorMock) ResetCheckCalls() {
	mock.lockCheck.Lock()
	mock.calls.Check = nil
	mock.lockCheck.Unlock()
}

// LoadSamples calls LoadSamplesFunc.
func (mock *DetectorMock) LoadSamples(exclReader io.Reader, spamReaders []io.Reader, hamReaders []io.Reader) (lib.LoadResult, error) {
	if mock.LoadSamplesFunc == nil {
		panic("DetectorMock.LoadSamplesFunc: method is nil but Detector.LoadSamples was just called")
	}
	callInfo := struct {
		ExclReader  io.Reader
		SpamReaders []io.Reader
		HamReaders  []io.Reader
	}{
		ExclReader:  exclReader,
		SpamReaders: spamReaders,
		HamReaders:  hamReaders,
	}
	mock.lockLoadSamples.Lock()
	mock.calls.LoadSamples = append(mock.calls.LoadSamples, callInfo)
	mock.lockLoadSamples.Unlock()
	return mock.LoadSamplesFunc(exclReader, spamReaders, hamReaders)
}

// LoadSamplesCalls gets all the calls that were made to LoadSamples.
// Check the length with:
//
//	len(mockedDetector.LoadSamplesCalls())
func (mock *DetectorMock) LoadSamplesCalls() []struct {
	ExclReader  io.Reader
	SpamReaders []io.Reader
	HamReaders  []io.Reader
} {
	var calls []struct {
		ExclReader  io.Reader
		SpamReaders []io.Reader
		HamReaders  []io.Reader
	}
	mock.lockLoadSamples.RLock()
	calls = mock.calls.LoadSamples
	mock.lockLoadSamples.RUnlock()
	return calls
}

// ResetLoadSamplesCalls reset all the calls that were made to LoadSamples.
func (mock *DetectorMock) ResetLoadSamplesCalls() {
	mock.lockLoadSamples.Lock()
	mock.calls.LoadSamples = nil
	mock.lockLoadSamples.Unlock()
}

// LoadStopWords calls LoadStopWordsFunc.
func (mock *DetectorMock) LoadStopWords(readers ...io.Reader) (lib.LoadResult, error) {
	if mock.LoadStopWordsFunc == nil {
		panic("DetectorMock.LoadStopWordsFunc: method is nil but Detector.LoadStopWords was just called")
	}
	callInfo := struct {
		Readers []io.Reader
	}{
		Readers: readers,
	}
	mock.lockLoadStopWords.Lock()
	mock.calls.LoadStopWords = append(mock.calls.LoadStopWords, callInfo)
	mock.lockLoadStopWords.Unlock()
	return mock.LoadStopWordsFunc(readers...)
}

// LoadStopWordsCalls gets all the calls that were made to LoadStopWords.
// Check the length with:
//
//	len(mockedDetector.LoadStopWordsCalls())
func (mock *DetectorMock) LoadStopWordsCalls() []struct {
	Readers []io.Reader
} {
	var calls []struct {
		Readers []io.Reader
	}
	mock.lockLoadStopWords.RLock()
	calls = mock.calls.LoadStopWords
	mock.lockLoadStopWords.RUnlock()
	return calls
}

// ResetLoadStopWordsCalls reset all the calls that were made to LoadStopWords.
func (mock *DetectorMock) ResetLoadStopWordsCalls() {
	mock.lockLoadStopWords.Lock()
	mock.calls.LoadStopWords = nil
	mock.lockLoadStopWords.Unlock()
}

// RemoveApprovedUsers calls RemoveApprovedUsersFunc.
func (mock *DetectorMock) RemoveApprovedUsers(ids ...string) {
	if mock.RemoveApprovedUsersFunc == nil {
		panic("DetectorMock.RemoveApprovedUsersFunc: method is nil but Detector.RemoveApprovedUsers was just called")
	}
	callInfo := struct {
		Ids []string
	}{
		Ids: ids,
	}
	mock.lockRemoveApprovedUsers.Lock()
	mock.calls.RemoveApprovedUsers = append(mock.calls.RemoveApprovedUsers, callInfo)
	mock.lockRemoveApprovedUsers.Unlock()
	mock.RemoveApprovedUsersFunc(ids...)
}

// RemoveApprovedUsersCalls gets all the calls that were made to RemoveApprovedUsers.
// Check the length with:
//
//	len(mockedDetector.RemoveApprovedUsersCalls())
func (mock *DetectorMock) RemoveApprovedUsersCalls() []struct {
	Ids []string
} {
	var calls []struct {
		Ids []string
	}
	mock.lockRemoveApprovedUsers.RLock()
	calls = mock.calls.RemoveApprovedUsers
	mock.lockRemoveApprovedUsers.RUnlock()
	return calls
}

// ResetRemoveApprovedUsersCalls reset all the calls that were made to RemoveApprovedUsers.
func (mock *DetectorMock) ResetRemoveApprovedUsersCalls() {
	mock.lockRemoveApprovedUsers.Lock()
	mock.calls.RemoveApprovedUsers = nil
	mock.lockRemoveApprovedUsers.Unlock()
}

// UpdateHam calls UpdateHamFunc.
func (mock *DetectorMock) UpdateHam(msg string) error {
	if mock.UpdateHamFunc == nil {
		panic("DetectorMock.UpdateHamFunc: method is nil but Detector.UpdateHam was just called")
	}
	callInfo := struct {
		Msg string
	}{
		Msg: msg,
	}
	mock.lockUpdateHam.Lock()
	mock.calls.UpdateHam = append(mock.calls.UpdateHam, callInfo)
	mock.lockUpdateHam.Unlock()
	return mock.UpdateHamFunc(msg)
}

// UpdateHamCalls gets all the calls that were made to UpdateHam.
// Check the length with:
//
//	len(mockedDetector.UpdateHamCalls())
func (mock *DetectorMock) UpdateHamCalls() []struct {
	Msg string
} {
	var calls []struct {
		Msg string
	}
	mock.lockUpdateHam.RLock()
	calls = mock.calls.UpdateHam
	mock.lockUpdateHam.RUnlock()
	return calls
}

// ResetUpdateHamCalls reset all the calls that were made to UpdateHam.
func (mock *DetectorMock) ResetUpdateHamCalls() {
	mock.lockUpdateHam.Lock()
	mock.calls.UpdateHam = nil
	mock.lockUpdateHam.Unlock()
}

// UpdateSpam calls UpdateSpamFunc.
func (mock *DetectorMock) UpdateSpam(msg string) error {
	if mock.UpdateSpamFunc == nil {
		panic("DetectorMock.UpdateSpamFunc: method is nil but Detector.UpdateSpam was just called")
	}
	callInfo := struct {
		Msg string
	}{
		Msg: msg,
	}
	mock.lockUpdateSpam.Lock()
	mock.calls.UpdateSpam = append(mock.calls.UpdateSpam, callInfo)
	mock.lockUpdateSpam.Unlock()
	return mock.UpdateSpamFunc(msg)
}

// UpdateSpamCalls gets all the calls that were made to UpdateSpam.
// Check the length with:
//
//	len(mockedDetector.UpdateSpamCalls())
func (mock *DetectorMock) UpdateSpamCalls() []struct {
	Msg string
} {
	var calls []struct {
		Msg string
	}
	mock.lockUpdateSpam.RLock()
	calls = mock.calls.UpdateSpam
	mock.lockUpdateSpam.RUnlock()
	return calls
}

// ResetUpdateSpamCalls reset all the calls that were made to UpdateSpam.
func (mock *DetectorMock) ResetUpdateSpamCalls() {
	mock.lockUpdateSpam.Lock()
	mock.calls.UpdateSpam = nil
	mock.lockUpdateSpam.Unlock()
}

// ResetCalls reset all the calls that were made to all mocked methods.
func (mock *DetectorMock) ResetCalls() {
	mock.lockAddApprovedUsers.Lock()
	mock.calls.AddApprovedUsers = nil
	mock.lockAddApprovedUsers.Unlock()

	mock.lockApprovedUsers.Lock()
	mock.calls.ApprovedUsers = nil
	mock.lockApprovedUsers.Unlock()

	mock.lockCheck.Lock()
	mock.calls.Check = nil
	mock.lockCheck.Unlock()

	mock.lockLoadSamples.Lock()
	mock.calls.LoadSamples = nil
	mock.lockLoadSamples.Unlock()

	mock.lockLoadStopWords.Lock()
	mock.calls.LoadStopWords = nil
	mock.lockLoadStopWords.Unlock()

	mock.lockRemoveApprovedUsers.Lock()
	mock.calls.RemoveApprovedUsers = nil
	mock.lockRemoveApprovedUsers.Unlock()

	mock.lockUpdateHam.Lock()
	mock.calls.UpdateHam = nil
	mock.lockUpdateHam.Unlock()

	mock.lockUpdateSpam.Lock()
	mock.calls.UpdateSpam = nil
	mock.lockUpdateSpam.Unlock()
}
