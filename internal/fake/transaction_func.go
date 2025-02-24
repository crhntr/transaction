// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"github.com/crhntr/transaction"
	pgx "github.com/jackc/pgx/v5"
)

type Func struct {
	Stub        func(context.Context, pgx.Tx) error
	mutex       sync.RWMutex
	argsForCall []struct {
		arg1 context.Context
		arg2 pgx.Tx
	}
	returns struct {
		result1 error
	}
	returnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Func) Spy(arg1 context.Context, arg2 pgx.Tx) error {
	fake.mutex.Lock()
	ret, specificReturn := fake.returnsOnCall[len(fake.argsForCall)]
	fake.argsForCall = append(fake.argsForCall, struct {
		arg1 context.Context
		arg2 pgx.Tx
	}{arg1, arg2})
	stub := fake.Stub
	returns := fake.returns
	fake.recordInvocation("Func", []interface{}{arg1, arg2})
	fake.mutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return returns.result1
}

func (fake *Func) CallCount() int {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return len(fake.argsForCall)
}

func (fake *Func) Calls(stub func(context.Context, pgx.Tx) error) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = stub
}

func (fake *Func) ArgsForCall(i int) (context.Context, pgx.Tx) {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return fake.argsForCall[i].arg1, fake.argsForCall[i].arg2
}

func (fake *Func) Returns(result1 error) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = nil
	fake.returns = struct {
		result1 error
	}{result1}
}

func (fake *Func) ReturnsOnCall(i int, result1 error) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = nil
	if fake.returnsOnCall == nil {
		fake.returnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.returnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Func) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Func) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ transaction.Func = new(Func).Spy
