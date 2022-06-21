package avchem_server

import (
	"log"
	"testing"
	"time"
)

func TestTest(x *testing.T) {

	t, err := time.Parse("15:04:05", "10:20:33")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(t.Hour(), t.Minute(), t.Second())

}
