package db

import (
	"fmt"
	"time"
)

type Store struct {
	exp   time.Time
	value string
}

var store = make(map[string]Store)

func CreateDb() {
	store = make(map[string]Store)
}

func Set(key, value string, exp int) {
	if exp == -1 {
		store[key] = Store{
			value: value,
		}
		return
	}

	store[key] = Store{
		value: value,
		exp:   time.Now().Add(time.Duration(exp) * time.Millisecond),
	}

}

func Get(key string) (string, bool) {
	rec, ok := store[key]

	fmt.Println("READ DATA")
	fmt.Println("READ DATA")
	fmt.Println("READ DATA")
	fmt.Println("IS ZERO")
	fmt.Println(rec.exp.IsZero())
	fmt.Println("time now")
	fmt.Println(time.Now().UnixMilli())
	fmt.Println("expiry")
	fmt.Println(rec.exp.UnixMilli())

	if rec.exp.IsZero() || (time.Now().UnixMilli() < rec.exp.UnixMilli()) {
		return rec.value, ok
	}

	return "", false
}
