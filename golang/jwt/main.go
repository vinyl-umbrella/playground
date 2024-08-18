package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Initialize secret key
	key, _ := initKey()
	fmt.Println(key)

	// Sign and verify for user1
	user1 := "user1"
	token1, _ := sign(user1, key)
	verifyResult1, _ := verify(token1, key)
	fmt.Println(user1, token1, verifyResult1)

	// Sign and verify for user2, but token expired
	user2 := "user2"
	token2, _ := sign(user2, key)
	time.Sleep(time.Second * 6)
	verifyResult2, _ := verify(token2, key)
	fmt.Println(user2, token2, verifyResult2)
}

func initKey() ([]byte, error) {
	// Generate a 32 bit secret key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func sign(userid string, secretKey []byte) (string, error) {
	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set some claims
	token.Claims = jwt.MapClaims{
		"sub": userid,
		"exp": time.Now().Add(time.Second * 5).Unix(),
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verify(tokenString string, secretKey []byte) (bool, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// Validate the token
	if token.Valid {
		return true, nil
	}

	return false, nil
}
