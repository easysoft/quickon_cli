#!/bin/sh
# Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

# Source code is available at https://github.com/easysoft/quickon_cli

# SCRIPT_COMMIT_SHA="bd770bd54308aad22b2a2e3ee585c2693a49a6de"
# SCRIPT_DATA="Thu Apr  3 11:27:15 CST 2025"

# Usage:
#   curl ... | ENV_VAR=... sh -
#       or
#   ENV_VAR=... ./install.sh
#
# Example:
#   Installing DevOPS with ZenTao:
#     curl ... | DEVOPS_TYPE="" sh -
#   - DEVOPS_TYPE
#     Install Type when install Zentao DevOPS.
#     Defaults to '', support 'max', 'biz', 'ipd'
#   - DEVOPS_VERSION
#     Install Version when install Zentao DevOPS.
#     Defaults to ''
#   - INSTALL_DOMAIN
#   - DEBUG
#     If set, print debug information
#   - STORAGE_TYPE
#     Storage Type when install Zentao DevOPS default use local as storage provider.
#     Defaults to '', support 'local', 'nfs'
#   - EX_DB_HOST
#     External Database Host when install Zentao DevOPS.
#     Defaults to ''
#   - EX_DB_PORT
#     External Database Port when install Zentao DevOPS.
#     Defaults to '3306'
#   - EX_DB_USER
#     External Database User when install Zentao DevOPS.
#     Defaults to ''
#   - EX_DB_PASSWORD
#     External Database Password when install Zentao DevOPS.
#     Defaults to ''


set -e
set -o noglob

[ -n "${DEBUG:+1}" ] && set -x

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

COS_URL=https://pkg.zentao.net/cli/devops
STORAGE_TYPE=${STORAGE_TYPE:-}

# --- define needed environment variables ---
setup_env() {
    # --- use sudo if we are not already root ---
    SUDO=sudo
    if [ $(id -u) -eq 0 ]; then
        SUDO=
    fi
    BIN_DIR=/usr/local/bin
    if ! $SUDO sh -c "touch ${BIN_DIR}/z-ro-test && rm -rf ${BIN_DIR}/z-ro-test"; then
      if [ -d /opt/bin ]; then
        BIN_DIR=/opt/bin
      fi
    fi
}

# --- verify an executable z binary is installed ---
verify_z_is_executable() {
    if [ ! -x ${BIN_DIR}/z ]; then
        fatal "Executable z binary not found at ${BIN_DIR}/z"
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
    TMP_DIR=$(mktemp -d -t z-install.XXXXXXXXXX)
    TMP_HASH=${TMP_DIR}/z.hash
    TMP_BIN=${TMP_DIR}/z.bin
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
  if [ -z "$VERSION" ]; then
    VERSION="stable"
  fi
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
    BIN_URL=${COS_URL}/${VERSION}/qcadmin_linux_${SUFFIX}
    info "Downloading binary"
    download ${TMP_BIN} ${BIN_URL}
}

# --- setup permissions and move binary to system directory ---
setup_binary() {
    chmod 755 ${TMP_BIN}
    info "Installing z to ${BIN_DIR}/z"
    $SUDO chown root:root ${TMP_BIN}
    $SUDO mv -f ${TMP_BIN} ${BIN_DIR}/qcadmin
    [ -f "${BIN_DIR}/z" ] && (
        $SUDO rm -f ${BIN_DIR}/z
    )
    info "Create soft link ${BIN_DIR}/z"
    $SUDO ln -s ${BIN_DIR}/qcadmin ${BIN_DIR}/z
    info "Installation is complete. Use z --help"
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

# --- install zentao devops
install_zentao_devops() {
  INSTALL_COMMAND="${BIN_DIR}/z init --provider devops"
  if [ -n "${INSTALL_DOMAIN}" ]; then
    INSTALL_COMMAND="${BIN_DIR}/z init --provider devops --domain ${INSTALL_DOMAIN}"
  fi
  if [ -n "${DEVOPS_TYPE}" ]; then
    INSTALL_COMMAND="${INSTALL_COMMAND} --type ${DEVOPS_TYPE}"
  fi
  if [ -n "${DEVOPS_VERSION}" ]; then
    INSTALL_COMMAND="${INSTALL_COMMAND} --version ${DEVOPS_VERSION}"
  fi
  if [ "${STORAGE_TYPE}" = "nfs" ]; then
    INSTALL_COMMAND="${INSTALL_COMMAND} --storage nfs"
  fi
  if [ -n "${EX_DB_HOST}" ] && [ -n "${EX_DB_PASSWORD}" ]; then
    INSTALL_COMMAND="${INSTALL_COMMAND} --ext-db-host ${EX_DB_HOST} --ext-db-password ${EX_DB_PASSWORD}"
    if [ -n "${EX_DB_PORT}" ] && [ "${EX_DB_PORT}" != "3306" ]; then
      INSTALL_COMMAND="${INSTALL_COMMAND} --ext-db-port ${EX_DB_PORT}"
    fi
    if [ -n "${EX_DB_USER}" ] && [ "${EX_DB_USER}" != "root" ]; then
      INSTALL_COMMAND="${INSTALL_COMMAND} --ext-db-user ${EX_DB_USER}"
    fi
  fi
  if [ -n "${DEBUG}" ]; then
    INSTALL_COMMAND="${INSTALL_COMMAND} --debug"
  fi
  eval "$INSTALL_COMMAND"
}

# --- run the install process --
{
  setup_env
  download_and_verify
  install_zentao_devops
}
