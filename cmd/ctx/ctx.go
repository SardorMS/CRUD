package main

import (
	"context"	
	"log"
)


func main() {
	root := context.Background()
	inner := context.WithValue(root, "key", "inner")
	outer := context.WithValue(root, "key", "outer")

	log.Println(outer.Value("key"))
	log.Println(inner.Value("inner-key"))
	log.Println(outer.Value("outer-key"))



}