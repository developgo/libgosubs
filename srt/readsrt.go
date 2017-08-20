package srt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//LoadSRT loads the provided file into the given object.
//It fixes the \ufeff problem that some parsers have.
func LoadSrt(v *SubRip, filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Cannot read file", filepath)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)
	z := &Subtitle{}
	for scanner.Scan() {
		line := scanner.Text()
		//Ufeff problem fix
		if strings.HasPrefix(line, "\ufeff") {
			line = strings.Replace(line, "\ufeff", "", -1)
		}
		//The most jackass SRT parser ever. Appends to an object.
		i, err := strconv.Atoi(line)
		if err == nil && z.Id == 0 {
			z.Id = int(i)
			//A bit more jackassery, because if z.Start and z.End are set, then welp.
		} else if strings.Contains(line, "-->") && z.Start == "" && z.End == "" {
			split := strings.Split(line, "-->")
			z.Start = split[0]
			z.End = split[1]
		} else if line != "" && z.Start != "" && z.End != "" && z.Id != 0 {
			z.Line = append(z.Line, line)
		} else if line == "" {
			//Clear object on newline
			v.subtitle.content = append(v.subtitle.content, *z)
			z = &Subtitle{}
		} else {
			//At some point, we need to start actually returning errors.
			//Wouldn't that be nice?
			fmt.Println("Error parsing .srt. Stray newline?")
		}

	}
	//Since the last subtitle often won't have a newline, append everything and clear object one last time
	v.subtitle.content = append(v.subtitle.content, *z)
	z = &Subtitle{}
	defer f.Close()
}

//Generic loader for srt files. Takes the path of the file being opened as the argument.
func ParseSrt(filename string) *SubRip {
	v := &SubRip{}
	LoadSrt(v, filename)
	return v
}
