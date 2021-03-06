package services

import (
	"database/sql"
	"errors"
	"fmt"
	"friend_management/intenal/db"
	"friend_management/intenal/model"
	"friend_management/intenal/util"

	"log"
	"strings"
)

//IUserService is a interface
type IUserService interface {
	ConnectFriends(req model.FriendConnectionRequest) (model.BasicResponse, error)
	GetUser(email string) (model.User, error)
	SendUpdate(sendRequest model.SendUpdateRequest) (model.SendUpdateResponse, error)
	GetAllUsers() ([]model.User, error)
	Blocked(subRequest model.SubscriptionRequest) (model.BasicResponse, error)
	Subscription(subRequest model.SubscriptionRequest) (model.BasicResponse, error)
	CommonFriends(commonFriends model.CommonFriendRequest) (model.FriendListResponse, error)
	FriendList(email model.FriendListRequest) (model.FriendListResponse, error)
	CreateNewUser(user model.User) (model.BasicResponse, error)
}

//Store is struct to implement interface
type Store struct {
	dbconn *sql.DB
}

//NewManager change sql.DB to Store
func NewManager(dbconn *sql.DB) *Store {
	return &Store{
		dbconn: dbconn,
	}
}

//ConnectFriends that func connect 2 user
func (st *Store) ConnectFriends(req model.FriendConnectionRequest) (model.BasicResponse, error) {
	basicResponse := model.BasicResponse{}
	//var service IUserService
	userA, errA := db.GetTheUser(st.dbconn, req.Friends[0])
	if errA != nil {
		fmt.Printf("Error QueryA: %s\n", errA)
		basicResponse.Success = false
		return basicResponse, errA
	}
	checkUserA, errCheckUserA := CheckUserExist(st.dbconn, userA.Email)
	if errCheckUserA != nil {
		basicResponse.Success = false
		return basicResponse, errCheckUserA
	}
	if !checkUserA {
		basicResponse.Success = false
		return basicResponse, nil
	}
	userB, errB := db.GetTheUser(st.dbconn, req.Friends[1])
	if errB != nil {
		fmt.Printf("Error QueryB: %s\n", errB)
		basicResponse.Success = false
		return basicResponse, errB
	}
	checkUserB, errCheckUserB := CheckUserExist(st.dbconn, userB.Email)
	if errCheckUserB != nil {
		basicResponse.Success = false
		return basicResponse, errCheckUserB
	}
	if !checkUserB {
		basicResponse.Success = false
		return basicResponse, nil
	}
	bBlock := util.Contains(userA.Blocked, userB.Email)
	aBlock := util.Contains(userB.Blocked, userA.Email)
	if aBlock || bBlock {
		basicResponse.Success = false
		return basicResponse, nil
	}

	bFriend := util.Contains(userA.Friends, userB.Email)
	aFriend := util.Contains(userB.Friends, userA.Email)
	if !bFriend || !aFriend {
		errUpdateA := AddFriends(st.dbconn, userB.Email, userA.Email)
		if errUpdateA != nil {
			fmt.Printf("Error QueryA: %s\n", errUpdateA)
		}
		log.Printf("B added to A friend's\n")
		errUpdateB := AddFriends(st.dbconn, userA.Email, userB.Email)
		if errUpdateB != nil {
			fmt.Printf("Error QueryB: %s\n", errUpdateB)
		}
		log.Printf("A added to B friend's\n")
	}

	basicResponse.Success = true
	return basicResponse, nil
}

//FriendList show friend list
func (st *Store) FriendList(email model.FriendListRequest) (model.FriendListResponse, error) {
	var friendList model.FriendListResponse
	user, err := db.GetUserByEmail(st.dbconn, email)
	if err != nil {
		friendList.Success = false
		return friendList, nil
	}
	checkUser, errCheckUser := CheckUserExist(st.dbconn, user.Email)
	if errCheckUser != nil {
		friendList.Success = false
		return friendList, errCheckUser
	}
	if !checkUser {
		friendList.Success = false
		return friendList, nil
	}
	friendList.Success = true
	friendList.Friends = user.Friends
	friendList.Count = len(user.Friends)
	return friendList, nil
}

//CommonFriends retrieve the common friends list between two email addresses
func (st *Store) CommonFriends(commonFriends model.CommonFriendRequest) (model.FriendListResponse, error) {
	var friendList model.FriendListResponse
	userA, errA := db.GetTheUser(st.dbconn, commonFriends.Friends[0])
	if errA != nil {
		fmt.Printf("Error QueryA: %s\n", errA)
		friendList.Success = false
		return friendList, errA
	}
	checkUserA, errCheckUserA := CheckUserExist(st.dbconn, userA.Email)
	if errCheckUserA != nil {
		friendList.Success = false
		return friendList, errCheckUserA
	}
	if !checkUserA {
		friendList.Success = false
		return friendList, nil
	}
	userB, errB := db.GetTheUser(st.dbconn, commonFriends.Friends[1])
	if errB != nil {
		fmt.Printf("Error QueryB: %s\n", errB)
		friendList.Success = false
		return friendList, errB
	}
	checkUserB, errCheckUserB := CheckUserExist(st.dbconn, userB.Email)
	if errCheckUserB != nil {
		friendList.Success = false
		return friendList, errCheckUserB
	}
	if !checkUserB {
		friendList.Success = false
		return friendList, nil
	}
	Commons := []string{}
	for _, a := range userA.Friends {
		for _, b := range userB.Friends {
			if a == b {
				Commons = append(Commons, a)
			}
		}
	}
	friendList.Success = true
	friendList.Friends = Commons
	friendList.Count = len(Commons)
	return friendList, nil
}

//Subscription subscribe to updates from an email address.
func (st *Store) Subscription(subRequest model.SubscriptionRequest) (model.BasicResponse, error) {
	var basicResponse model.BasicResponse
	userRequestor, errGetRequestor := db.GetTheUser(st.dbconn, subRequest.Requestor)
	if errGetRequestor != nil {
		basicResponse.Success = false
		return basicResponse, errGetRequestor
	}
	checkRequestor, errCheckRequestor := CheckUserExist(st.dbconn, userRequestor.Email)
	if errCheckRequestor != nil {
		basicResponse.Success = false
		return basicResponse, errCheckRequestor
	}
	if !checkRequestor {
		basicResponse.Success = false
		return basicResponse, nil
	}

	userTarget, errGetTarget := db.GetTheUser(st.dbconn, subRequest.Target)
	if errGetTarget != nil {
		basicResponse.Success = false
		return basicResponse, errGetTarget
	}
	checkTarget, errCheckTarget := CheckUserExist(st.dbconn, userTarget.Email)
	if errCheckTarget != nil {
		basicResponse.Success = false
		return basicResponse, errCheckTarget
	}
	if !checkTarget {
		basicResponse.Success = false
		return basicResponse, nil
	}

	isUserRequestor := util.Contains(userRequestor.Subscription, userTarget.Email)
	if !isUserRequestor {
		err := db.CreateSubscribeFriend(st.dbconn, userRequestor.Email, userTarget.Email)
		if err != nil {
			basicResponse.Success = false
			return basicResponse, err
		}
	}

	basicResponse.Success = true
	return basicResponse, nil
}

//Blocked is  an API to block updates from an email address
func (st *Store) Blocked(subRequest model.SubscriptionRequest) (model.BasicResponse, error) {
	var basicResponse model.BasicResponse
	userRequestor, errGetRequestor := db.GetTheUser(st.dbconn, subRequest.Requestor)
	if errGetRequestor != nil {
		basicResponse.Success = false
		return basicResponse, errGetRequestor
	}
	checkRequestor, errCheckRequestor := CheckUserExist(st.dbconn, userRequestor.Email)
	if errCheckRequestor != nil {
		basicResponse.Success = false
		return basicResponse, errCheckRequestor
	}
	if !checkRequestor {
		basicResponse.Success = false
		return basicResponse, nil
	}
	userTarget, errGetTarget := db.GetTheUser(st.dbconn, subRequest.Target)
	if errGetTarget != nil {
		basicResponse.Success = false
		return basicResponse, errGetTarget
	}
	checkTarget, errCheckTarget := CheckUserExist(st.dbconn, userTarget.Email)
	if errCheckTarget != nil {
		basicResponse.Success = false
		return basicResponse, errCheckTarget
	}
	if !checkTarget {
		basicResponse.Success = false
		return basicResponse, nil
	}
	isUserRequestor := util.Contains(userRequestor.Blocked, userTarget.Email)
	if !isUserRequestor {
		err := db.CreateBlockedFriend(st.dbconn, userRequestor.Email, userTarget.Email)
		if err != nil {
			basicResponse.Success = false
			return basicResponse, err
		}
	}

	basicResponse.Success = true
	return basicResponse, nil

}

//SendUpdate retrieve all email addresses that can receive updates from an email address.
func (st *Store) SendUpdate(sendRequest model.SendUpdateRequest) (model.SendUpdateResponse, error) {
	var sendResponse model.SendUpdateResponse
	sender, errGetSender := db.GetTheUser(st.dbconn, sendRequest.Sender)
	if errGetSender != nil {
		sendResponse.Success = false
		return sendResponse, errGetSender
	}
	checkSender, errCheckSender := CheckUserExist(st.dbconn, sendRequest.Sender)
	if errCheckSender != nil {
		sendResponse.Success = false
		return sendResponse, errCheckSender
	}
	if !checkSender {
		sendResponse.Success = false
		return sendResponse, nil
	}
	Recipients := []string{}
	allUser, errGetAllUser := db.GetListUsers(st.dbconn)
	if errGetAllUser != nil {
		sendResponse.Success = false
		return sendResponse, errGetAllUser
	}
	for _, u := range allUser {
		var isBlock = util.Contains(u.Blocked, sender.Email)
		if !isBlock {
			isFriend := util.Contains(u.Friends, sender.Email)
			isSubscriber := util.Contains(u.Subscription, sender.Email)
			isMentioned := strings.Contains(sendRequest.Text, u.Email)
			if isFriend || isSubscriber || isMentioned {
				Recipients = append(Recipients, u.Email)
			}

		}
	}
	sendResponse.Success = true
	sendResponse.Recipients = Recipients
	return sendResponse, nil
}

//GetUser get user bu email
func (st *Store) GetUser(email string) (model.User, error) {
	user, err := db.GetTheUser(st.dbconn, email)
	if err != nil {
		return user, err
	}
	checkUser, errCheckUser := CheckUserExist(st.dbconn, user.Email)
	if errCheckUser != nil {
		return user, errCheckUser
	}
	if !checkUser {
		return user, errors.New("Cannot fetch user")
	}
	return user, nil
}

//GetAllUsers get all user
func (st *Store) GetAllUsers() ([]model.User, error) {
	users, err := db.GetListUsers(st.dbconn)
	if err != nil {
		return users, err
	}
	return users, nil
}

//AddFriends add a new friend
func AddFriends(db *sql.DB, emailFriend string, email string) error {

	result, err := db.Exec("Update users set friends=array_append(friends,$1)  where email = $2 ",
		emailFriend, email)
	if err != nil {
		return err
	}

	result.RowsAffected()
	return nil
}

//CreateNewUser creates a new user
func (st *Store) CreateNewUser(user model.User) (model.BasicResponse, error) {
	var res model.BasicResponse
	checkUser, errCheckUser := CheckUserExist(st.dbconn, user.Email)
	if errCheckUser != nil {
		res.Success = false
		return res, errCheckUser
	}
	if checkUser {
		res.Success = false
		return res, nil
	}
	err := db.CreateNewUser(st.dbconn, user)
	if err != nil {
		res.Success = false
		return res, err
	}
	res.Success = true
	return res, nil
}

//CheckUserExist will check the user that is exist or not
func CheckUserExist(dbconn *sql.DB, email string) (bool, error) {
	var count int
	err := dbconn.QueryRow("select count(*) from users where email = $1", email).Scan(&count)
	if err != nil {
		log.Println("Query error: ", err.Error())
		return false, err
	}
	if count < 1 {
		return false, nil
	}
	return true, nil
}
