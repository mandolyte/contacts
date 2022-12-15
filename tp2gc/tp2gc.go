package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

func usage(msg string) {
	fmt.Println(msg)
	fmt.Print("Usage: tp2gc -i input -o output")
	flag.PrintDefaults()
	os.Exit(0)
}

func debug(msg string) {
	if *dbg {
		log.Println(msg)
	}
}

var input = flag.String("i", "", "Input CSV filename")
var output = flag.String("o", "", "Output CSV filename")
var dbg = flag.Bool("debug", false, "Show debug console output")

func main() {

	log.Print("Starting.\n")
	flag.Parse()

	if *input == "" {
		usage("input Touchpoint CSV filename is missing")
	}
	if *output == "" {
		usage("outout Google Contacts CSV filename is missing")
	}

	/*
		Open the Touchpoint input file
	*/
	// open the CSV file
	f, ferr := os.Open(*input)
	if ferr != nil {
		log.Fatal("os.Open() Error:" + ferr.Error())
	}
	r := csv.NewReader(f)
	// ignore expectations of fields per row
	r.FieldsPerRecord = -1

	// read all the input rows into a table
	inrows, err := r.ReadAll()
	if err != nil {
		log.Fatal("r.ReadAll() Error:" + err.Error())
	}

	// output rows
	var outrows [][]string
	// add the header
	outrows = append(outrows, makeHeader())

	for i := range inrows {
		if i == 0 {
			continue
		}
		var row []string

		// title -> name prefix
		row = append(row, inrows[i][1])
		// first name -> given name
		row = append(row, inrows[i][2])
		// last name -> familuy name
		row = append(row, inrows[i][3])
		// address -> address 1 - street
		row = append(row, inrows[i][4])
		// address2 -> address 1 - po box
		row = append(row, inrows[i][5])
		// city -> address 1 - city
		row = append(row, inrows[i][6])
		// state -> address 1 - region
		row = append(row, inrows[i][7])
		// country -> address 1 - country
		row = append(row, inrows[i][8])
		// zip -> address 1 - postal code
		row = append(row, inrows[i][9])
		// email -> email 1 - value
		row = append(row, inrows[i][10])
		// birthDate -> birthday
		row = append(row, inrows[i][11])
		// home phone type
		row = append(row, "Home")
		// home phone value
		row = append(row, inrows[i][14])
		// mobile phone type
		row = append(row, "Mobile")
		// mobile phone value
		row = append(row, inrows[i][15])
		// work phone type
		row = append(row, "Work")
		// work phone value
		row = append(row, inrows[i][16])
		// member status -> group membership
		row = append(row, inrows[i][17])
		// notes
		// lets combine some things to stuff into the notes field
		note := "Family Position is: " + inrows[i][22] + "\n"
		note += "Fellowship Leader is: " + inrows[i][26] + "\n"
		note += "Fellowship Class is: " + inrows[i][28] + "\n"
		note += "Employer is: " + inrows[i][30]
		row = append(row, note)
		outrows = append(outrows, row)
	}

	f.Close()

	// write out the chronologically concatenated CSV
	w, werr := os.Create(*output)
	if werr != nil {
		log.Fatal("os.Open() Error:" + werr.Error())
	}

	csvw := csv.NewWriter(w)
	// write the header row first
	csvw.WriteAll(outrows)
	w.Close()

	log.Print("All Done!\n")
}

func makeHeader() []string {
	h := make([]string, 19)
	h[0] = "Name Prefix"
	h[1] = "Given Name"
	h[2] = "Family Name"
	h[3] = "Address 1 - Street"
	h[4] = "Address 1 - PO Box"
	h[5] = "Address 1 - City"
	h[6] = "Address 1 - Region"
	h[7] = "Address 1 - Country"
	h[8] = "Address 1 - Postal Code"
	h[9] = "E-mail 1 - Value"
	h[10] = "BirthDate"
	h[11] = "Phone 1 - Type"
	h[12] = "Phone 1 - Value"
	h[13] = "Phone 2 - Type"
	h[14] = "Phone 2 - Value"
	h[15] = "Phone 3 - Type"
	h[16] = "Phone 3 - Value"
	h[17] = "Group Membership"
	h[18] = "Notes"

	return h
}

/* Documentation on mapping

| Index | TP | GC |
| -- | -- | -- |
| 0 | PeopleId |  |
| 1 | Title | Name Prefix |
| 2 | FirstName | Given Name |
| 3 | LastName | Family Name |
| 4 | Address | Address 1 - Street |
| 5 | Address2 | Address 1 - PO Box |
| 6 | City | Address 1 - City |
| 7 | State | Address 1 - Region |
| 8 | Country | Address 1 - Country |
| 9 | Zip | Address 1 - Postal Code |
| 10 | Email | E-mail 1 - Value |
| 11 | BirthDate | |
| 12 | BirthDay |  |
| 13 | JoinDate |  |
| 14 | HomePhone | Phone 1 - Value |
| 15 | CellPhone | Phone 2 - Value |
| 16 | WorkPhone | Phone 3 - Value |
| 17 | MemberStatus | Group Membership |
| 18 | Age |  |
| 19 | Married |  |
| 20 | Wedding |  |
| 21 | FamilyId |  |
| 22 | FamilyPosition | Occupation |
| 23| Gender | Gender |
| 24 | School |  |
| 25 | Grade |  |
| 26 | FellowshipLeader | Hobby |
| 27 | AttendPctBF |  |
| 28 | FellowshipClass | Notes |
| 29 | AltName |  |
| 30 | Employer |  |
| 31 | OtherId |  |
| 32 | Campus |  |
| 33 | DecisionDate |  |
*/
