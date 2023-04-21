#!/bin/sh
# Determine OS platform
# shellcheck source=/dev/null

create_link() {
  # shellcheck disable=SC2015
	[ -f "/usr/bin/q" ] && rm -rf /usr/bin/q || true
	# shellcheck disable=SC2015
	[ -f "/usr/local/bin/q" ] && rm -rf /usr/local/bin/q || true
	# shellcheck disable=SC2015
	[ -f "/usr/local/bin/qcadmin" ] && rm -rf /usr/local/bin/qcadmin || true
  # shellcheck disable=SC2015
  [ -f "/usr/bin/qcadmin" ] && ln -s /usr/bin/qcadmin /usr/bin/q || true
}

summary() {
	echo "----------------------------------------------------------------------"
	echo "quickon package has been successfully installed."
	echo ""
	echo " Please follow the next steps to start the software:"
	echo ""
	echo "    q init --help"
	echo ""
	echo ""
	echo "----------------------------------------------------------------------"
}

{
  create_link
  summary
}
