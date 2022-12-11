package main

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// CSVDIR must run in the right directory for this to work!
var CSVDIR = "CSV"

func main() {
	/*
		First, open the directory and get all the CSV files
	*/
	finfo, rderr := ioutil.ReadDir("CSV")
	if rderr != nil {
		log.Fatal("Error:", rderr)
	}
	dirlist := make([]string, len(finfo))

	/*
		Read CSV filenames into array and
		sort descending. The latest CSV
		will be the first and will supersede any dups
	*/
	var i int
	for _, fi := range finfo {
		dirlist[i] = fi.Name()
		i++
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dirlist)))

	/*
		Process each file in the array
	*/
	// first a map to collect all the rows
	m := make(map[string][]string)
	var file int
	for _, fname := range dirlist {
		log.Printf("Working on CSV:%s", fname)

		// open the CSV file
		f, ferr := os.Open(CSVDIR + "/" + fname)
		if ferr != nil {
			log.Fatal("os.Open() Error:" + ferr.Error())
		}
		file++
		r := csv.NewReader(f)
		// ignore expectations of fields per row
		r.FieldsPerRecord = -1

		// read loop for CSV
		var line int
		for {
			// read the csv file
			cells, rerr := r.Read()

			// Error handling
			if rerr == io.EOF {
				break
			}
			if rerr != nil {
				log.Fatal("r.Read() Error:" + rerr.Error())
			}
			line++
			if file == 1 && line == 1 {
				m[cells[0]] = cells // capture header on first file only
				continue
			}
			if line == 1 {
				continue // skip header on each file
			}
			//log.Printf("... on line:%v", line)
			// do some normailization work:
			// a. make column D normal case
			// b. replace column A with column B + " " + column D
			cells[1] = strings.TrimSpace(cells[1])
			cells[3] = strings.TrimSpace(cells[3])
			cells[3] = string(cells[3][0]) + strings.ToLower(cells[3][1:])
			cells[0] = cells[1] + " " + cells[3]

			// for each row check if key (column 1) is in map
			// discard if in map (ie, it's been superseded)
			// if not in map, the store entire line as value
			// for the column 1 key
			_, exists := m[cells[0]]
			if exists {
				continue
			}

			// the name key may contain an indirection
			// detect with a regexp and skip it
			sm, smerr := regexp.MatchString("- see", cells[0])
			if smerr != nil {
				log.Fatal("regexp error" + smerr.Error())
			}
			if sm {
				continue
			}
			m[cells[0]] = cells
		}
		f.Close()
	}

	/*
		After processing all the rows in all the files,
		then create a slice matching the len of the map.

		Then loop thru the map storing all lines in the slice.

		Then sort ascending the slice and finally, write out
		all lines to the concatenated result CSV
	*/

	// remove the header row first!
	var hdrVal = "Name"
	hdr := m[hdrVal]
	delete(m, hdrVal)

	allrows := make([][]string, len(m))
	i = 0
	for _, v := range m {
		allrows[i] = v
		i++
	}

	// write out the chronologically concatenated CSV
	w, werr := os.Create("cat_chrono.csv")
	if werr != nil {
		log.Fatal("os.Open() Error:" + werr.Error())
	}

	csvw := csv.NewWriter(w)
	// write the header row first
	csvw.Write(hdr)
	csvw.WriteAll(allrows)
	w.Close()

	log.Print("All Done!\n")
}
