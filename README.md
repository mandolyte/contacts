# Contact Processing for importing

## Background
This project is designed to replace the former methods used to 
create a Google Contacts (GC) import file from a TouchPoint (TP) contact database.

## 2022-12-11
In this iteration, TP has been used to generate the data from its database
with various filters applied to limit the data to that desired for importing
into GC.

It may be that the existing two Go source will not be needed. The first attempt will
be to simply take the TP export, map it into GC columns, and then import into GC.

Attempt #1 will be to see if GC will import existing TP export directly if the 
TP column headers are simply renamed to be the mapped GC columns.

Here the mapping for the columns:

| TP | GC |
| -- | -- |
| PeopleId |  |
| Title | Name Prefix |
| FirstName | Name |
| LastName | Given Name |
| Address | Address 1 - Street |
| Address2 | Address 1 - PO Box |
| City | Address 1 - City |
| State | Address 1 - Region |
| Country | Address 1 - Country |
| Zip | Address 1 - Postal Code |
| Email | E-mail 1 - Value |
| BirthDate | Birthday |
| BirthDay |  |
| JoinDate |  |
| HomePhone | Phone 1 - Value |
| CellPhone | Phone 2 - Value |
| WorkPhone | Phone 3 - Value |
| MemberStatus | Group Membership |
| Age |  |
| Married |  |
| Wedding |  |
| FamilyId |  |
| FamilyPosition | Location |
| Gender | Gender |
| School |  |
| Grade |  |
| FellowshipLeader | Directory Server |
| AttendPctBF |  |
| FellowshipClass | Billing Information |
| AltName |  |
| Employer |  |
| OtherId |  |
| Campus |  |
| DecisionDate |  |

**Results**
1. Name and Given Name are not combined. Thus I need to make one column containing first and last names.
2. Some of the columns came thru on their own and where shown as "information about" the person. These were: birthday, peopleId, age, married, familyId, attendPctBF. 
3. But others that I mapped explicitly did not make it. These were: FellowshipLeader, FellowshipClass, Gender, FamilyPosition
4. "Group Membership" was mapped as a label - interesting!
5. Phone, email, and address worked. But they are imported without the "home", "mobile", etc. type attribute. I need to add these values as new columns.


## 2012-12-12

Iteration 2 will simply redo some of the column headings in an attempt to improve the results.
1. map first name to given name and last name to family name.
2. map FellowshipLeader to 


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

**Results**
1. Title is coming thru and doesn't seem worthwhile... change it back to Title (mapped to Name Prefix)
2. FellowshipClass is coming thru with "Notes", but Leader isn't on "Subject". Try "Hobby" for latter.
3. Looks like all numbers and dates are coming thru. These seem to be treated differently.
4. Reverting to birthdate to see if both dates come thru.

**Results**
This looks nearly perfect and no programming needed - Woot! 

So I do need to do some programming to add labels to the phone numbers. So that and perhaps adding some things into the notes field is about it.

# Appendix A

Here are all the columns used in the GC export.
*Note: items with an exclamation mark are mapped.*

Name
!Given Name
Additional Name
!Family Name
Yomi Name
Given Name Yomi
Additional Name Yomi
Family Name Yomi
!Name Prefix
Name Suffix
Initials
Nickname
Short Name
Maiden Name
!Birthday
!Gender
Location
Billing Information
Directory Server
Mileage
!Occupation
!Hobby
Sensitivity
Priority
Subject
!Notes
Language
Photo
!Group Membership
E-mail 1 - Type
!E-mail 1 - Value
E-mail 2 - Type
E-mail 2 - Value
E-mail 3 - Type
E-mail 3 - Value
E-mail 4 - Type
E-mail 4 - Value
Phone 1 - Type
!Phone 1 - Value
Phone 2 - Type
!Phone 2 - Value
Phone 3 - Type
!Phone 3 - Value
Phone 4 - Type
Phone 4 - Value
Address 1 - Type
Address 1 - Formatted
!Address 1 - Street
!Address 1 - City
!Address 1 - PO Box
!Address 1 - Region
!Address 1 - Postal Code
!Address 1 - Country
Address 1 - Extended Address
Address 2 - Type
Address 2 - Formatted
Address 2 - Street
Address 2 - City
Address 2 - PO Box
Address 2 - Region
Address 2 - Postal Code
Address 2 - Country
Address 2 - Extended Address
Organization 1 - Type
Organization 1 - Name
Organization 1 - Yomi Name
Organization 1 - Title
Organization 1 - Department
Organization 1 - Symbol
Organization 1 - Location
Organization 1 - Job Description
Website 1 - Type
Website 1 - Value