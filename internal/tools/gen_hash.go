package main

import (
	"fmt"

	"github.com/yourusername/go_server/internal/utils"
)

func main() {
	password := "superadmin123"
	hash, err := utils.HashPassword(password)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Gunakan hash ini di SQL:\n%s\n", hash)
}
