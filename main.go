package main

import (
	"os"
	"log"
)

func main(){
	src, err := os.Open("table.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()
	log.Println(src)
	
	var buffer []byte
	var linebuf []byte
	buffer = make([]byte, 1)

	bytesCount := 0
}
