// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"sync"
	"time"

	"github.com/operator-framework/operator-registry/pkg/api"
	"github.com/operator-framework/operator-registry/pkg/client"
)

type FakeInterface struct {
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	closeReturns struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	GetBundleStub        func(context.Context, string, string, string) (*api.Bundle, error)
	getBundleMutex       sync.RWMutex
	getBundleArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}
	getBundleReturns struct {
		result1 *api.Bundle
		result2 error
	}
	getBundleReturnsOnCall map[int]struct {
		result1 *api.Bundle
		result2 error
	}
	GetBundleInPackageChannelStub        func(context.Context, string, string) (*api.Bundle, error)
	getBundleInPackageChannelMutex       sync.RWMutex
	getBundleInPackageChannelArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}
	getBundleInPackageChannelReturns struct {
		result1 *api.Bundle
		result2 error
	}
	getBundleInPackageChannelReturnsOnCall map[int]struct {
		result1 *api.Bundle
		result2 error
	}
	GetBundleThatProvidesStub        func(context.Context, string, string, string) (*api.Bundle, error)
	getBundleThatProvidesMutex       sync.RWMutex
	getBundleThatProvidesArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}
	getBundleThatProvidesReturns struct {
		result1 *api.Bundle
		result2 error
	}
	getBundleThatProvidesReturnsOnCall map[int]struct {
		result1 *api.Bundle
		result2 error
	}
	GetReplacementBundleInPackageChannelStub        func(context.Context, string, string, string) (*api.Bundle, error)
	getReplacementBundleInPackageChannelMutex       sync.RWMutex
	getReplacementBundleInPackageChannelArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}
	getReplacementBundleInPackageChannelReturns struct {
		result1 *api.Bundle
		result2 error
	}
	getReplacementBundleInPackageChannelReturnsOnCall map[int]struct {
		result1 *api.Bundle
		result2 error
	}
	HealthCheckStub        func(context.Context, time.Duration) (bool, error)
	healthCheckMutex       sync.RWMutex
	healthCheckArgsForCall []struct {
		arg1 context.Context
		arg2 time.Duration
	}
	healthCheckReturns struct {
		result1 bool
		result2 error
	}
	healthCheckReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeInterface) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.closeReturns
	return fakeReturns.result1
}

func (fake *FakeInterface) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeInterface) CloseCalls(stub func() error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeInterface) CloseReturns(result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeInterface) CloseReturnsOnCall(i int, result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeInterface) GetBundle(arg1 context.Context, arg2 string, arg3 string, arg4 string) (*api.Bundle, error) {
	fake.getBundleMutex.Lock()
	ret, specificReturn := fake.getBundleReturnsOnCall[len(fake.getBundleArgsForCall)]
	fake.getBundleArgsForCall = append(fake.getBundleArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("GetBundle", []interface{}{arg1, arg2, arg3, arg4})
	fake.getBundleMutex.Unlock()
	if fake.GetBundleStub != nil {
		return fake.GetBundleStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBundleReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeInterface) GetBundleCallCount() int {
	fake.getBundleMutex.RLock()
	defer fake.getBundleMutex.RUnlock()
	return len(fake.getBundleArgsForCall)
}

func (fake *FakeInterface) GetBundleCalls(stub func(context.Context, string, string, string) (*api.Bundle, error)) {
	fake.getBundleMutex.Lock()
	defer fake.getBundleMutex.Unlock()
	fake.GetBundleStub = stub
}

func (fake *FakeInterface) GetBundleArgsForCall(i int) (context.Context, string, string, string) {
	fake.getBundleMutex.RLock()
	defer fake.getBundleMutex.RUnlock()
	argsForCall := fake.getBundleArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeInterface) GetBundleReturns(result1 *api.Bundle, result2 error) {
	fake.getBundleMutex.Lock()
	defer fake.getBundleMutex.Unlock()
	fake.GetBundleStub = nil
	fake.getBundleReturns = struct {
		result1 *api.Bundle
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) GetBundleReturnsOnCall(i int, result1 *api.Bundle, result2 error) {
	fake.getBundleMutex.Lock()
	defer fake.getBundleMutex.Unlock()
	fake.GetBundleStub = nil
	if fake.getBundleReturnsOnCall == nil {
		fake.getBundleReturnsOnCall = make(map[int]struct {
			result1 *api.Bundle
			result2 error
		})
	}
	fake.getBundleReturnsOnCall[i] = struct {
		result1 *api.Bundle
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) GetBundleInPackageChannel(arg1 context.Context, arg2 string, arg3 string) (*api.Bundle, error) {
	fake.getBundleInPackageChannelMutex.Lock()
	ret, specificReturn := fake.getBundleInPackageChannelReturnsOnCall[len(fake.getBundleInPackageChannelArgsForCall)]
	fake.getBundleInPackageChannelArgsForCall = append(fake.getBundleInPackageChannelArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetBundleInPackageChannel", []interface{}{arg1, arg2, arg3})
	fake.getBundleInPackageChannelMutex.Unlock()
	if fake.GetBundleInPackageChannelStub != nil {
		return fake.GetBundleInPackageChannelStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBundleInPackageChannelReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeInterface) GetBundleInPackageChannelCallCount() int {
	fake.getBundleInPackageChannelMutex.RLock()
	defer fake.getBundleInPackageChannelMutex.RUnlock()
	return len(fake.getBundleInPackageChannelArgsForCall)
}

func (fake *FakeInterface) GetBundleInPackageChannelCalls(stub func(context.Context, string, string) (*api.Bundle, error)) {
	fake.getBundleInPackageChannelMutex.Lock()
	defer fake.getBundleInPackageChannelMutex.Unlock()
	fake.GetBundleInPackageChannelStub = stub
}

func (fake *FakeInterface) GetBundleInPackageChannelArgsForCall(i int) (context.Context, string, string) {
	fake.getBundleInPackageChannelMutex.RLock()
	defer fake.getBundleInPackageChannelMutex.RUnlock()
	argsForCall := fake.getBundleInPackageChannelArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeInterface) GetBundleInPackageChannelReturns(result1 *api.Bundle, result2 error) {
	fake.getBundleInPackageChannelMutex.Lock()
	defer fake.getBundleInPackageChannelMutex.Unlock()
	fake.GetBundleInPackageChannelStub = nil
	fake.getBundleInPackageChannelReturns = struct {
		result1 *api.Bundle
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) GetBundleInPackageChannelReturnsOnCall(i int, result1 *api.Bundle, result2 error) {
	fake.getBundleInPackageChannelMutex.Lock()
	defer fake.getBundleInPackageChannelMutex.Unlock()
	fake.GetBundleInPackageChannelStub = nil
	if fake.getBundleInPackageChannelReturnsOnCall == nil {
		fake.getBundleInPackageChannelReturnsOnCall = make(map[int]struct {
			result1 *api.Bundle
			result2 error
		})
	}
	fake.getBundleInPackageChannelReturnsOnCall[i] = struct {
		result1 *api.Bundle
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) GetBundleThatProvides(arg1 context.Context, arg2 string, arg3 string, arg4 string) (*api.Bundle, error) {
	fake.getBundleThatProvidesMutex.Lock()
	ret, specificReturn := fake.getBundleThatProvidesReturnsOnCall[len(fake.getBundleThatProvidesArgsForCall)]
	fake.getBundleThatProvidesArgsForCall = append(fake.getBundleThatProvidesArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("GetBundleThatProvides", []interface{}{arg1, arg2, arg3, arg4})
	fake.getBundleThatProvidesMutex.Unlock()
	if fake.GetBundleThatProvidesStub != nil {
		return fake.GetBundleThatProvidesStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBundleThatProvidesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeInterface) GetBundleThatProvidesCallCount() int {
	fake.getBundleThatProvidesMutex.RLock()
	defer fake.getBundleThatProvidesMutex.RUnlock()
	return len(fake.getBundleThatProvidesArgsForCall)
}

func (fake *FakeInterface) GetBundleThatProvidesCalls(stub func(context.Context, string, string, string) (*api.Bundle, error)) {
	fake.getBundleThatProvidesMutex.Lock()
	defer fake.getBundleThatProvidesMutex.Unlock()
	fake.GetBundleThatProvidesStub = stub
}

func (fake *FakeInterface) GetBundleThatProvidesArgsForCall(i int) (context.Context, string, string, string) {
	fake.getBundleThatProvidesMutex.RLock()
	defer fake.getBundleThatProvidesMutex.RUnlock()
	argsForCall := fake.getBundleThatProvidesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeInterface) GetBundleThatProvidesReturns(result1 *api.Bundle, result2 error) {
	fake.getBundleThatProvidesMutex.Lock()
	defer fake.getBundleThatProvidesMutex.Unlock()
	fake.GetBundleThatProvidesStub = nil
	fake.getBundleThatProvidesReturns = struct {
		result1 *api.Bundle
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) GetBundleThatProvidesReturnsOnCall(i int, result1 *api.Bundle, result2 error) {
	fake.getBundleThatProvidesMutex.Lock()
	defer fake.getBundleThatProvidesMutex.Unlock()
	fake.GetBundleThatProvidesStub = nil
	if fake.getBundleThatProvidesReturnsOnCall == nil {
		fake.getBundleThatProvidesReturnsOnCall = make(map[int]struct {
			result1 *api.Bundle
			result2 error
		})
	}
	fake.getBundleThatProvidesReturnsOnCall[i] = struct {
		result1 *api.Bundle
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) GetReplacementBundleInPackageChannel(arg1 context.Context, arg2 string, arg3 string, arg4 string) (*api.Bundle, error) {
	fake.getReplacementBundleInPackageChannelMutex.Lock()
	ret, specificReturn := fake.getReplacementBundleInPackageChannelReturnsOnCall[len(fake.getReplacementBundleInPackageChannelArgsForCall)]
	fake.getReplacementBundleInPackageChannelArgsForCall = append(fake.getReplacementBundleInPackageChannelArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("GetReplacementBundleInPackageChannel", []interface{}{arg1, arg2, arg3, arg4})
	fake.getReplacementBundleInPackageChannelMutex.Unlock()
	if fake.GetReplacementBundleInPackageChannelStub != nil {
		return fake.GetReplacementBundleInPackageChannelStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getReplacementBundleInPackageChannelReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeInterface) GetReplacementBundleInPackageChannelCallCount() int {
	fake.getReplacementBundleInPackageChannelMutex.RLock()
	defer fake.getReplacementBundleInPackageChannelMutex.RUnlock()
	return len(fake.getReplacementBundleInPackageChannelArgsForCall)
}

func (fake *FakeInterface) GetReplacementBundleInPackageChannelCalls(stub func(context.Context, string, string, string) (*api.Bundle, error)) {
	fake.getReplacementBundleInPackageChannelMutex.Lock()
	defer fake.getReplacementBundleInPackageChannelMutex.Unlock()
	fake.GetReplacementBundleInPackageChannelStub = stub
}

func (fake *FakeInterface) GetReplacementBundleInPackageChannelArgsForCall(i int) (context.Context, string, string, string) {
	fake.getReplacementBundleInPackageChannelMutex.RLock()
	defer fake.getReplacementBundleInPackageChannelMutex.RUnlock()
	argsForCall := fake.getReplacementBundleInPackageChannelArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeInterface) GetReplacementBundleInPackageChannelReturns(result1 *api.Bundle, result2 error) {
	fake.getReplacementBundleInPackageChannelMutex.Lock()
	defer fake.getReplacementBundleInPackageChannelMutex.Unlock()
	fake.GetReplacementBundleInPackageChannelStub = nil
	fake.getReplacementBundleInPackageChannelReturns = struct {
		result1 *api.Bundle
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) GetReplacementBundleInPackageChannelReturnsOnCall(i int, result1 *api.Bundle, result2 error) {
	fake.getReplacementBundleInPackageChannelMutex.Lock()
	defer fake.getReplacementBundleInPackageChannelMutex.Unlock()
	fake.GetReplacementBundleInPackageChannelStub = nil
	if fake.getReplacementBundleInPackageChannelReturnsOnCall == nil {
		fake.getReplacementBundleInPackageChannelReturnsOnCall = make(map[int]struct {
			result1 *api.Bundle
			result2 error
		})
	}
	fake.getReplacementBundleInPackageChannelReturnsOnCall[i] = struct {
		result1 *api.Bundle
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) HealthCheck(arg1 context.Context, arg2 time.Duration) (bool, error) {
	fake.healthCheckMutex.Lock()
	ret, specificReturn := fake.healthCheckReturnsOnCall[len(fake.healthCheckArgsForCall)]
	fake.healthCheckArgsForCall = append(fake.healthCheckArgsForCall, struct {
		arg1 context.Context
		arg2 time.Duration
	}{arg1, arg2})
	fake.recordInvocation("HealthCheck", []interface{}{arg1, arg2})
	fake.healthCheckMutex.Unlock()
	if fake.HealthCheckStub != nil {
		return fake.HealthCheckStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.healthCheckReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeInterface) HealthCheckCallCount() int {
	fake.healthCheckMutex.RLock()
	defer fake.healthCheckMutex.RUnlock()
	return len(fake.healthCheckArgsForCall)
}

func (fake *FakeInterface) HealthCheckCalls(stub func(context.Context, time.Duration) (bool, error)) {
	fake.healthCheckMutex.Lock()
	defer fake.healthCheckMutex.Unlock()
	fake.HealthCheckStub = stub
}

func (fake *FakeInterface) HealthCheckArgsForCall(i int) (context.Context, time.Duration) {
	fake.healthCheckMutex.RLock()
	defer fake.healthCheckMutex.RUnlock()
	argsForCall := fake.healthCheckArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeInterface) HealthCheckReturns(result1 bool, result2 error) {
	fake.healthCheckMutex.Lock()
	defer fake.healthCheckMutex.Unlock()
	fake.HealthCheckStub = nil
	fake.healthCheckReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) HealthCheckReturnsOnCall(i int, result1 bool, result2 error) {
	fake.healthCheckMutex.Lock()
	defer fake.healthCheckMutex.Unlock()
	fake.HealthCheckStub = nil
	if fake.healthCheckReturnsOnCall == nil {
		fake.healthCheckReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.healthCheckReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeInterface) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.getBundleMutex.RLock()
	defer fake.getBundleMutex.RUnlock()
	fake.getBundleInPackageChannelMutex.RLock()
	defer fake.getBundleInPackageChannelMutex.RUnlock()
	fake.getBundleThatProvidesMutex.RLock()
	defer fake.getBundleThatProvidesMutex.RUnlock()
	fake.getReplacementBundleInPackageChannelMutex.RLock()
	defer fake.getReplacementBundleInPackageChannelMutex.RUnlock()
	fake.healthCheckMutex.RLock()
	defer fake.healthCheckMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeInterface) recordInvocation(key string, args []interface{}) {
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

var _ client.Interface = new(FakeInterface)
