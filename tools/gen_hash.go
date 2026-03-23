package main

import (
	"fmt"

	"github.com/yourusername/go_server/utils"
)

func main() {
	password := "superadmin123"
	hash, _ := utils.HashPassword(password)
	fmt.Println("Gunakan hash ini di SQL:")
	fmt.Println(hash)
}
