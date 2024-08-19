package main

import (
	"crypto/rand"
	"fmt"
	"strings"
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

	// JWT none Attack.
	splittedToken1 := strings.Split(token1, ".")
	// {"alg": "none","typ": "JWT"} => eyJhbGciOiAibm9uZSIsInR5cCI6ICJKV1QifQ
	token1NoneAlg := "eyJhbGciOiAibm9uZSIsInR5cCI6ICJKV1QifQ." + splittedToken1[1] + "." + splittedToken1[2]

	verifyResult1NoneAlg, _ := verify(token1NoneAlg, key)
	fmt.Println(user1, token1NoneAlg, verifyResult1NoneAlg)

	// Sign and verify for user2, but token expired
	user2 := "user2"
	token2, _ := sign(user2, key)
	time.Sleep(time.Second * 6)
	verifyResult2, _ := verify(token2, key)
	fmt.Println(user2, token2, verifyResult2)
}

// Generate a 32 bit secret key
func initKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Create a new token
func sign(userid string, secretKey []byte) (string, error) {
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

// Verify the token
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
