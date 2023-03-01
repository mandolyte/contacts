# Combined Community and Shepherding Resport

## Process Steps

1. Download Shepherding report
2. Download Community Group report
 - Go to TP, click Involvements
 - Click Search
 - Enter "CG" for name
 - Should return all community groups
 - Click to download and select to Export Members (this will export from all the groups)
 - Convert to CSV; for example, by importing into Google sheets, then exporting
 - Upload to this folder
3. Move them to this folder with appropriate names (see sample commands below)
4. Then execute the `run.sh` script to combine them
5. Upload to Shepherding Committee folders or email


Sample Command:
```sh
$ go run cg-report.go -cg Community_2022-12-20.csv -sg Shepherding_2022-12-20.csv -o x.csv
2022/12/20 13:21:43 Starting.
2022/12/20 13:21:43 All Done!
```

Community Columns:
0. InvolvmentId
1. Involvement
2. FirstName
3. LastName
4. Gender
5. Grade
6. ShirtSize
7. Request
8. Amount
9. AmountPaid
10. HasBalance
11. Groups
12. Email
13. HomePhone
14. CellPhone
15. WorkPhone
16. Age
17. BirthDate
18. JoinDate
19. MemberStatus
20. SchoolOther
21. LastAttend
22. AttendPct
23. AttendStr
24. MemberType
25. MemberInfo
26. InactiveDate
27. Medical
28. PeopleId
29. EnrollDate
30. Tickets
