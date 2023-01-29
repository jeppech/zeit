package main

import (
	"os"
	"os/signal"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/jeppech/zeit"
)

func main() {

	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	for i := 0; i < 1000; i++ {
		t1 := zeit.Now()
		t2, _ := zeit.Parse("23:59:59")
		zeit.RangeFromZeit(t1, t2)

		loc, _ := time.LoadLocation("Europe/Copenhagen")
		t3 := zeit.NowInLoc(loc)
		t4 := t3.Add(8 * time.Hour)
		zeit.RangeFromZeit(t3, t4)

		r, _ := zeit.ParseRange("08:00:00", "17:00:00")

		r.Split(30 * time.Minute)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
