package yr

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ToBeAss/funtemps/conv"
)

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConvertTemperature(filepath string) {
	newFile, err := os.Create("kjevik-temp-fahr-20220318-20230318.csv")
	logError(err)
	defer newFile.Close()
	processFile(filepath, "convert", newFile)
}

func CalculateAverage(filepath string, input string) {
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
