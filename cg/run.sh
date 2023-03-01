#!/bin/sh
if [ "$1"x = "x" ]
then
	echo enter ISO date, for example 2023-02-28 of the files
	exit
fi

go run cg-report.go \
	-cg Community_$1.csv \
	-sg Shepherding_$1.csv \
	-o Combined_$1.csv
