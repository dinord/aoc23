package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type irange struct {
	start int // included
	end   int // not included
}

type rangeMap struct {
	src irange
	dst irange
}

type categoryMap struct {
	src       string
	dst       string
	rangeMaps []rangeMap
}

type puzzle struct {
	seeds   []irange
	srcMaps map[string]categoryMap
}

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

func parseSeeds(line string) (seeds []irange, err error) {
	const Prefix = "seeds: "
	if strings.Index(line, Prefix) != 0 {
		return nil, fmt.Errorf("Expected prefix `%s`, got: %s", Prefix, line)
	}
	numbers, err := parseInts(line[len(Prefix):])
	if err != nil {
		return nil, err
	}

	count := len(numbers)
	if count%2 != 0 {
		return nil, fmt.Errorf("Expected start-length pairs of seed locations, got: %s", line)
	}

	seeds = make([]irange, count/2)
	for i, _ := range seeds {
		seeds[i].start = numbers[2*i]
		seeds[i].end = seeds[i].start + numbers[2*i+1]
	}
	return seeds, nil
}

func parseSrcDst(line string) (src string, dst string, err error) {
	var mapName string
	_, err = fmt.Sscanf(line, "%s map:", &mapName)
	if err != nil {
		return "", "", err
	}

	srcDst := strings.Split(mapName, "-to-")
	if len(srcDst) != 2 {
		return "", "", fmt.Errorf("Expected `key-to-value`, got: %s", line)
	}
	return srcDst[0], srcDst[1], err
}

func parseRangeMap(line string) (r rangeMap, err error) {
	ints, err := parseInts(line)
	if err != nil {
		return r, err
	}
	if len(ints) != 3 {
		return r, fmt.Errorf("Expected `<d> <s> <l>`. Got: %s", line)
	}
	length := ints[2]
	src := irange{start: ints[1], end: ints[1] + length}
	dst := irange{start: ints[0], end: ints[0] + length}
	return rangeMap{src: src, dst: dst}, nil
}

func scanRangeMaps(scanner *bufio.Scanner) (rs []rangeMap, err error) {
	rs = make([]rangeMap, 0, 10) // arbitrary capacity
	for scanner.Scan() && scanner.Text() != "" {
		r, err := parseRangeMap(scanner.Text())
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return rs, nil
}

func loadPuzzle(inputPath string) (p puzzle, err error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return p, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)

	if scanner.Scan() {
		p.seeds, err = parseSeeds(scanner.Text())
	} else {
		return p, fmt.Errorf("Expected line with seeds!")
	}
	if err != nil {
		return p, err
	}

	p.srcMaps = make(map[string]categoryMap)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		var cm categoryMap
		cm.src, cm.dst, err = parseSrcDst(scanner.Text())
		if err != nil {
			return p, err
		}
		cm.rangeMaps, err = scanRangeMaps(scanner)
		if err != nil {
			return p, err
		}
		p.srcMaps[cm.src] = cm

	}

	if err := scanner.Err(); err != nil {
		return p, err
	}
	return p, nil
}

func (r irange) apply(m rangeMap) (unmapped []irange, mapped []irange) {
	if m.src.start >= r.end || m.src.end <= r.start {
		return []irange{r}, []irange{}
	}

	// -RR-
	// MMMM,
	// RRRR
	// MMMM
	if m.src.start <= r.start && m.src.end >= r.end {
		startOffset := r.start - m.src.start
		endOffset := m.src.end - r.end
		dst := irange{start: m.dst.start + startOffset, end: m.dst.end - endOffset}
		return []irange{}, []irange{dst}
	}
	// --RRRR
	// MMMM
	if m.src.start <= r.start && m.src.end < r.end {
		startOffset := r.start - m.src.start
		left := irange{start: m.dst.start + startOffset, end: m.dst.end}
		right := irange{start: m.src.end, end: r.end}
		return []irange{right}, []irange{left}
	}
	// RRRR
	// --MMMM
	if m.src.start > r.start && m.src.end >= r.end {
		left := irange{start: r.start, end: m.src.start}
		endOffset := m.src.end - r.end
		right := irange{start: m.dst.start, end: m.dst.end - endOffset}
		return []irange{left}, []irange{right}
	}
	// RRRR
	// -MM-
	if m.src.start > r.start && m.src.end < r.end {
		left := irange{start: r.start, end: m.src.start}
		mid := irange{start: m.dst.start, end: m.dst.end}
		right := irange{start: m.src.end, end: r.end}
		return []irange{left, right}, []irange{mid}
	}
	panic("Either origin range or map ranges aren't real ranges!")
}

func (r irange) applyAll(maps []rangeMap) []irange {
	original := []irange{r}
	applied := make([]irange, 0, 10) // arbitrary capacity
	for _, m := range maps {
		updated := make([]irange, 0, 10)
		for _, r := range original {
			unmapped, mapped := r.apply(m)
			updated = append(updated, unmapped...)
			applied = append(applied, mapped...)
		}
		original = updated
	}

	applied = append(applied, original...)
	return applied
}

func applyToAll(rs []irange, maps []rangeMap) []irange {
	applied := make([]irange, 0, len(rs))
	for _, r := range rs {
		applied = append(applied, r.applyAll(maps)...)
	}
	return applied
}

func findSeedLocations(seeds []irange, srcMaps map[string]categoryMap) (locations []irange, err error) {
	const Location string = "location"
	key := "seed"
	seedValues := seeds
	maxSteps := len(srcMaps)
	for i := 0; i < maxSteps; i++ {
		cm, ok := srcMaps[key]
		if !ok {
			return nil, fmt.Errorf("Missing map from: %s", key)
		}
		seedValues = applyToAll(seedValues, cm.rangeMaps)
		key = cm.dst
		if key == "location" {
			return seedValues, nil
		}
	}
	return nil, fmt.Errorf("Seed-to-location maps do not form a chain: %v", srcMaps)
}

func computeLowestSeedLocation(inputPath string) (count int, err error) {
	puzzle, err := loadPuzzle(inputPath)
	if err != nil {
		return -1, err
	}

	if len(puzzle.seeds) == 0 {
		return -1, fmt.Errorf("Expected at least one seed in the puzzle, got none!")
	}

	// Run a simple algorithm on read data without any preprocessing.
	// Do not sort range maps and use binary search, not worth it.
	locations, err := findSeedLocations(puzzle.seeds, puzzle.srcMaps)
	if err != nil {
		return -1, err
	}
	minLocation := math.MaxInt
	for _, loc := range locations {
		if loc.start < minLocation {
			minLocation = loc.start
		}
	}
	return minLocation, nil
}

func main() {
	inputPathFlag := flag.String("input_path", "", "Path to puzzle input file")
	flag.Parse()

	if *inputPathFlag == "" {
		log.Fatal("Flag --input_path must be non-empty!")
	}

	points, err := computeLowestSeedLocation(*inputPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(points)

}
