 package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type cubeSet struct {
	red   int
	green int
	blue  int
}

func parseCubeSet(s string) (set cubeSet, err error) {
	parts := strings.Split(strings.Trim(s, " "), ",")
	value := -1
	color := ""
	for _, part := range parts {
		_, err := fmt.Sscanf(strings.Trim(part, " "), "%d %s", &value, &color)
		if err != nil {
			return cubeSet{}, err
		}
		switch color {
		case "red":
			set.red = value
		case "green":
			set.green = value
		case "blue":
			set.blue = value
		default:
			return cubeSet{}, fmt.Errorf("Failed to parse unknown color (not RGB): %s", color)
		}
	}
	return set, nil
}

func parseGame(s string) (id int, sets []cubeSet, err error) {
	idAndParts := strings.Split(s, ":")
	if len(idAndParts) != 2 {
		return -1, []cubeSet{}, fmt.Errorf("Expected `Game %d: ..., got: %s", s)
	}
	_, err = fmt.Sscanf(s, "Game %d", &id)
	if err != nil {
		return -1, []cubeSet{}, err
	}

	parts := strings.Split(idAndParts[1], ";")
	sets = make([]cubeSet, len(parts))
	for i, part := range parts {
		sets[i], err = parseCubeSet(part)
		if err != nil {
			return -1, []cubeSet{}, err
		}
	}
	return id, sets, nil
}

func minFeasibleSet(sets []cubeSet) (min cubeSet, err error) {
	if len(sets) == 0 {
		return cubeSet{}, errors.New("minFeasibleSet: Expected at least one cubeSet!")
	}
	min = sets[0]
	for _, set := range sets {
		if set.red > min.red {
			min.red = set.red
		}
		if set.green > min.green {
			min.green = set.green
		}
		if set.blue > min.blue {
			min.blue = set.blue
		}
	}
	return min, nil
}

func computePowerSum(inputPath string, limits cubeSet) int {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatal("Unable to open input file: ", err)
	}
	defer inputFile.Close()

	powerSum := 0
	reader := bufio.NewReader(inputFile)
	for {
		line, _, err := reader.ReadLine()

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Failed to read line from input file: ", err)
		}

		_, cubeSets, err := parseGame(string(line))
		if err != nil {
			log.Fatal("Failed to parse game: ", err)
		}
		min, err := minFeasibleSet(cubeSets)
		if err != nil {
			log.Fatal(err)
		}
		powerSum += (min.red * min.green * min.blue)
	}
	return powerSum
}

func main() {
	inputPathFlag := flag.String("input_path", "", "Path to puzzle input file")
	flag.Parse()
	if *inputPathFlag == "" {
		log.Fatal("Flag --input_path must be non-empty!")
	}

	cubeLimits := cubeSet{red: 12, green: 13, blue: 14}
	fmt.Println(computePowerSum(*inputPathFlag, cubeLimits))
}
