package main

import (
	"BYR_Iptables/web"
	"fmt"
)

func main() {
	user1 := web.User{}
	if user1 == (web.User{}) {
		fmt.Println("xxx")
	}
}
