package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "log"
import "math"
import "os"
import "slices"
import "strings"

var inputPathFlag = flag.String("input_path", "", "Path to puzzle input file")

var digitToValue = map[string]int{
	"zero":  0,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
	"0":     0,
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
}

func findFirstDigit(line string, reverseKey bool) (value int, index int) {
	var firstValue int = -1
	firstDigitIndex := math.MaxInt
	for d, v := range digitToValue {
		digitBytes := []byte(strings.Clone(d))
		if reverseKey {
			slices.Reverse(digitBytes)
		}
		digitIndex := strings.Index(line, string(digitBytes))
		if digitIndex != -1 && digitIndex < firstDigitIndex {
			firstValue = v
			firstDigitIndex = digitIndex
		}
	}
	return firstValue, firstDigitIndex
}

func computeCalibrationValue(inputPath string) int {
	inputFile, err := os.Open(*inputPathFlag)
	if err != nil {
		log.Fatal("Unable to open input file: ", err)
	}
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)
	var calibrationValue int = 0
	for {
		line, _, err := reader.ReadLine()

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Failed to read line from input file: ", err)
		}

		firstDigit, firstIndex := findFirstDigit(string(line), false)
		slices.Reverse(line)
		lastDigit, lastIndex := findFirstDigit(string(line), true)

		if firstIndex == -1 || lastIndex == -1 {
			log.Fatal("Expecting at least one digit per line, found none in: ", string(line))
		}
		calibrationValue += (firstDigit*10 + lastDigit)
	}
	return calibrationValue
}

func main() {
	flag.Parse()

	if *inputPathFlag == "" {
		log.Fatal("Flag --input_path must be non-empty!")
	}
	fmt.Println(computeCalibrationValue(*inputPathFlag))
}
