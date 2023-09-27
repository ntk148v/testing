#!/bin/bash
rm -rf output.txt
go test -bench=. -benchtime=10s -benchmem > output.txt
