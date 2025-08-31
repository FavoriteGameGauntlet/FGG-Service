package main

import (
	"FavoriteGameGauntlet/user"
	"fmt"
)

func main() {
	user := user.FindUser("Jegern")
	timer := user.GetCurrentTimer()
	fmt.Println(timer.GetTimerId())
}
