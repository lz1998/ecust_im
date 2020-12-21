package user

import "testing"

func TestCreateUser(t *testing.T) {
	user, err := CreateUser("123", "hello")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", user)
	}
}

func TestListUser(t *testing.T) {
	users, err := ListUser([]int64{10000})
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", users)
	}
}

func TestUpdateUser(t *testing.T) {
	users, err := ListUser([]int64{10000})
	if err != nil {
		t.Error(err)
		return
	}
	if len(users) < 1 {
		return
	}
	user := users[0]
	user.Password = "aaa"
	user.Nickname = "asdf"
	if err := UpdateUser([]*EcustUser{user}); err != nil {
		t.Error(err)
	}
}
