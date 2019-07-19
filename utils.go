package main

import (
	"github.com/hackebrot/turtle"
	"math/rand"
)

var Emojis []string

func UtilsInit() {
	for  _, value := range turtle.Emojis {
		Emojis = append(Emojis, value.String())
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CreateFilename(pattern string) string {
	final := ""
	Numbers := "0123456789"
	Letters := "abcdefghijklmnopqrstuvwxyz"
	for _, c := range pattern {
		switch c {
			case 'e': {
				final += Emojis[rand.Intn(len(Emojis))]
				break
			}
			case 'n': {
				final += string(Numbers[rand.Intn(len(Numbers))])
			}
			case 'c': {
				final += string(Letters[rand.Intn(len(Letters))])
			}
		}
	}
	return final
}
