package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"minyr/yr"
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
				yr.ConvertTemperature("kjevik-temp-celsius-20220318-20230318.csv")
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
				yr.CalculateAverage("kjevik-temp-celsius-20220318-20230318.csv", input)
				inputHistory = append(inputHistory, input)
			} else if input == "f" {
				fmt.Println("Printer ut gjennomsnittstemperatur i Fahrenheit")
				yr.CalculateAverage("kjevik-temp-fahr-20220318-20230318.csv", input)
				inputHistory = append(inputHistory, input)
			} else {
				fmt.Println("Svar c eller f")
			}
		} else if input == "convert" {
			inputHistory = append(inputHistory, input)
			fmt.Println("Konverterer alle maalingene fra Celsius til Fahrenheit")
			if _, err := os.Stat("kjevik-temp-fahr-20220318-20230318.csv"); err == nil {
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Filen med Fahrenheit maalinger eksisterer allerede. Vil du erstatte den?")
			} else {
				fmt.Println("Genererer fil")
				yr.ConvertTemperature("kjevik-temp-celsius-20220318-20230318.csv")
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
