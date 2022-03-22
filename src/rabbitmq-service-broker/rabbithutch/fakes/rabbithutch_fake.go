// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"rabbitmq-service-broker/rabbithutch"
	"sync"
)

type FakeRabbitHutch struct {
	AssignPermissionsToStub        func(string, string) error
	assignPermissionsToMutex       sync.RWMutex
	assignPermissionsToArgsForCall []struct {
		arg1 string
		arg2 string
	}
	assignPermissionsToReturns struct {
		result1 error
	}
	assignPermissionsToReturnsOnCall map[int]struct {
		result1 error
	}
	CreatePolicyStub        func(string, string, int, map[string]interface{}) error
	createPolicyMutex       sync.RWMutex
	createPolicyArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 int
		arg4 map[string]interface{}
	}
	createPolicyReturns struct {
		result1 error
	}
	createPolicyReturnsOnCall map[int]struct {
		result1 error
	}
	CreateUserAndGrantPermissionsStub        func(string, string, string) (string, error)
	createUserAndGrantPermissionsMutex       sync.RWMutex
	createUserAndGrantPermissionsArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	createUserAndGrantPermissionsReturns struct {
		result1 string
		result2 error
	}
	createUserAndGrantPermissionsReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	DeleteUserStub        func(string) error
	deleteUserMutex       sync.RWMutex
	deleteUserArgsForCall []struct {
		arg1 string
	}
	deleteUserReturns struct {
		result1 error
	}
	deleteUserReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteUserAndConnectionsStub        func(string) error
	deleteUserAndConnectionsMutex       sync.RWMutex
	deleteUserAndConnectionsArgsForCall []struct {
		arg1 string
	}
	deleteUserAndConnectionsReturns struct {
		result1 error
	}
	deleteUserAndConnectionsReturnsOnCall map[int]struct {
		result1 error
	}
	ProtocolPortsStub        func() (map[string]int, error)
	protocolPortsMutex       sync.RWMutex
	protocolPortsArgsForCall []struct {
	}
	protocolPortsReturns struct {
		result1 map[string]int
		result2 error
	}
	protocolPortsReturnsOnCall map[int]struct {
		result1 map[string]int
		result2 error
	}
	UserListStub        func() ([]string, error)
	userListMutex       sync.RWMutex
	userListArgsForCall []struct {
	}
	userListReturns struct {
		result1 []string
		result2 error
	}
	userListReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	VHostCreateStub        func(string) error
	vHostCreateMutex       sync.RWMutex
	vHostCreateArgsForCall []struct {
		arg1 string
	}
	vHostCreateReturns struct {
		result1 error
	}
	vHostCreateReturnsOnCall map[int]struct {
		result1 error
	}
	VHostDeleteStub        func(string) error
	vHostDeleteMutex       sync.RWMutex
	vHostDeleteArgsForCall []struct {
		arg1 string
	}
	vHostDeleteReturns struct {
		result1 error
	}
	vHostDeleteReturnsOnCall map[int]struct {
		result1 error
	}
	VHostExistsStub        func(string) (bool, error)
	vHostExistsMutex       sync.RWMutex
	vHostExistsArgsForCall []struct {
		arg1 string
	}
	vHostExistsReturns struct {
		result1 bool
		result2 error
	}
	vHostExistsReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRabbitHutch) AssignPermissionsTo(arg1 string, arg2 string) error {
	fake.assignPermissionsToMutex.Lock()
	ret, specificReturn := fake.assignPermissionsToReturnsOnCall[len(fake.assignPermissionsToArgsForCall)]
	fake.assignPermissionsToArgsForCall = append(fake.assignPermissionsToArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.AssignPermissionsToStub
	fakeReturns := fake.assignPermissionsToReturns
	fake.recordInvocation("AssignPermissionsTo", []interface{}{arg1, arg2})
	fake.assignPermissionsToMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRabbitHutch) AssignPermissionsToCallCount() int {
	fake.assignPermissionsToMutex.RLock()
	defer fake.assignPermissionsToMutex.RUnlock()
	return len(fake.assignPermissionsToArgsForCall)
}

func (fake *FakeRabbitHutch) AssignPermissionsToCalls(stub func(string, string) error) {
	fake.assignPermissionsToMutex.Lock()
	defer fake.assignPermissionsToMutex.Unlock()
	fake.AssignPermissionsToStub = stub
}

func (fake *FakeRabbitHutch) AssignPermissionsToArgsForCall(i int) (string, string) {
	fake.assignPermissionsToMutex.RLock()
	defer fake.assignPermissionsToMutex.RUnlock()
	argsForCall := fake.assignPermissionsToArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeRabbitHutch) AssignPermissionsToReturns(result1 error) {
	fake.assignPermissionsToMutex.Lock()
	defer fake.assignPermissionsToMutex.Unlock()
	fake.AssignPermissionsToStub = nil
	fake.assignPermissionsToReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) AssignPermissionsToReturnsOnCall(i int, result1 error) {
	fake.assignPermissionsToMutex.Lock()
	defer fake.assignPermissionsToMutex.Unlock()
	fake.AssignPermissionsToStub = nil
	if fake.assignPermissionsToReturnsOnCall == nil {
		fake.assignPermissionsToReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.assignPermissionsToReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) CreatePolicy(arg1 string, arg2 string, arg3 int, arg4 map[string]interface{}) error {
	fake.createPolicyMutex.Lock()
	ret, specificReturn := fake.createPolicyReturnsOnCall[len(fake.createPolicyArgsForCall)]
	fake.createPolicyArgsForCall = append(fake.createPolicyArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 int
		arg4 map[string]interface{}
	}{arg1, arg2, arg3, arg4})
	stub := fake.CreatePolicyStub
	fakeReturns := fake.createPolicyReturns
	fake.recordInvocation("CreatePolicy", []interface{}{arg1, arg2, arg3, arg4})
	fake.createPolicyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRabbitHutch) CreatePolicyCallCount() int {
	fake.createPolicyMutex.RLock()
	defer fake.createPolicyMutex.RUnlock()
	return len(fake.createPolicyArgsForCall)
}

func (fake *FakeRabbitHutch) CreatePolicyCalls(stub func(string, string, int, map[string]interface{}) error) {
	fake.createPolicyMutex.Lock()
	defer fake.createPolicyMutex.Unlock()
	fake.CreatePolicyStub = stub
}

func (fake *FakeRabbitHutch) CreatePolicyArgsForCall(i int) (string, string, int, map[string]interface{}) {
	fake.createPolicyMutex.RLock()
	defer fake.createPolicyMutex.RUnlock()
	argsForCall := fake.createPolicyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeRabbitHutch) CreatePolicyReturns(result1 error) {
	fake.createPolicyMutex.Lock()
	defer fake.createPolicyMutex.Unlock()
	fake.CreatePolicyStub = nil
	fake.createPolicyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) CreatePolicyReturnsOnCall(i int, result1 error) {
	fake.createPolicyMutex.Lock()
	defer fake.createPolicyMutex.Unlock()
	fake.CreatePolicyStub = nil
	if fake.createPolicyReturnsOnCall == nil {
		fake.createPolicyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createPolicyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) CreateUserAndGrantPermissions(arg1 string, arg2 string, arg3 string) (string, error) {
	fake.createUserAndGrantPermissionsMutex.Lock()
	ret, specificReturn := fake.createUserAndGrantPermissionsReturnsOnCall[len(fake.createUserAndGrantPermissionsArgsForCall)]
	fake.createUserAndGrantPermissionsArgsForCall = append(fake.createUserAndGrantPermissionsArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.CreateUserAndGrantPermissionsStub
	fakeReturns := fake.createUserAndGrantPermissionsReturns
	fake.recordInvocation("CreateUserAndGrantPermissions", []interface{}{arg1, arg2, arg3})
	fake.createUserAndGrantPermissionsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRabbitHutch) CreateUserAndGrantPermissionsCallCount() int {
	fake.createUserAndGrantPermissionsMutex.RLock()
	defer fake.createUserAndGrantPermissionsMutex.RUnlock()
	return len(fake.createUserAndGrantPermissionsArgsForCall)
}

func (fake *FakeRabbitHutch) CreateUserAndGrantPermissionsCalls(stub func(string, string, string) (string, error)) {
	fake.createUserAndGrantPermissionsMutex.Lock()
	defer fake.createUserAndGrantPermissionsMutex.Unlock()
	fake.CreateUserAndGrantPermissionsStub = stub
}

func (fake *FakeRabbitHutch) CreateUserAndGrantPermissionsArgsForCall(i int) (string, string, string) {
	fake.createUserAndGrantPermissionsMutex.RLock()
	defer fake.createUserAndGrantPermissionsMutex.RUnlock()
	argsForCall := fake.createUserAndGrantPermissionsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeRabbitHutch) CreateUserAndGrantPermissionsReturns(result1 string, result2 error) {
	fake.createUserAndGrantPermissionsMutex.Lock()
	defer fake.createUserAndGrantPermissionsMutex.Unlock()
	fake.CreateUserAndGrantPermissionsStub = nil
	fake.createUserAndGrantPermissionsReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitHutch) CreateUserAndGrantPermissionsReturnsOnCall(i int, result1 string, result2 error) {
	fake.createUserAndGrantPermissionsMutex.Lock()
	defer fake.createUserAndGrantPermissionsMutex.Unlock()
	fake.CreateUserAndGrantPermissionsStub = nil
	if fake.createUserAndGrantPermissionsReturnsOnCall == nil {
		fake.createUserAndGrantPermissionsReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.createUserAndGrantPermissionsReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitHutch) DeleteUser(arg1 string) error {
	fake.deleteUserMutex.Lock()
	ret, specificReturn := fake.deleteUserReturnsOnCall[len(fake.deleteUserArgsForCall)]
	fake.deleteUserArgsForCall = append(fake.deleteUserArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.DeleteUserStub
	fakeReturns := fake.deleteUserReturns
	fake.recordInvocation("DeleteUser", []interface{}{arg1})
	fake.deleteUserMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRabbitHutch) DeleteUserCallCount() int {
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	return len(fake.deleteUserArgsForCall)
}

func (fake *FakeRabbitHutch) DeleteUserCalls(stub func(string) error) {
	fake.deleteUserMutex.Lock()
	defer fake.deleteUserMutex.Unlock()
	fake.DeleteUserStub = stub
}

func (fake *FakeRabbitHutch) DeleteUserArgsForCall(i int) string {
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	argsForCall := fake.deleteUserArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRabbitHutch) DeleteUserReturns(result1 error) {
	fake.deleteUserMutex.Lock()
	defer fake.deleteUserMutex.Unlock()
	fake.DeleteUserStub = nil
	fake.deleteUserReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) DeleteUserReturnsOnCall(i int, result1 error) {
	fake.deleteUserMutex.Lock()
	defer fake.deleteUserMutex.Unlock()
	fake.DeleteUserStub = nil
	if fake.deleteUserReturnsOnCall == nil {
		fake.deleteUserReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteUserReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) DeleteUserAndConnections(arg1 string) error {
	fake.deleteUserAndConnectionsMutex.Lock()
	ret, specificReturn := fake.deleteUserAndConnectionsReturnsOnCall[len(fake.deleteUserAndConnectionsArgsForCall)]
	fake.deleteUserAndConnectionsArgsForCall = append(fake.deleteUserAndConnectionsArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.DeleteUserAndConnectionsStub
	fakeReturns := fake.deleteUserAndConnectionsReturns
	fake.recordInvocation("DeleteUserAndConnections", []interface{}{arg1})
	fake.deleteUserAndConnectionsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRabbitHutch) DeleteUserAndConnectionsCallCount() int {
	fake.deleteUserAndConnectionsMutex.RLock()
	defer fake.deleteUserAndConnectionsMutex.RUnlock()
	return len(fake.deleteUserAndConnectionsArgsForCall)
}

func (fake *FakeRabbitHutch) DeleteUserAndConnectionsCalls(stub func(string) error) {
	fake.deleteUserAndConnectionsMutex.Lock()
	defer fake.deleteUserAndConnectionsMutex.Unlock()
	fake.DeleteUserAndConnectionsStub = stub
}

func (fake *FakeRabbitHutch) DeleteUserAndConnectionsArgsForCall(i int) string {
	fake.deleteUserAndConnectionsMutex.RLock()
	defer fake.deleteUserAndConnectionsMutex.RUnlock()
	argsForCall := fake.deleteUserAndConnectionsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRabbitHutch) DeleteUserAndConnectionsReturns(result1 error) {
	fake.deleteUserAndConnectionsMutex.Lock()
	defer fake.deleteUserAndConnectionsMutex.Unlock()
	fake.DeleteUserAndConnectionsStub = nil
	fake.deleteUserAndConnectionsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) DeleteUserAndConnectionsReturnsOnCall(i int, result1 error) {
	fake.deleteUserAndConnectionsMutex.Lock()
	defer fake.deleteUserAndConnectionsMutex.Unlock()
	fake.DeleteUserAndConnectionsStub = nil
	if fake.deleteUserAndConnectionsReturnsOnCall == nil {
		fake.deleteUserAndConnectionsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteUserAndConnectionsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) ProtocolPorts() (map[string]int, error) {
	fake.protocolPortsMutex.Lock()
	ret, specificReturn := fake.protocolPortsReturnsOnCall[len(fake.protocolPortsArgsForCall)]
	fake.protocolPortsArgsForCall = append(fake.protocolPortsArgsForCall, struct {
	}{})
	stub := fake.ProtocolPortsStub
	fakeReturns := fake.protocolPortsReturns
	fake.recordInvocation("ProtocolPorts", []interface{}{})
	fake.protocolPortsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRabbitHutch) ProtocolPortsCallCount() int {
	fake.protocolPortsMutex.RLock()
	defer fake.protocolPortsMutex.RUnlock()
	return len(fake.protocolPortsArgsForCall)
}

func (fake *FakeRabbitHutch) ProtocolPortsCalls(stub func() (map[string]int, error)) {
	fake.protocolPortsMutex.Lock()
	defer fake.protocolPortsMutex.Unlock()
	fake.ProtocolPortsStub = stub
}

func (fake *FakeRabbitHutch) ProtocolPortsReturns(result1 map[string]int, result2 error) {
	fake.protocolPortsMutex.Lock()
	defer fake.protocolPortsMutex.Unlock()
	fake.ProtocolPortsStub = nil
	fake.protocolPortsReturns = struct {
		result1 map[string]int
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitHutch) ProtocolPortsReturnsOnCall(i int, result1 map[string]int, result2 error) {
	fake.protocolPortsMutex.Lock()
	defer fake.protocolPortsMutex.Unlock()
	fake.ProtocolPortsStub = nil
	if fake.protocolPortsReturnsOnCall == nil {
		fake.protocolPortsReturnsOnCall = make(map[int]struct {
			result1 map[string]int
			result2 error
		})
	}
	fake.protocolPortsReturnsOnCall[i] = struct {
		result1 map[string]int
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitHutch) UserList() ([]string, error) {
	fake.userListMutex.Lock()
	ret, specificReturn := fake.userListReturnsOnCall[len(fake.userListArgsForCall)]
	fake.userListArgsForCall = append(fake.userListArgsForCall, struct {
	}{})
	stub := fake.UserListStub
	fakeReturns := fake.userListReturns
	fake.recordInvocation("UserList", []interface{}{})
	fake.userListMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRabbitHutch) UserListCallCount() int {
	fake.userListMutex.RLock()
	defer fake.userListMutex.RUnlock()
	return len(fake.userListArgsForCall)
}

func (fake *FakeRabbitHutch) UserListCalls(stub func() ([]string, error)) {
	fake.userListMutex.Lock()
	defer fake.userListMutex.Unlock()
	fake.UserListStub = stub
}

func (fake *FakeRabbitHutch) UserListReturns(result1 []string, result2 error) {
	fake.userListMutex.Lock()
	defer fake.userListMutex.Unlock()
	fake.UserListStub = nil
	fake.userListReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitHutch) UserListReturnsOnCall(i int, result1 []string, result2 error) {
	fake.userListMutex.Lock()
	defer fake.userListMutex.Unlock()
	fake.UserListStub = nil
	if fake.userListReturnsOnCall == nil {
		fake.userListReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.userListReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitHutch) VHostCreate(arg1 string) error {
	fake.vHostCreateMutex.Lock()
	ret, specificReturn := fake.vHostCreateReturnsOnCall[len(fake.vHostCreateArgsForCall)]
	fake.vHostCreateArgsForCall = append(fake.vHostCreateArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.VHostCreateStub
	fakeReturns := fake.vHostCreateReturns
	fake.recordInvocation("VHostCreate", []interface{}{arg1})
	fake.vHostCreateMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRabbitHutch) VHostCreateCallCount() int {
	fake.vHostCreateMutex.RLock()
	defer fake.vHostCreateMutex.RUnlock()
	return len(fake.vHostCreateArgsForCall)
}

func (fake *FakeRabbitHutch) VHostCreateCalls(stub func(string) error) {
	fake.vHostCreateMutex.Lock()
	defer fake.vHostCreateMutex.Unlock()
	fake.VHostCreateStub = stub
}

func (fake *FakeRabbitHutch) VHostCreateArgsForCall(i int) string {
	fake.vHostCreateMutex.RLock()
	defer fake.vHostCreateMutex.RUnlock()
	argsForCall := fake.vHostCreateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRabbitHutch) VHostCreateReturns(result1 error) {
	fake.vHostCreateMutex.Lock()
	defer fake.vHostCreateMutex.Unlock()
	fake.VHostCreateStub = nil
	fake.vHostCreateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) VHostCreateReturnsOnCall(i int, result1 error) {
	fake.vHostCreateMutex.Lock()
	defer fake.vHostCreateMutex.Unlock()
	fake.VHostCreateStub = nil
	if fake.vHostCreateReturnsOnCall == nil {
		fake.vHostCreateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.vHostCreateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) VHostDelete(arg1 string) error {
	fake.vHostDeleteMutex.Lock()
	ret, specificReturn := fake.vHostDeleteReturnsOnCall[len(fake.vHostDeleteArgsForCall)]
	fake.vHostDeleteArgsForCall = append(fake.vHostDeleteArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.VHostDeleteStub
	fakeReturns := fake.vHostDeleteReturns
	fake.recordInvocation("VHostDelete", []interface{}{arg1})
	fake.vHostDeleteMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRabbitHutch) VHostDeleteCallCount() int {
	fake.vHostDeleteMutex.RLock()
	defer fake.vHostDeleteMutex.RUnlock()
	return len(fake.vHostDeleteArgsForCall)
}

func (fake *FakeRabbitHutch) VHostDeleteCalls(stub func(string) error) {
	fake.vHostDeleteMutex.Lock()
	defer fake.vHostDeleteMutex.Unlock()
	fake.VHostDeleteStub = stub
}

func (fake *FakeRabbitHutch) VHostDeleteArgsForCall(i int) string {
	fake.vHostDeleteMutex.RLock()
	defer fake.vHostDeleteMutex.RUnlock()
	argsForCall := fake.vHostDeleteArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRabbitHutch) VHostDeleteReturns(result1 error) {
	fake.vHostDeleteMutex.Lock()
	defer fake.vHostDeleteMutex.Unlock()
	fake.VHostDeleteStub = nil
	fake.vHostDeleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) VHostDeleteReturnsOnCall(i int, result1 error) {
	fake.vHostDeleteMutex.Lock()
	defer fake.vHostDeleteMutex.Unlock()
	fake.VHostDeleteStub = nil
	if fake.vHostDeleteReturnsOnCall == nil {
		fake.vHostDeleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.vHostDeleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitHutch) VHostExists(arg1 string) (bool, error) {
	fake.vHostExistsMutex.Lock()
	ret, specificReturn := fake.vHostExistsReturnsOnCall[len(fake.vHostExistsArgsForCall)]
	fake.vHostExistsArgsForCall = append(fake.vHostExistsArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.VHostExistsStub
	fakeReturns := fake.vHostExistsReturns
	fake.recordInvocation("VHostExists", []interface{}{arg1})
	fake.vHostExistsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRabbitHutch) VHostExistsCallCount() int {
	fake.vHostExistsMutex.RLock()
	defer fake.vHostExistsMutex.RUnlock()
	return len(fake.vHostExistsArgsForCall)
}

func (fake *FakeRabbitHutch) VHostExistsCalls(stub func(string) (bool, error)) {
	fake.vHostExistsMutex.Lock()
	defer fake.vHostExistsMutex.Unlock()
	fake.VHostExistsStub = stub
}

func (fake *FakeRabbitHutch) VHostExistsArgsForCall(i int) string {
	fake.vHostExistsMutex.RLock()
	defer fake.vHostExistsMutex.RUnlock()
	argsForCall := fake.vHostExistsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRabbitHutch) VHostExistsReturns(result1 bool, result2 error) {
	fake.vHostExistsMutex.Lock()
	defer fake.vHostExistsMutex.Unlock()
	fake.VHostExistsStub = nil
	fake.vHostExistsReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitHutch) VHostExistsReturnsOnCall(i int, result1 bool, result2 error) {
	fake.vHostExistsMutex.Lock()
	defer fake.vHostExistsMutex.Unlock()
	fake.VHostExistsStub = nil
	if fake.vHostExistsReturnsOnCall == nil {
		fake.vHostExistsReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.vHostExistsReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitHutch) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.assignPermissionsToMutex.RLock()
	defer fake.assignPermissionsToMutex.RUnlock()
	fake.createPolicyMutex.RLock()
	defer fake.createPolicyMutex.RUnlock()
	fake.createUserAndGrantPermissionsMutex.RLock()
	defer fake.createUserAndGrantPermissionsMutex.RUnlock()
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	fake.deleteUserAndConnectionsMutex.RLock()
	defer fake.deleteUserAndConnectionsMutex.RUnlock()
	fake.protocolPortsMutex.RLock()
	defer fake.protocolPortsMutex.RUnlock()
	fake.userListMutex.RLock()
	defer fake.userListMutex.RUnlock()
	fake.vHostCreateMutex.RLock()
	defer fake.vHostCreateMutex.RUnlock()
	fake.vHostDeleteMutex.RLock()
	defer fake.vHostDeleteMutex.RUnlock()
	fake.vHostExistsMutex.RLock()
	defer fake.vHostExistsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeRabbitHutch) recordInvocation(key string, args []interface{}) {
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

var _ rabbithutch.RabbitHutch = new(FakeRabbitHutch)
