#!/bin/sh
# Determine OS platform
# shellcheck source=/dev/null

# shellcheck disable=SC2015
[ -f "/usr/bin/q" ] && rm -rf /usr/bin/q || true
[ -f "/usr/local/bin/q" ] && rm -rf /usr/local/bin/q || true
