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
	fmt.Print("Usage: cg-report -sg sg_input -cg cg_input -o output")
	fmt.Println(("This report takes two inputs:"))
	fmt.Println(("- the shepherding group report as a CSV file"))
	fmt.Println(("- the community group report as a CSV file"))
	fmt.Println(("The output will be a join of the two showing both"))
	fmt.Println(("elder and community group associations or absence."))
	flag.PrintDefaults()
	os.Exit(0)
}

func debug(msg string) {
	if *dbg {
		log.Println(msg)
	}
}

var sg_input = flag.String("sg", "", "Shepherding CSV filename")
var cg_input = flag.String("cg", "", "Community CSV filename")
var output = flag.String("o", "", "Output CSV filename")
var dbg = flag.Bool("debug", false, "Show debug console output")

func main() {

	log.Print("Starting.\n")
	flag.Parse()

	if *sg_input == "" {
		usage("input Touchpoint CSV filename is missing")
	}
	if *cg_input == "" {
		usage("input Touchpoint CSV filename is missing")
	}
	if *output == "" {
		usage("outout Google Contacts CSV filename is missing")
	}

	/*
		Open the Touchpoint shepherding input file
	*/
	f, ferr := os.Open(*sg_input)
	if ferr != nil {
		log.Fatal("os.Open() Error:" + ferr.Error())
	}
	r := csv.NewReader(f)
	// ignore expectations of fields per row
	r.FieldsPerRecord = -1

	// read all the input rows into a table
	sgrows, err := r.ReadAll()
	if err != nil {
		log.Fatal("r.ReadAll() Error:" + err.Error())
	}
	f.Close()

	/*
		Open the Touchpoint community input file
	*/
	f, ferr = os.Open(*cg_input)
	if ferr != nil {
		log.Fatal("os.Open() Error:" + ferr.Error())
	}
	r = csv.NewReader(f)
	// ignore expectations of fields per row
	r.FieldsPerRecord = -1

	// read all the input rows into a table
	cgrows, err := r.ReadAll()
	if err != nil {
		log.Fatal("r.ReadAll() Error:" + err.Error())
	}
	f.Close()

	// output rows
	var outrows [][]string
	// add the header
	outrows = append(outrows, makeHeader())

	for i := range sgrows {
		if i == 0 {
			continue
		}

		var cg_leader string
		for j := range cgrows {
			if j == 0 {
				continue
			}
			if cgrows[j][28] == sgrows[i][0] {
				cg_leader = cgrows[j][1]
				break
			}
		}

		var row []string

		// first and last names -> name
		_name := sgrows[i][2] + " " + sgrows[i][3]
		row = append(row, _name)

		// zip -> address 1 - postal code
		row = append(row, sgrows[i][9])

		// email -> email 1 - value
		row = append(row, sgrows[i][10])

		// phone
		_phone := sgrows[i][15] // prefer mobile
		if _phone == "" {
			_phone = sgrows[i][14] // then home
			if sgrows[i][14] == "" {
				_phone = sgrows[i][16] // then work
			}
		}
		row = append(row, _phone)

		// member status -> group membership
		row = append(row, sgrows[i][17])

		// cg group
		row = append(row, cg_leader)

		// sg leader
		row = append(row, sgrows[i][26])

		outrows = append(outrows, row)
	}
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
	h[0] = "Name"
	h[1] = "Zip"
	h[2] = "Email"
	h[3] = "Phone"
	h[4] = "Member Status"
	h[5] = "Community Group"
	h[6] = "Shepherding Group"

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
