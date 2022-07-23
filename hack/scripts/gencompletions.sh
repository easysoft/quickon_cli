#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run ./main.go completion "$sh" >"completions/qcadmin.$sh"
  cp -a "completions/qcadmin.$sh" "completions/q.$sh"
  sed -i "" "s#qcadmin#q#g" "completions/q.$sh"
done
