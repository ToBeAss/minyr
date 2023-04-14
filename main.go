package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ToBeAss/funtemps/conv"
)

func main() {
	var input string
	inputHistory := []string{"blank"}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input = scanner.Text()
		if input == "q" || input == "exit" {
			inputHistory = append(inputHistory, input)
			os.Exit(0)
		}
		if inputHistory[len(inputHistory)-1] == "convert" {
			if input == "y" {
				fmt.Println("Genererer filen paa nytt")
				convertTemperature("kjevik-temp-celsius-20220318-20230318.csv")
				inputHistory = append(inputHistory, input)
				fmt.Println("Fil generert")
			} else if input == "n" {
				inputHistory = append(inputHistory, input)
				fmt.Println("Endringer ble ikke lagret")
			} else {
				fmt.Println("Svar y for ja eller n for nei")
			}
		} else if inputHistory[len(inputHistory)-1] == "average" {
			if input == "c" {
				fmt.Println("Printer ut gjennomsnittstemperatur i Celsius")
				calculateAverage("kjevik-temp-celsius-20220318-20230318.csv", input)
				inputHistory = append(inputHistory, input)
			} else if input == "f" {
				fmt.Println("Printer ut gjennomsnittstemperatur i Fahrenheit")
				calculateAverage("kjevik-temp-fahr-20220318-20230318.csv", input)
				inputHistory = append(inputHistory, input)
			} else {
				fmt.Println("Svar c eller f")
			}
		} else if input == "convert" {
			inputHistory = append(inputHistory, input)
			fmt.Println("Konverterer alle maalingene fra Celsius til Fahrenheit")
			if _, err := os.Stat("kjevik-temp-fahr-20220318-20230318.csv"); err == nil {
				logError(err)
				fmt.Println("Filen med Fahrenheit maalinger eksisterer allerede. Vil du erstatte den?")
			} else {
				fmt.Println("Genererer fil")
				convertTemperature("kjevik-temp-celsius-20220318-20230318.csv")
				inputHistory = append(inputHistory, "blank")
				fmt.Println("Fil generert")
			}
		} else if input == "average" {
			inputHistory = append(inputHistory, input)
			fmt.Println("Genererer temperaturgjennomsnittet for hele perioden. Skal det staa i Celsius 'c' eller Fahrenheit 'f'?")
		} else {
			fmt.Println("Venligst skriv convert, average eller exit")
		}
	}
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func convertTemperature(filepath string) {
	newFile, err := os.Create("kjevik-temp-fahr-20220318-20230318.csv")
	logError(err)
	defer newFile.Close()
	processFile(filepath, "convert", newFile)
}

func calculateAverage(filepath string, input string) {
	sum, count := processFile(filepath, "average", nil)
	var average float64
	if count > 0 {
		average = sum / float64(count)
	}
	fmt.Printf("Den gjennomsnittlige temperaturen på Kjevik i perioden er %.2f°%v \n", average, strings.ToUpper(input))
}

func processFile(filepath string, operation string, newFile *os.File) (float64, int) {
	file, err := os.Open(filepath)
	if err != nil {
		logError(err)
	}
	defer file.Close()

	var buffer []byte
	var linebuf []byte
	buffer = make([]byte, 1)

	lineCount := 0
	var temperatureSum float64
	var count int

	for {
		_, err := file.Read(buffer) // Leser av ett og ett tegn
		if err != io.EOF {
			logError(err)
		}

		if buffer[0] != 0x0A { // Frem til tegnet for linjeskift, legges hvert tegn til i linebuf
			linebuf = append(linebuf, buffer[0])
		} else { // Ved linjeskift velges det hva som skal skrives til den nye filen og linebuf nullstilles
			lineCount++

			elementArray := strings.Split(string(linebuf), ";")
			if len(elementArray) > 3 {
				if operation == "convert" {
					newFile.Write(writeToFile(lineCount, elementArray, linebuf))
				} else if operation == "average" {
					temperatureSum, count = sumTemperature(lineCount, elementArray, temperatureSum, count)
				}
			}
			linebuf = nil
		}

		if err == io.EOF {
			break
		}
	}
	return temperatureSum, count
}

func writeToFile(lineCount int, elementArray []string, linebuf []byte) []byte {
	var newLine string
	if lineCount == 1 { // Første linje forblir lik
		newLine = string(linebuf) + "\n"
	} else if len(elementArray[3]) != 0 { // Sjekker om tredje rute har innhold
		celsius, err := strconv.ParseFloat(elementArray[3], 64)
		logError(err)
		fahr := conv.CelsiusToFahrenheit(celsius)
		newLine = elementArray[0] + ";" + elementArray[1] + ";" + elementArray[2] + ";" + fahr + "\n"
	} else { // Siste linje blir erstattet i dette tilfellet
		newLine = "Data er basert paa gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET);endringen er gjort av Tobias Molland;;"
	}
	return []byte(newLine)
}

func sumTemperature(lineCount int, elementArray []string, temperatureSum float64, count int) (float64, int) {
	if lineCount > 1 && elementArray[3] != "" {
		temperature, err := strconv.ParseFloat(elementArray[3], 64) // Deler linjen opp i "ruter" ved semikolon (;)
		if err != nil {
			logError(err)
		}
		temperatureSum += temperature
		count++
	}
	return temperatureSum, count
}
