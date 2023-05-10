#!/usr/bin/env bash

version=$(cat VERSION)

version=$(cat VERSION)
# shellcheck disable=SC2002
macosAMD64sha=$(cat dist/checksums.txt | grep qcadmin_darwin_amd64 | awk '{print $1}')
# shellcheck disable=SC2002
macosARM64sha=$(cat dist/checksums.txt | grep qcadmin_darwin_arm64| awk '{print $1}')
# shellcheck disable=SC2002
linuxAMD64sha=$(cat dist/checksums.txt | grep qcadmin_linux_amd64 | awk '{print $1}')
# shellcheck disable=SC2002
linuxARM64sha=$(cat dist/checksums.txt | grep qcadmin_linux_arm64 | awk '{print $1}')

cat > docs/qcadmin.rb <<EOF
class Qcadmin < Formula
    desc "qcadmin is an open-source lightweight cli tool for managing quickon."
    homepage "https://github.com/easysoft/quickon_cli"
    version "${version}"

    on_macos do
      if Hardware::CPU.arm?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_arm64"
        sha256 "${macosARM64sha}"

        def install
            bin.install "qcadmin_darwin_arm64" => "qcadmin"
        end
      end

      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_amd64"
        sha256 "${macosAMD64sha}"

        def install
            bin.install "qcadmin_darwin_amd64" => "qcadmin"
        end
      end
    end

    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_amd64"
        sha256 "${linuxAMD64sha}"

        def install
            bin.install "qcadmin_linux_amd64" => "qcadmin"
        end
      end

      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_arm64"
        sha256 "${linuxARM64sha}"

        def install
            bin.install "qcadmin_linux_arm64" => "qcadmin"
        end
      end
    end
end
EOF
