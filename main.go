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
				processFiles(inputHistory)
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
				//fmt.Printf("%.2f", calculateAverage())
				inputHistory = append(inputHistory, input)
			} else if input == "f" {
				fmt.Println("Printer ut gjennomsnittstemperatur i Fahrenheit")
				//fmt.Printf("%.2f", calculateAverage())
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
				processFiles(inputHistory)
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

func processFiles(inputHistory []string) {
	src, err := os.Open("kjevik-temp-celsius-20220318-20230318.csv")
	logError(err)
	defer src.Close()

	newFile, err := os.Create("kjevik-temp-fahr-20220318-20230318.csv")
	logError(err)
	defer newFile.Close()

	var buffer []byte
	var linebuf []byte
	buffer = make([]byte, 1)

	byteCount := 0
	lineCount := 0
	var temperatureSlice []float64

	for {
		_, err := src.Read(buffer) // Leser av ett og ett tegn
		if err != io.EOF {
			logError(err)
		}
		byteCount++

		if buffer[0] != 0x0A { // Frem til tegnet for linjeskift, legges hvert tegn til i linebuf
			linebuf = append(linebuf, buffer[0])
		} else { // Ved linjeskift velges det hva som skal skrives til den nye filen og linebuf nullstilles
			lineCount++

			elementArray := strings.Split(string(linebuf), ";") // Deler linjen opp i "ruter" ved semikolon (;)
			if len(elementArray) > 3 {
				newFile.Write(writeToFile(lineCount, elementArray, linebuf))
				getAllTemperatures(lineCount, elementArray, temperatureSlice)
			}
			linebuf = nil
		}

		if err == io.EOF {
			break
		}
	}
	if inputHistory[len(inputHistory)-1] == "average" {
		sum := .0
		for i := 0; i < len(temperatureSlice); i++ {
			sum += (temperatureSlice[i])
		}
		average := sum / float64(len(temperatureSlice))
		fmt.Printf("Gjennomsnittstemperaturen for perioden var %.2f grader Celsius", average)
	}
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

func getAllTemperatures(lineCount int, elementArray []string, temperatureSlice []float64) {
	log.Println("getAllTemperatures")
	if lineCount > 1 && len(elementArray[3]) > 0 {
		log.Println("Past logic")
		float, err := strconv.ParseFloat(elementArray[3], 64)
		if err != nil {
			logError(err)
		}
		log.Println(float)
		temperatureSlice = append(temperatureSlice, float)
	}
}

/*func calculateAverage() float64 {
	sum := .0
	for i := 0; i < len(temperatureSlice); i++ {
		sum += (temperatureSlice[i])
	}
	average := sum / float64(len(temperatureSlice))
	//log.Println(temperatureSlice[0])
	return average
}*/
