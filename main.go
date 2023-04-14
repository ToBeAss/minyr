package main

import (
	"os"
	"bufio"
	"io"
	"log"
	"fmt"
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

	var input string
	var inputHistory []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
    		input = scanner.Text()
    		if input == "q" || input == "exit" {
    			inputHistory = append(inputHistory, input)
			os.Exit(0)
    		}
		if len(inputHistory) > 0 && inputHistory[len(inputHistory)-1] == "convert" {
			if input == "y" {
				inputHistory = append(inputHistory, input)
				fmt.Println("Genererer filen paa nytt")
				ReadFile()
			} else if input == "n" {
				inputHistory = append(inputHistory, input)
				fmt.Println("Endringer ble ikke lagret")
			} else {
				fmt.Println("Svar y for ja eller n for nei")
			}
		} else if len(inputHistory) > 0 && inputHistory[len(inputHistory)-1] == "average" {
			if input == "c" {
    				inputHistory = append(inputHistory, input)
				fmt.Println("Printer ut gjennomsnittstemperatur i Celsius")
			} else if input == "f" {
	    			inputHistory = append(inputHistory, input)
				fmt.Println("Printer ut gjennomsnittstemperatur i Fahrenheit")
			} else {
				fmt.Println("Svar c eller f")
			}
		} else if input == "convert" {
			inputHistory = append(inputHistory, input)
			fmt.Println("Konverterer alle maalingene fra Celsius til Fahrenheit")
			if _, err := os.Stat("kjevik-temp-fahr-20220318-20230318.csv"); err == nil {
				fmt.Println("Filen med Fahrenheit maalinger eksisterer allerede. Vil du erstatte den?")
			} else {
				fmt.Println("Genererer fil")
				ReadFile()
			}
		} else if input == "average" {
			inputHistory = append(inputHistory, input)
			fmt.Println("Genererer temperaturgjennomsnittet for hele perioden. Skal det staa i Celsius 'c' eller Fahrenheit 'f'?")
    		} else {
        		fmt.Println("Venligst skriv convert, average eller exit")
    		}
	}

	var buffer []byte
	var linebuf []byte
	var newLine string
	buffer = make([]byte, 1)

	bytesCount := 0
	linesCount := 0

func ReadFile() {
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
}
