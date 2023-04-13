package main

import (
	"os"
	"log"
	"io"
	"strings"
	"strconv"
	"github.com/ToBeAss/funtemps/conv"
)

func main(){
	src, err := os.Open("kjevik-temp-celsius-20220318-20230318.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	newFile, err := os.Create("kjevik-temp-fahr-20220318-20230318.csv")
	if err != nil {
    		log.Fatal(err)
	}
	defer newFile.Close()

	var buffer []byte
	var linebuf []byte
	var newLine string
	buffer = make([]byte, 1)

	bytesCount := 0
	linesCount := 0
	for {
		_, err := src.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		bytesCount++
		if buffer[0] == 0x0A {
			linesCount++;
			elementArray := strings.Split(string(linebuf), ";")
			if len(elementArray) > 3 {
				if linesCount > 1 {
					if len(elementArray[3]) == 0 {
						newLine = "Data er basert paa gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET);endringen er gjort av Tobias Molland;;"
						newFile.Write([]byte(newLine))
					} else {
						celsius, err := strconv.ParseFloat(elementArray[3], 64)
						if err != nil {
							log.Fatal(err)
						}
							fahr := conv.CelsiusToFahrenheit(celsius)
							newLine = elementArray[0] + ";" + elementArray[1] + ";" + elementArray[2] + ";" + fahr + "\n"
							newFile.Write([]byte(newLine))
					}
				} else {
					newLine = string(linebuf) + "\n"
					newFile.Write([]byte(newLine))
				}
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
