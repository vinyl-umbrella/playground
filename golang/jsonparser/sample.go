package main

import (
	"fmt"

	"github.com/buger/jsonparser"
)

type User struct {
	Username string `json:"username"`
}

func main() {
	j := []byte(`{"username": "a", "username": "b"}`)
	username, _ := jsonparser.GetString(j, "username")
	fmt.Println(username) // a
}
