package _generated

import (
	"testing"
)

func TestConvertDataFromAEmbeddedStructToANonEmbeddedStruct(t *testing.T) {
	getUserRequestWithEmbeddedStruct := GetUserRequestWithEmbeddedStruct{
		Common: Common{
			RequestID: 10,
			Token:     "token",
		},
		UserID: 1000,
	}

	bytes, err := getUserRequestWithEmbeddedStruct.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}

	getUserRequest := GetUserRequest{}
	_, err = getUserRequest.UnmarshalMsg(bytes)
	if err != nil {
		t.Fatal(err)
	}
	if getUserRequest.RequestID != getUserRequestWithEmbeddedStruct.RequestID {
		t.Fatal("not same request id")
	}
	if getUserRequest.UserID != getUserRequestWithEmbeddedStruct.UserID {
		t.Fatal("not same user id")
	}
	if getUserRequest.Token != getUserRequestWithEmbeddedStruct.Token {
		t.Fatal("not same token")
	}

	return
}

func TestConvertDataFromANonEmbeddedStructToAEmbeddedStruct(t *testing.T) {
	getUserRequest := GetUserRequest{
		RequestID: 10,
		Token:     "token",
		UserID:    1000,
	}

	bytes, err := getUserRequest.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}

	getUserRequestWithEmbeddedStruct := GetUserRequestWithEmbeddedStruct{}
	_, err = getUserRequestWithEmbeddedStruct.UnmarshalMsg(bytes)
	if err != nil {
		t.Fatal(err)
	}
	if getUserRequest.RequestID != getUserRequestWithEmbeddedStruct.RequestID {
		t.Fatal("not same request id")
	}
	if getUserRequest.UserID != getUserRequestWithEmbeddedStruct.UserID {
		t.Fatal("not same user id")
	}
	if getUserRequest.Token != getUserRequestWithEmbeddedStruct.Token {
		t.Fatal("not same token")
	}

	return
}
