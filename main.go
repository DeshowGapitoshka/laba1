package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	Time int
	Line int
}

func mediana(t []int) time.Duration {
	n := len(t)
	if n == 0 {
		return 0
	}
	sort.Ints(t)
	mid := n / 2
	if n%2 == 1 {
		return time.Duration(t[mid])
	}
	return time.Duration((t[mid-1] + t[mid]) / 2)
}

func laba() int {
	stack := make(map[string]([]Session))
	var active int
	var peak int
	var peakLine int
	var lineNum int
	var closed int
	var total int64
	args := os.Args
	if len(args) != 3 {
		os.Exit(1)
	}

	inPath := strings.Trim(args[1], `'"`)
	outPath := strings.Trim(args[2], `'"`)
	file, err := os.Open(inPath)
	if err != nil {
		os.Exit(1)
	}
	fileWrite, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer fileWrite.Close()
	defer file.Close()
	scanner := bufio.NewScanner(file)
	start := time.Now()
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if len(line) < 3 {
			continue
		}
		parts := strings.Fields(line)
		user := parts[1]
		time, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Fprintf(fileWrite, "ERROR %d bad_event\n", lineNum)
			continue
		}
		switch parts[0] {
		case "LOGIN":
			stack[user] = append(stack[user], Session{Time: time, Line: lineNum})
			active++
			if active > peak {
				peak = active
				peakLine = lineNum
			}

		case "LOGOUT":
			st := stack[user]
			if len(st) == 0 {
				fmt.Fprintf(fileWrite, "ERROR %d logout_without_login\n", lineNum)
				continue
			}
			sess := st[len(st)-1]
			stack[user] = st[:len(st)-1]
			total += int64(time - sess.Time)
			active--
			closed++

		default:
			fmt.Fprintf(fileWrite, "ERROR %d bad_event\n", lineNum)
		}
	}
	var opens []int
	for _, st := range stack {
		for _, sess := range st {
			opens = append(opens, sess.Line)
		}
	}
	sort.Ints(opens)
	t := time.Since(start)
	for _, line := range opens {
		fmt.Fprintf(fileWrite, "ERROR %d session_not_closed\n", line)
	}
	fmt.Fprintf(fileWrite, "SESSIONS %d\n", closed)
	fmt.Fprintf(fileWrite, "peak_concurrent=%d at line=%d\n", peak, peakLine)
	fmt.Fprintf(fileWrite, "avg_duration=%f\n", float64(total)/float64(closed))
	return int(t)
}

func main() {
	var times []int
	for i := 0; i < 3; i++ {
		laba()
	}
	for i := 0; i < 7; i++ {
		times = append(times, laba())
	}
	fmt.Printf("Mediana time:%v", mediana(times))
}
