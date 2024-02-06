package main

import (
	"fmt"
	"log"
)

// A type intended as a *value* value of the "interval" map
type interval struct {
	LeftBound  float64
	RightBound float64
	Length     float64
}

// A function that returns a list of keys of a given map
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// Main Func
func main() {
	var stopSignCode float64
	var stringToEncode string

	// A map, representing each symbol probability in a message given to encode
	intervalMap := make(map[string]interval)

	fmt.Println(`Give me a string to encode (don't forget to end it with "!" stop sign"): `)
	_, err := fmt.Scanln(&stringToEncode)
	if err != nil {
		log.Fatal(err)
	}
	//stringToEncode = "abaabaaca!"
	fmt.Printf("Encoding: %v\n", stringToEncode)
	stopSignCode, intervalMap = Encode(stringToEncode)
	result, err := Decode(stopSignCode, intervalMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoded message is: %v\n", result)
}

// A function to encode a given string, returning the stop sign code, and a map needed to decode the string
func Encode(inputStr string) (float64, map[string]interval) {

	// First step: lets come up with each symbol count in the given string
	strLen := len(inputStr)
	symbolCountMap := make(map[string]int)
	for i := 0; i < strLen; i++ {
		symbolCountMap[string(inputStr[i])]++
	}

	// Taking all the symbols the string contains
	strKeys := Keys(symbolCountMap)

	// Creating a new map which keys are symbols and values are intervals of the [0; 1] line segment
	intervalMap := make(map[string]interval)
	var PreviousRightBound float64 = 0.0
	for i := range strKeys {
		var newMapValue interval
		newMapValue.LeftBound = PreviousRightBound
		newMapValue.RightBound = newMapValue.LeftBound + float64(symbolCountMap[strKeys[i]])/10
		newMapValue.Length = newMapValue.RightBound - newMapValue.LeftBound
		intervalMap[strKeys[i]] = newMapValue
		PreviousRightBound = newMapValue.RightBound
	}
	//fmt.Printf("INTERVAL MAP: %v\n", intervalMap)

	// Consideting we will need out intervalMap to divige the line segment again and again, let's create a copy that we can rewrite on each step
	tempIterMap := make(map[string]interval)
	for k, v := range intervalMap {
		tempIterMap[k] = v
	}
	//fmt.Printf("INTERVAL MAP: %v[COPY]\n\n", tempIterMap)

	// Here comes the music.
	for i := 0; i < strLen; i++ {

		// On each iteration we take the next symbol in the string and make it boundaries the new boundaries of our whole line segment
		currentSymbol := string(inputStr[i])
		fmt.Printf("\nCurrent Symbol: %v(%v), [%v of %v]    ", currentSymbol, intervalMap[currentSymbol].Length, i, strLen)

		var nextInterval interval
		nextInterval.LeftBound = tempIterMap[currentSymbol].LeftBound
		nextInterval.RightBound = tempIterMap[currentSymbol].RightBound
		nextInterval.Length = nextInterval.RightBound - nextInterval.LeftBound
		fmt.Printf("Next Interval: %v\n", nextInterval)

		// Then we allocate new intervals of each letter according to the ratio of the character interval to the segment(that became different after the previous step)
		PreviousRightBound = nextInterval.LeftBound
		for j := range strKeys {
			symbolInterval := intervalMap[strKeys[j]].Length

			var newMapValue interval
			newMapValue.LeftBound = PreviousRightBound
			newMapValue.RightBound = newMapValue.LeftBound + float64(nextInterval.Length*symbolInterval)
			newMapValue.Length = newMapValue.RightBound - newMapValue.LeftBound
			PreviousRightBound = newMapValue.RightBound

			tempIterMap[strKeys[j]] = newMapValue

			fmt.Printf("new interval for %v is %v\n", strKeys[j], newMapValue)
		}
		fmt.Printf("%v\n", tempIterMap)
	}

	// When we approach the last symbol of our string, it is the stop symbol, that is meant to be returned atfer using Encode function
	stopSignCode := (tempIterMap["!"].RightBound + tempIterMap["!"].LeftBound) / 2
	fmt.Printf("\nStop Sign Code: %v\n", stopSignCode)

	return stopSignCode, intervalMap
}

// A function, taking the stop sign code and the interval map needed to decode a string
func Decode(stopSignCode float64, intervalMap map[string]interval) (string, error) {
	var outText string

	// Taking all the symbols the string contains
	strKeys := Keys(intervalMap)

	//fmt.Printf("INTERVAL MAP: %v\n", intervalMap)

	// Consideting we will need out intervalMap to divige the line segment again and again, let's create a copy that we can rewrite on each step
	tempIterMap := make(map[string]interval)
	for k, v := range intervalMap {
		tempIterMap[k] = v
	}
	//fmt.Printf("INTERVAL MAP: %v[COPY]\n\n", tempIterMap)

	// The cycle will stop when it decodes the stop sign
	for {
		var currentSymbol string

		// Loop to decode the current symbol on each iteration
		for d := range tempIterMap {
			if stopSignCode >= tempIterMap[d].LeftBound && stopSignCode <= tempIterMap[d].RightBound {
				currentSymbol = d
				break
			}
		}
		outText += currentSymbol
		if currentSymbol == "!" {
			return outText, nil
		}

		// On each iteration we take the next symbol in the string and make it boundaries the new boundaries of our whole line segment
		var nextInterval interval
		nextInterval.LeftBound = tempIterMap[currentSymbol].LeftBound
		nextInterval.RightBound = tempIterMap[currentSymbol].RightBound
		nextInterval.Length = nextInterval.RightBound - nextInterval.LeftBound
		PreviousRightBound := nextInterval.LeftBound
		fmt.Printf("\nDecoded Symbol: %v", currentSymbol)
		fmt.Printf("Next Interval: %v\n", nextInterval)

		// Then we allocate new intervals of each letter according to the ratio of the character interval to the segment(that became different after the previous step)
		for j := range strKeys {
			symbolInterval := intervalMap[strKeys[j]].Length

			var newMapValue interval
			newMapValue.LeftBound = PreviousRightBound
			// fmt.Printf("will add %v", float64(nextInterval.Length*symbolInterval))
			newMapValue.RightBound = newMapValue.LeftBound + float64(nextInterval.Length*symbolInterval)
			newMapValue.Length = newMapValue.RightBound - newMapValue.LeftBound

			fmt.Printf("new interval for %v is %v\n", strKeys[j], newMapValue)
			tempIterMap[strKeys[j]] = newMapValue
			PreviousRightBound = newMapValue.RightBound
		}
		fmt.Printf("%v\n", tempIterMap)
	}
}
