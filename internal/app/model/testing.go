package model

import "testing"

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User {
		Name: "Alex111",
		Balance: 500,
	}
}