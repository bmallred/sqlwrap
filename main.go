package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	a := len(os.Args)
	if a == 1 {
		log.Fatal("Must specific a file")
	}

	for i := 1; i < a; i++ {
		path := os.Args[i]

		// Open the file initially for read only
		fi, err := os.OpenFile(path, os.O_RDONLY, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer fi.Close()

		// Loop through all the lines of the file look for INSERT statements to wrap
		lines := []string{}
		scanner := bufio.NewScanner(fi)
		for scanner.Scan() {
			line := scanner.Text()
			idx := strings.Index(line, "INSERT")
			if idx == 0 || len(lines) == 0 {
				lines = append(lines, line)
			} else {
				previous := len(lines) - 1
				lines[previous] = fmt.Sprintf("%s' + CHAR(13) + N'%s", lines[previous], line)
			}
		}

		// Check for errors in the scanner
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Close for sanity's sake
		err = fi.Close()
		if err != nil {
			log.Fatal(err)
		}

		// Truncate the file and open for write
		fi, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer fi.Close()

		// Write each of the new lines to the file
		writer := bufio.NewWriter(fi)
		for _, line := range lines {
			fmt.Fprintln(writer, line)
		}
		writer.Flush()
	}
}
