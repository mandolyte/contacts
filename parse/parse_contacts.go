package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type expecting int

const (
	family expecting = iota
	childrenOrAddress
	addressOrPhone
	phoneOrEmail
	email
)

var currentPlace expecting

type info int

const (
	headFirstName info = iota
	headLastName
	wifeFirstName
	wifeLastName
	childrenNote
	address
	homePhone
	headPhone
	wifePhone
	headEmail
	wifeEmail
)

var familyData map[info]string
var familycount int
var w *csv.Writer

var input = flag.String("i", "", "Input CSV filename")
var output = flag.String("o", "", "Output CSV filename")
var dbg = flag.Bool("debug", false, "Show debug console output")

func main() {
	log.Print("Starting.\n")
	flag.Parse()

	if *input == "" {
		usage("input CSV filename is missing")
	}
	if *output == "" {
		usage("input CSV filename is missing")
	}

	familyData = make(map[info]string)

	// open input file
	fi, ferr := os.Open(*input)
	if ferr != nil {
		log.Fatalf("Error os.Open():%v\n", ferr)
	}
	defer fi.Close()

	// open output file
	fo, foerr := os.Create(*output)
	if foerr != nil {
		log.Fatal("os.Create() Error:" + foerr.Error())
	}
	defer fo.Close()
	w = csv.NewWriter(fo)

	// write out the headers ...
	werr := w.Write(headers)
	if werr != nil {
		log.Fatalf("Error w.Write(headers):%v\n", werr)
	}

	scanner := bufio.NewScanner(fi)
	scanner.Split(bufio.ScanLines)

	linecount := 0
	currentPlace = family
	for scanner.Scan() {
		curline := scanner.Text()
		linecount++
		//debug(fmt.Sprintf("[%v]:%v", linecount, curline))
		if len(curline) < 2 {
			//log.Printf("Short line, skipping line %v", linecount)
			continue
		}
		switch currentPlace {
		case family:
			currentPlace = processFamily(curline)
		case childrenOrAddress:
			currentPlace = processChildrenOrAddress(curline)
		case addressOrPhone:
			currentPlace = processAddressOrPhone(curline)
		case phoneOrEmail:
			currentPlace = processPhoneOrEmail(curline)
		case email:
			currentPlace = processEmail(curline)
		default:
			log.Fatalf("Unknown switch value:%v", currentPlace)
		}
	}
	if serr := scanner.Err(); serr != nil {
		log.Fatalf("Errror scanner.Err():%v", serr)
	}
	w.Flush()
	log.Printf("Total number of familes: %v\n", familycount)
	log.Print("All Done!\n")
}

func processFamily(line string) expecting {
	//debug("Enter processFamily")
	if familycount > 0 {
		processFamilyData(w)
	}
	familycount++

	/* Two basic examples:
	Admiraal, Jeff* & Nancy*
	Allen, Angela*

	The case where wife retains maiden name:
	DeBord, Tom & Nicole Lewis

	1. Must eliminate the asterisks
	2. look for the ampersand and split into both husband and wife
	*/

	// split on comma
	names := strings.Split(line, ",")
	for n := range names {
		names[n] = strings.TrimSpace(names[n])
		names[n] = strings.Replace(names[n], "*", "", -1)
	}
	// there should always be 2 elements, the second containing
	// either one or two first names
	if len(names) != 2 {
		log.Fatalf("No first name on line <%v>", line)
	}
	firsts := strings.Split(names[1], "&")
	for n := range firsts {
		firsts[n] = strings.TrimSpace(firsts[n])
	}
	if len(firsts) == 1 {
		// only head of household case
		familyData[headFirstName] = firsts[0]
		familyData[headLastName] = names[0]
	} else if len(firsts) == 2 {
		familyData[headFirstName] = firsts[0]
		familyData[headLastName] = names[0]
		wifenames := strings.Split(firsts[1], " ")
		if len(wifenames) > 1 {
			familyData[wifeFirstName] = wifenames[0]
			familyData[wifeLastName] = wifenames[1]
		} else {
			familyData[wifeFirstName] = firsts[1]
			familyData[wifeLastName] = names[0]
		}
	}
	//debug("Exit processFamily")

	return childrenOrAddress
}

func processChildrenOrAddress(line string) expecting {
	//debug("Enter processChildrenOrAddress")
	/* If there are children, then the next lines look like this:
	Tristan* (18)
	Aaron* (16)
	Cara (12)
	Christina (9)

	If there's no pattern like this, then it must be the address
	*/
	var childline = regexp.MustCompile(`\(\d+\)|^[A-Z][a-z].*$`)
	matched := childline.MatchString(line)
	if matched {
		familyData[childrenNote] += line + ";"
		//debug("Children Note is:\n" + familyData[childrenNote])
		return childrenOrAddress
	}
	//debug("Exit processChildrenOrAddress")

	return processAddressOrPhone(line)
}

func processAddressOrPhone(line string) expecting {
	//debug("Enter processAddressOrPhone")
	/* An address looks like this:
	4712 Flagstone Dr.
	Mason                         , OH 45040

	Phone numbers look like this:
	H 513-683-6886
	C Angela: 513-235-6261
	*/

	// the pattern for phone is easier, test for it
	if strings.HasPrefix(line, "H ") {
		return processPhoneOrEmail(line)
	}

	if strings.HasPrefix(line, "C ") {
		return processPhoneOrEmail(line)
	}

	// if no phone number, look for an email before
	// assuming it is an address line
	if strings.Contains(line, "@") {
		return processEmail(line)
	}

	// ok, must be an address
	re := regexp.MustCompile(`( {2,})`)
	line = re.ReplaceAllLiteralString(line, " ")
	familyData[address] += line + "\n"
	////debug("Address is:\n" + familyData[address])

	//debug("Exit processAddressOrPhone")

	// might be another line of address, so send it back here...
	return addressOrPhone
}

func processPhoneOrEmail(line string) expecting {
	//debug("Enter processPhoneOrEmail")
	/* Phone numbers look like this:
	H 513-683-6886
	C Angela: 513-235-6261
	*/

	if strings.HasPrefix(line, "H ") {
		familyData[homePhone] = line[2:]
		return phoneOrEmail
	}

	if strings.HasPrefix(line, "C ") {
		phtemp := strings.Split(line[2:], ":")
		if len(phtemp) != 2 {
			log.Fatalf("No colon in phone line!, line is: %v\n", line)
		}
		if phtemp[0] == familyData[headFirstName] {
			familyData[headPhone] = phtemp[1][1:]
		} else if phtemp[0] == familyData[wifeFirstName] {
			familyData[wifePhone] = phtemp[1][1:]
		}
		return phoneOrEmail
	}

	//debug("Exit processPhoneOrEmail")
	// must be an email
	return processEmail(line)
}

func processEmail(line string) expecting {
	//debug("Enter processEmail:" + line)
	/* Emails like this:
	Jeff: taccsan@yahoo.com
	*/

	// check for the colon
	emtemp := strings.Split(line, ":")
	if len(emtemp) < 2 {
		// go on to process a new family
		return processFamily(line)
	}
	//debug("first name for email is:" + emtemp[0])
	if emtemp[0] == familyData[headFirstName] {
		familyData[headEmail] = emtemp[1][1:]
		return email
	}

	if emtemp[0] == familyData[wifeFirstName] {
		familyData[wifeEmail] = emtemp[1][1:]
		return email
	}

	return family
}

func usage(msg string) {
	fmt.Println(msg)
	fmt.Print("Usage: parse_contacts -i input -o output")
	flag.PrintDefaults()
	os.Exit(0)
}

func debug(msg string) {
	if *dbg {
		log.Println(msg)
	}
}

/* some data stuff */
// the headers used by Google Contact Imports
var headers = []string{
	"Name", "Given Name", "Additional Name", "Family Name",
	"Yomi Name", "Given Name Yomi", "Additional Name Yomi",
	"Family Name Yomi", "Name Prefix", "Name Suffix", "Initials",
	"Nickname", "Short Name", "Maiden Name", "Birthday", "Gender",
	"Location", "Billing Information", "Directory Server",
	"Mileage", "Occupation", "Hobby", "Sensitivity", "Priority", "Subject",
	"Notes", "Group Membership",
	"E-mail 1 - Type", "E-mail 1 - Value",
	"E-mail 2 - Type", "E-mail 2 - Value",
	"Phone 1 - Type", "Phone 1 - Value",
	"Phone 2 - Type", "Phone 2 - Value",
	"Phone 3 - Type", "Phone 3 - Value",
	"Address 1 - Type", "Address 1 - Formatted",
	"Address 1 - Street", "Address 1 - City",
	"Address 1 - PO Box", "Address 1 - Region",
	"Address 1 - Postal Code", "Address 1 - Country",
	"Address 1 - Extended Address",
	"Organization 1 - Type", "Organization 1 - Name",
	"Organization 1 - Yomi Name", "Organization 1 - Title",
	"Organization 1 - Department", "Organization 1 - Symbol",
	"Organization 1 - Location",
	"Organization 1 - Job Description",
	"Custom Field 1 - Type", "Custom Field 1 - Value",
}

func processFamilyData(w *csv.Writer) {
	//debug("Enter processFamilyData")
	//debug("Family Data:\n" + fmt.Sprintf("%v\n", familyData))
	processHead(w)
	if familyData[wifeFirstName] != "" {
		processWife(w)
	}

	// reset the familyData
	for k := range familyData {
		familyData[k] = ""
	}
	//debug("Exit processFamilyData")

}

func processHead(w *csv.Writer) {
	var cols [56]string
	cols[0] = familyData[headFirstName] + " " + familyData[headLastName]
	cols[1] = familyData[headFirstName]
	cols[3] = familyData[headLastName]
	cols[25] = strings.TrimSpace(familyData[childrenNote])
	cols[26] = "NCCC-DIR-IMPORT"
	cols[27] = familyData[headFirstName]
	cols[28] = familyData[headEmail]
	cols[31] = "Home"
	cols[32] = familyData[homePhone]
	cols[33] = "Mobile"
	cols[34] = familyData[headPhone]
	cols[37] = "Home"
	cols[38] = strings.TrimSpace(familyData[address])
	row := make([]string, 56)
	for n := range cols {
		//row = append(row, cols[n])
		row[n] = cols[n]
	}
	werr := w.Write(row)
	if werr != nil {
		log.Fatalf("Error w.Write(headers):%v\n", werr)
	}

}

func processWife(w *csv.Writer) {
	var cols [56]string
	cols[0] = familyData[wifeFirstName] + " " + familyData[wifeLastName]
	cols[1] = familyData[wifeFirstName]
	cols[3] = familyData[wifeLastName]
	cols[25] = strings.TrimSpace(familyData[childrenNote])
	cols[26] = "NCCC-DIR-IMPORT"
	cols[27] = familyData[wifeFirstName]
	cols[28] = familyData[wifeEmail]
	cols[31] = "Home"
	cols[32] = familyData[homePhone]
	cols[33] = "Mobile"
	cols[34] = familyData[wifePhone]
	cols[37] = "Home"
	cols[38] = strings.TrimSpace(familyData[address])
	row := make([]string, 56)
	for n := range cols {
		//row = append(row, cols[n])
		row[n] = cols[n]
	}
	werr := w.Write(row)
	if werr != nil {
		log.Fatalf("Error w.Write(headers):%v\n", werr)
	}

}
