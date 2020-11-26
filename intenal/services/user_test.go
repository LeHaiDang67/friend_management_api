package services

import (
	"errors"
	"friend_management/intenal/db"
	"friend_management/intenal/model"
	"friend_management/intenal/util"
	"testing"

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
		expectedResult: []model.User([]model.User{model.User{Email: "c@gmail.com", Friends: []string(nil), Subscription: []string(nil), Blocked: []string(nil)}, model.User{Email: "a@gmail.com", Friends: []string{"b@gmail.com"}, Subscription: []string(nil), Blocked: []string(nil)}, model.User{Email: "b@gmail.com", Friends: []string{"a@gmail.com"}, Subscription: []string(nil), Blocked: []string(nil)}, model.User{Email: "test-email@gmail.com", Friends: []string{"hero@gmail.com"}, Subscription: []string(nil), Blocked: []string(nil)}}),
	}
	//require.NoError(t, util.LoadFixture(db, "./testdata/01_get_user.sql"))
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
			desc:           "Should return no user",
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
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/02_connect_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.ConnectFriends(tt.request)
			// then

			require.Nil(t, err)
			require.Equal(t, tt.expectedResult, result)

		})
	}
}

func TestFriendList(t *testing.T) {
	db := db.InitDatabase()
	defer db.Close()
	testCases := []struct {
		name           string
		givenUserEmail string
		expectedResult model.FriendListResponse
	}{
		{
			name:           "Retrieve success ",
			givenUserEmail: "a@gmail.com",
			expectedResult: model.FriendListResponse{
				Success: true,
				Friends: []string{"b@gmail.com"},
				Count:   1,
			},
		},
		{
			name:           "Retrieve failed ",
			givenUserEmail: "andy@example.com",
			expectedResult: model.FriendListResponse{
				Success: true,
			},
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/01_get_user.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.FriendList(tt.givenUserEmail)
			// then

			require.Nil(t, err)
			require.Equal(t, tt.expectedResult, result)

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
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/03_common_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.CommonFriends(tt.commonFriends)
			// then
			require.Nil(t, err)
			require.Equal(t, tt.expectedResult, result)

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
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/02_connect_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.Subscription(tt.subscribeRequest)
			// then
			require.Nil(t, err)
			require.Equal(t, tt.expectedResult, result)

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
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/02_connect_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.Blocked(tt.blockedRequest)
			// then
			require.Nil(t, err)
			require.Equal(t, tt.expectedResult, result)

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
		},
	}
	require.NoError(t, util.LoadFixture(db, "./testdata/03_common_friend.sql"))
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mn := NewManager(db)
			result, err := mn.SendUpdate(tt.sendRequest)
			// then
			require.Nil(t, err)
			require.Equal(t, tt.expectedResult, result)

		})
	}

}
