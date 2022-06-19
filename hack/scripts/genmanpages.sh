#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run ./main.go man | gzip -c -9 >manpages/q.1.gz
go run ./main.go man | gzip -c -9 >manpages/qcadmin.1.gz
