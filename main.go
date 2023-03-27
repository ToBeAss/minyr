package main

import (
	"os"
	"log"
	"io"
	"strings"
	"github.com/ToBeAss/funtemps/conv"
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
	for {
		_, err := src.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		bytesCount++
		if buffer[0] == 0x0A {
			log.Println(string(linebuf))
			// Her
			elementArray := strings.Split(string(linebuf), ";")
			if len(elementArray) > 3 {
				celsius := elementArray[3]
				fahr := conv.CelsiusToFahrenheit(celsius)
				log.Println(elementArray[3])
			}
			linebuf = nil
		} else {
			linebuf = append(linebuf, buffer[0])
		}
		if err == io.EOF {
			break
		}
	}
}
