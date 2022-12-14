// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package server

import (
	"foodie/server/apierr"
	"sync"
	"time"

	"github.com/rs/xid"
)

// Ensure, that AuthorizerMock does implement Authorizer.
// If this is not the case, regenerate this file with moq.
var _ Authorizer = &AuthorizerMock{}

// AuthorizerMock is a mock implementation of Authorizer.
//
//	func TestSomethingThatUsesAuthorizer(t *testing.T) {
//
//		// make and configure a mocked Authorizer
//		mockedAuthorizer := &AuthorizerMock{
//			IssueFunc: func(id xid.ID, admin bool, tstamp time.Time) ([]byte, error) {
//				panic("mock out the Issue method")
//			},
//			ParseFunc: func(data []byte, tstamp time.Time) (xid.ID, bool, *apierr.Error) {
//				panic("mock out the Parse method")
//			},
//		}
//
//		// use mockedAuthorizer in code that requires Authorizer
//		// and then make assertions.
//
//	}
type AuthorizerMock struct {
	// IssueFunc mocks the Issue method.
	IssueFunc func(id xid.ID, admin bool, tstamp time.Time) ([]byte, error)

	// ParseFunc mocks the Parse method.
	ParseFunc func(data []byte, tstamp time.Time) (xid.ID, bool, *apierr.Error)

	// calls tracks calls to the methods.
	calls struct {
		// Issue holds details about calls to the Issue method.
		Issue []struct {
			// ID is the id argument value.
			ID xid.ID
			// Admin is the admin argument value.
			Admin bool
			// Tstamp is the tstamp argument value.
			Tstamp time.Time
		}
		// Parse holds details about calls to the Parse method.
		Parse []struct {
			// Data is the data argument value.
			Data []byte
			// Tstamp is the tstamp argument value.
			Tstamp time.Time
		}
	}
	lockIssue sync.RWMutex
	lockParse sync.RWMutex
}

// Issue calls IssueFunc.
func (mock *AuthorizerMock) Issue(id xid.ID, admin bool, tstamp time.Time) ([]byte, error) {
	callInfo := struct {
		ID     xid.ID
		Admin  bool
		Tstamp time.Time
	}{
		ID:     id,
		Admin:  admin,
		Tstamp: tstamp,
	}
	mock.lockIssue.Lock()
	mock.calls.Issue = append(mock.calls.Issue, callInfo)
	mock.lockIssue.Unlock()
	if mock.IssueFunc == nil {
		var (
			bytesOut []byte
			errOut   error
		)
		return bytesOut, errOut
	}
	return mock.IssueFunc(id, admin, tstamp)
}

// IssueCalls gets all the calls that were made to Issue.
// Check the length with:
//
//	len(mockedAuthorizer.IssueCalls())
func (mock *AuthorizerMock) IssueCalls() []struct {
	ID     xid.ID
	Admin  bool
	Tstamp time.Time
} {
	var calls []struct {
		ID     xid.ID
		Admin  bool
		Tstamp time.Time
	}
	mock.lockIssue.RLock()
	calls = mock.calls.Issue
	mock.lockIssue.RUnlock()
	return calls
}

// Parse calls ParseFunc.
func (mock *AuthorizerMock) Parse(data []byte, tstamp time.Time) (xid.ID, bool, *apierr.Error) {
	callInfo := struct {
		Data   []byte
		Tstamp time.Time
	}{
		Data:   data,
		Tstamp: tstamp,
	}
	mock.lockParse.Lock()
	mock.calls.Parse = append(mock.calls.Parse, callInfo)
	mock.lockParse.Unlock()
	if mock.ParseFunc == nil {
		var (
			iDOut    xid.ID
			bOut     bool
			errorOut *apierr.Error
		)
		return iDOut, bOut, errorOut
	}
	return mock.ParseFunc(data, tstamp)
}

// ParseCalls gets all the calls that were made to Parse.
// Check the length with:
//
//	len(mockedAuthorizer.ParseCalls())
func (mock *AuthorizerMock) ParseCalls() []struct {
	Data   []byte
	Tstamp time.Time
} {
	var calls []struct {
		Data   []byte
		Tstamp time.Time
	}
	mock.lockParse.RLock()
	calls = mock.calls.Parse
	mock.lockParse.RUnlock()
	return calls
}
