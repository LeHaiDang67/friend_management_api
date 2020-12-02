package services

import (
	"errors"
	"friend_management/intenal/db"
	"friend_management/intenal/model"
	"friend_management/intenal/util"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestGetAllUser(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCase := struct {
		name           string
		expectedResult []model.User
	}{
		name:           "Get users success",
		expectedResult: []model.User([]model.User{model.User{Email: "a@gmail.com", Friends: []string{"b@gmail.com"}, Subscription: []string(nil), Blocked: []string(nil)}}),
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/01_get_user.sql"))
	t.Run(testCase.name, func(t *testing.T) {
		mn := NewManager(db)
		result, err := mn.GetAllUsers()
		require.Nil(t, err)
		require.Equal(t, testCase.expectedResult, result)
	})
}

func TestGetUser(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCases := []struct {
		desc           string
		givenUserEmail string
		expectedResult model.User
		expectedError  error
	}{
		{
			desc:           "Should return user",
			givenUserEmail: "a@gmail.com",
			expectedResult: model.User{
				Email:   "a@gmail.com",
				Friends: []string{"b@gmail.com"},
			},
			expectedError: nil,
		},
		{
			desc:           "User not already existed",
			givenUserEmail: "andy@example.com",
			expectedResult: model.User{},
			expectedError:  errors.New("sad"),
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/01_get_user.sql"))
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.GetUser(tc.givenUserEmail)
			if err != nil {
				require.Equal(t, tc.expectedError, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestConnectFriends(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCases := []struct {
		name           string
		request        model.FriendConnectionRequest
		expectedResult model.BasicResponse
		expectedError  error
	}{
		{
			name: "Make friend successfully",
			request: model.FriendConnectionRequest{
				Friends: []string{"a@gmail.com", "b@gmail.com"},
			},
			expectedResult: model.BasicResponse{
				Success: true,
			},
		},
		{
			name: "User not existed",
			request: model.FriendConnectionRequest{
				Friends: []string{"test1@gmail.com", "test2@gmail.com"},
			},
			expectedResult: model.BasicResponse{
				Success: true,
			},
			expectedError: errors.New("Sad"),
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/02_connect_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.ConnectFriends(tt.request)
			// then
			if err != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestFriendList(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCases := []struct {
		name           string
		givenUserEmail model.FriendListRequest
		expectedResult model.FriendListResponse
		expectedError  error
	}{
		{
			name: "Retrieve success ",
			givenUserEmail: model.FriendListRequest{
				Email: "a@gmail.com"},
			expectedResult: model.FriendListResponse{
				Success: true,
				Friends: []string{"b@gmail.com"},
				Count:   1,
			},
		},
		{
			name: "Retrieve failed ",
			givenUserEmail: model.FriendListRequest{
				Email: "andy@gmail.com"},
			expectedResult: model.FriendListResponse{
				Success: true,
			},
			expectedError: errors.New("Sad"),
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/01_get_user.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.FriendList(tt.givenUserEmail)
			// then
			if err != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestCommonFriends(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCases := []struct {
		name           string
		commonFriends  model.CommonFriendRequest
		expectedResult model.FriendListResponse
		expectedError  error
	}{
		{
			name: "Retrieve success",
			commonFriends: model.CommonFriendRequest{
				Friends: []string{"a@gmail.com", "b@gmail.com"},
			},
			expectedResult: model.FriendListResponse{
				Success: true,
				Friends: []string{"c@gmail.com"},
				Count:   1,
			},
		},
		{
			name: "Users not existed",
			commonFriends: model.CommonFriendRequest{
				Friends: []string{"messi@gmail.com", "ronaldo@gmail.com"},
			},
			expectedResult: model.FriendListResponse{
				Success: true,
				Friends: []string{},
				Count:   0,
			},
			expectedError: errors.New("Sad"),
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/03_common_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.CommonFriends(tt.commonFriends)
			// then
			if err != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestSubscription(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCases := []struct {
		name             string
		subscribeRequest model.SubscriptionRequest
		expectedResult   model.BasicResponse
		expectedError    error
	}{
		{
			name: "Retrieve success",
			subscribeRequest: model.SubscriptionRequest{
				Requestor: "a@gmail.com",
				Target:    "b@gmail.com",
			},
			expectedResult: model.BasicResponse{
				Success: true,
			},
		},
		{
			name: "Users not existed",
			subscribeRequest: model.SubscriptionRequest{
				Requestor: "b@gmail.com",
				Target:    "dang@gmail.com",
			},
			expectedResult: model.BasicResponse{
				Success: true,
			},
			expectedError: errors.New("Sad"),
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/02_connect_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.Subscription(tt.subscribeRequest)
			// then
			if err != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestBlocked(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCases := []struct {
		name           string
		blockedRequest model.SubscriptionRequest
		expectedResult model.BasicResponse
		expectedError  error
	}{
		{
			name: "Retrieve success",
			blockedRequest: model.SubscriptionRequest{
				Requestor: "a@gmail.com",
				Target:    "b@gmail.com",
			},
			expectedResult: model.BasicResponse{
				Success: true,
			},
		},
		{
			name: "Users not existed",
			blockedRequest: model.SubscriptionRequest{
				Requestor: "a@gmail.com",
				Target:    "dang@gmail.com",
			},
			expectedResult: model.BasicResponse{
				Success: true,
			},
			expectedError: errors.New("Sad"),
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/02_connect_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.Blocked(tt.blockedRequest)
			// then
			if err != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestSendUpdate(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCases := []struct {
		name           string
		sendRequest    model.SendUpdateRequest
		expectedResult model.SendUpdateResponse
		expectedError  error
	}{
		{
			name: "Retrieve success",
			sendRequest: model.SendUpdateRequest{
				Sender: "a@gmail.com",
				Text:   "Hello World! b@gmail.com",
			},
			expectedResult: model.SendUpdateResponse{
				Success:    true,
				Recipients: []string{"b@gmail.com"},
			},
		},
		{
			name: "Users not existed",
			sendRequest: model.SendUpdateRequest{
				Sender: "dang@gmail.com",
				Text:   "Hello World! messi@gmail.com",
			},
			expectedResult: model.SendUpdateResponse{
				Success:    true,
				Recipients: []string{},
			},
			expectedError: errors.New("Sad"),
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/03_common_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.SendUpdate(tt.sendRequest)
			// then
			if err != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expectedResult, result)
			}
		})
	}

}

func TestCreateNewUser(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCases := []struct {
		name           string
		sendRequest    model.User
		expectedResult model.BasicResponse
		expectedError  error
	}{
		{
			name: "Retrieve success",
			sendRequest: model.User{
				Email: "c@gmail.com",
			},
			expectedResult: model.BasicResponse{
				Success: true,
			},
		},
		{
			name: "User already existed",
			sendRequest: model.User{
				Email: "a@gmail.com",
			},
			expectedResult: model.BasicResponse{
				Success: true,
			},
			expectedError: &pq.Error{Severity: "ERROR", Code: "23505", Message: "duplicate key value violates unique constraint \"users_pkey\"", Detail: "Key (email)=(a@gmail.com) already exists.", Hint: "", Position: "", InternalPosition: "", InternalQuery: "", Where: "", Schema: "public", Table: "users", Column: "", DataTypeName: "", Constraint: "users_pkey", File: "nbtinsert.c", Line: "432", Routine: "_bt_check_unique"},
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/03_common_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.CreateNewUser(tt.sendRequest)
			// then
			if err != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
