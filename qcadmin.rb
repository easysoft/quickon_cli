class Qcadmin < Formula
    desc "qcadmin is an open-source lightweight cli tool for managing quickon."
    homepage "https://github.com/easysoft/quickon_cli"
    version "3.0.0-beta.1"

    on_macos do
      if Hardware::CPU.arm?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_arm64"
        sha256 "287d16f7f7852239759d2e59a3c169b243a1211cdaac5fdc7590abf97af8b947"

        def install
            bin.install "qcadmin_darwin_arm64" => "qcadmin"
        end
      end

      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_amd64"
        sha256 "05ce5f60b5d57b3913eb7897c5c9a8def40e33178936a2d4124bf9f98ce46566"

        def install
            bin.install "qcadmin_darwin_amd64" => "qcadmin"
        end
      end
    end

    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_amd64"
        sha256 "a2070c478488732cf7c829745005fe0b0604682a294ee833d1e215758be0750d"

        def install
            bin.install "qcadmin_linux_amd64" => "qcadmin"
        end
      end

      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_arm64"
        sha256 "d1ed9abf7e139bc6b0babb72408a08c61a6fa99a81444b85c779cc67a43a221c"

        def install
            bin.install "qcadmin_linux_arm64" => "qcadmin"
        end
      end
    end
end
