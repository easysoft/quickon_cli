#!/bin/sh
set -e
set -o noglob

# --- helper functions for logs ---
info()
{
    echo '[INFO] ' "$@"
}
warn()
{
    echo '[WARN] ' "$@" >&2
}
fatal()
{
    echo '[ERROR] ' "$@" >&2
    exit 1
}

COS_URL=https://pkg.qucheng.com/qucheng/cli

# --- define needed environment variables ---
setup_env() {
    # --- use sudo if we are not already root ---
    SUDO=sudo
    if [ $(id -u) -eq 0 ]; then
        SUDO=
    fi
		BIN_DIR=/usr/local/bin
    if ! $SUDO sh -c "touch ${BIN_DIR}/q-ro-test && rm -rf ${BIN_DIR}/q-ro-test"; then
      if [ -d /opt/bin ]; then
        BIN_DIR=/opt/bin
      fi
    fi
}

# --- verify an executable qcadmin binary is installed ---
verify_qcadmin_is_executable() {
    if [ ! -x ${BIN_DIR}/q ]; then
        fatal "Executable qcadmin binary not found at ${BIN_DIR}/qcadmin"
    fi
}

# --- verify existence of network downloader executable ---
verify_downloader() {
    # Return failure if it doesn't exist or is no executable
    [ -x "$(command -v $1)" ] || return 1

    # Set verified executable as our downloader program and return success
    DOWNLOADER=$1
    return 0
}

# --- create temporary directory and cleanup when done ---
setup_tmp() {
    TMP_DIR=$(mktemp -d -t qcadmin-install.XXXXXXXXXX)
    TMP_HASH=${TMP_DIR}/qcadmin.hash
    TMP_BIN=${TMP_DIR}/qcadmin.bin
    cleanup() {
        code=$?
        set +e
        trap - EXIT
        rm -rf ${TMP_DIR}
        exit $code
    }
    trap cleanup INT EXIT
}

setup_quickon() {
  [ -d "/opt/quickon/backup" ] || (
    mkdir -p /opt/quickon/backup
    chmod 777 /opt/quickon/backup
  )
}

# --- use desired qcadmin version if defined or find version from channel ---
get_release_version() {
		VERSION="edge"
    info "Using ${VERSION} as release"
}

# --- set arch and suffix, fatal if architecture not supported ---
setup_verify_arch() {
    if [ -z "$ARCH" ]; then
        ARCH=$(uname -m)
    fi
    case $ARCH in
        amd64|x86_64)
            ARCH=amd64
            SUFFIX=${ARCH}
            ;;
        arm64|aarch64)
            ARCH=arm64
            SUFFIX=${ARCH}
            ;;
        arm*)
            ARCH=arm
            SUFFIX=${ARCH}hf
            ;;
        *)
            fatal "Unsupported architecture $ARCH"
    esac
}


# --- download from url ---
download() {
    [ $# -eq 2 ] || fatal 'download needs exactly 2 arguments'

    case $DOWNLOADER in
        curl)
            curl -o $1 -sfL $2
            ;;
        wget)
            wget -qO $1 $2
            ;;
        *)
            fatal "Incorrect executable '$DOWNLOADER'"
            ;;
    esac

    # Abort if download command failed
    [ $? -eq 0 ] || fatal 'Download failed'
}

# --- download binary from cos url ---
download_binary() {
    BIN_URL=${COS_URL}/${VERSION}/qcadmin_linux_${SUFFIX} # qcadmin_linux_amd64
    info "Downloading binary ${COS_URL}/${VERSION}/q"
    download ${TMP_BIN} ${BIN_URL}
}

# --- setup permissions and move binary to system directory ---
setup_binary() {
    chmod 755 ${TMP_BIN}
    info "Installing qcadmin to ${BIN_DIR}/qcadmin"
    $SUDO chown root:root ${TMP_BIN}
    $SUDO mv -f ${TMP_BIN} ${BIN_DIR}/qcadmin
		[ -f "${BIN_DIR}/q" ] && (
			$SUDO rm -f ${BIN_DIR}/q
		)
		info "Create qcadmin soft link ${BIN_DIR}/q"
		$SUDO ln -s ${BIN_DIR}/qcadmin ${BIN_DIR}/q
		info "Installation is complete. Use q --help"
}

# --- download and verify qcadmin ---
download_and_verify() {
    setup_verify_arch
    verify_downloader curl || verify_downloader wget || fatal 'Can not find curl or wget for downloading files'
    setup_tmp
    setup_quickon
    get_release_version
		# Skip download if qcadmin binary exists, support upgrade
    download_binary
    setup_binary
}

# --- run the install process --
{
	setup_env "$@"
	download_and_verify
}
