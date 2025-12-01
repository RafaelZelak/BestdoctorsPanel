package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin"
	// Usando Cost 14 como definido no projeto
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	fmt.Println(string(hash))
}
