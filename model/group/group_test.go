package group

import "testing"

func TestCreateGroup(t *testing.T) {
	group, err := CreateGroup(&EcustGroup{
		GroupName: "group_test",
		OwnerId:   10000,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", group)
	}
}

func TestListGroup(t *testing.T) {
	groups, err := ListGroup([]int64{1})
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", groups)
	}
}

func TestUpdateGroup(t *testing.T) {
	// TODO
}
