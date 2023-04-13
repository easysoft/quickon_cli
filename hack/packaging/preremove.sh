#!/bin/sh
# Determine OS platform
# shellcheck source=/dev/null
. /etc/os-release

[ -f "/usr/bin/q" ] && rm -rf /usr/bin/q || true
