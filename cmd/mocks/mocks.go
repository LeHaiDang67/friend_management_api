package mocks

import (
	"friend_management/intenal/model"

	"github.com/stretchr/testify/mock"
)

//Service is...
type Service struct {
	mock.Mock
}

// ConnectFriends is...
func (st *Service) ConnectFriends(req model.FriendConnectionRequest) (model.BasicResponse, error) {
	returnVals := st.Called(req)
	r0 := returnVals.Get(0).(model.BasicResponse)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return r0, r1
}

//GetUser is...
func (st *Service) GetUser(email string) (model.User, error) {
	returnVals := st.Called(email)
	r0 := returnVals.Get(0).(model.User)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return r0, r1

}

//GetAllUsers is...
func (st *Service) GetAllUsers() ([]model.User, error) {
	returnVals := st.Called()
	r0 := returnVals.Get(0).([]model.User)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return r0, r1
}

//SendUpdate is...
func (st *Service) SendUpdate(sendRequest model.SendUpdateRequest) (model.SendUpdateResponse, error) {
	returnVals := st.Called(sendRequest)
	r0 := returnVals.Get(0).(model.SendUpdateResponse)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return r0, r1
}

//Blocked is ...
func (st *Service) Blocked(subRequest model.SubscriptionRequest) (model.BasicResponse, error) {
	returnVals := st.Called(subRequest)
	r0 := returnVals.Get(0).(model.BasicResponse)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return r0, r1
}

//Subscription is ...
func (st *Service) Subscription(subRequest model.SubscriptionRequest) (model.BasicResponse, error) {
	returnVals := st.Called(subRequest)
	r0 := returnVals.Get(0).(model.BasicResponse)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return r0, r1
}

//CommonFriends is...
func (st *Service) CommonFriends(commonFriends model.CommonFriendRequest) (model.FriendListResponse, error) {
	returnVals := st.Called(commonFriends)
	r0 := returnVals.Get(0).(model.FriendListResponse)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return r0, r1
}

//FriendList is...
func (st *Service) FriendList(email string) (model.FriendListResponse, error) {
	returnVals := st.Called(email)
	r0 := returnVals.Get(0).(model.FriendListResponse)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return r0, r1
}

//CreateNewUser is...
func (st *Service) CreateNewUser(user model.User) (model.BasicResponse, error) {
	returnVals := st.Called(user)
	r0 := returnVals.Get(0).(model.BasicResponse)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return r0, r1
}
