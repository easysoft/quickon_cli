#!/bin/sh
# Determine OS platform
# shellcheck source=/dev/null
. /etc/os-release

ensure_sudo() {
	if [ "$(id -u)" = "0" ]; then
		echo "Sudo permissions detected"
	else
		echo "No sudo permission detected, please run as sudo"
		exit 1
	fi
}

{
  ensure_sudo
}
