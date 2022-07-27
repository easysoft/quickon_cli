#!/bin/sh
set -ex

rm -rf completions
mkdir completions
go build -o /tmp/qcadmin
/tmp/qcadmin version | grep -q linux && touch /tmp/.qcadmin-linux
for sh in bash zsh fish; do
	/tmp/qcadmin completion "$sh" >"completions/qcadmin.$sh"
  cp -a "completions/qcadmin.$sh" "completions/q.$sh"
  [ -f "/tmp/.qcadmin-linux" ] && (
    sed -i "s#qcadmin#q#g" "completions/q.$sh"
  ) || (
    sed -i "" "s#qcadmin#q#g" "completions/q.$sh"
  )
done
rm -rf /tmp/qcadmin /tmp/.qcadmin-linux
