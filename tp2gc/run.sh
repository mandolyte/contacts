#!/bin/sh
if [ "$1"x = "x" ]
then
	echo enter ISO date, for example 2023-02-28 of the files
	exit
fi

go run tp2gc.go \
	-i $1_People.csv \
	-o $1_Google_Contacts.csv
