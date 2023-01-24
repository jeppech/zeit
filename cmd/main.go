package main

import (
	"fmt"
	"time"

	"github.com/jeppech/zeit"
)

func main() {
	now := time.Now().UTC()
	fmt.Println(now) // 2023-01-24 14:48:58.684023 +0000 UTC

	loc, _ := time.LoadLocation("Europe/Copenhagen")
	z_now := zeit.NowInLoc(loc)
	fmt.Println(z_now) // 14:48:58
}
