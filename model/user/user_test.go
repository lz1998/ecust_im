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
