package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// TODO: Reuse this from common library instead of copying from day 4.
func parseInts(s string) ([]int, error) {
	ints := make([]int, 0, 10) // arbitrary capacity
	tokens := strings.Split(strings.Trim(s, " "), " ")
	for _, t := range tokens {
		if t == "" {
			continue
		}
		i, err := strconv.Atoi(t)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func scanIntsWithPrefix(scanner *bufio.Scanner, prefix string) ([]int, error) {
	if !scanner.Scan() {
		return nil, errors.New("Expected another token!")
	}
	token := scanner.Text()
	if strings.Index(token, prefix) != 0 {
		return nil, fmt.Errorf("Expected token with prefix `%s`, got: %s", prefix, token)
	}
	return parseInts(token[len(prefix):])
}

// Computes the number of viable strategies for covering
// `distMillim` in strictly less than `timeMillis`.
//
// The racer's initial velocity is zero.
// At the start of the race, the racer can choose to either
// 1) keep holding the pedal and not moving, while gaining
// 1 millim / millis of speed per millis, or 2) release the
// pedal and continue moving at the obtained velocity.
// Viable strategies are all combinations of (1) and (2)
// that cover `distMillim` in less than `timeMilis`.
func numViableStrategies(timeMillis int, distMillim int) int {
	// x := number of millis the pedal is held
	// t := `timeMillis`
	// d := `distMillim
	// viable strategies: x * (t - x) > d > 0

	// If the quadratic equation has less than two real solutions,
	// there are no viable strategies.
	d := timeMillis * timeMillis - 4 * distMillim
	if d <= 0 {
		return 0
	}

	sqrtd := math.Sqrt(float64(d))
	leftBound := (float64(timeMillis) - sqrtd) / 2.0	
	rightBound := (float64(timeMillis) + sqrtd) / 2.0
	minViable := int(math.Floor(leftBound)) + 1
	maxViable := int(math.Ceil(rightBound)) - 1

	if minViable > maxViable {
		return 0
	}
	return maxViable - minViable + 1
}

func computeStrategyProduct(inputPath string) (int, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return -1, err
	}
	defer inputFile.Close()
	
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)

	timesMillis, err := scanIntsWithPrefix(scanner, "Time:")
	if err != nil {
		return -1, err
	}
	distancesMillim, err := scanIntsWithPrefix(scanner, "Distance:")
	if err != nil {
		return -1, err
	}

	if len(timesMillis) != len(distancesMillim) {
		return -1, errors.New("Different number of times and distances!")
	}

	p := 1
	for i := 0; i < len(timesMillis); i++ {
		p *= numViableStrategies(timesMillis[i], distancesMillim[i])
	}
	return p, nil
}

func main() {
	inputPathFlag := flag.String("input_path", "", "Path to puzzle input file")
	flag.Parse()

	if *inputPathFlag == "" {
		log.Fatal("Flag --input_path must be non-empty!")
	}

	points, err := computeStrategyProduct(*inputPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(points)

}
