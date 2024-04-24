class Qcadmin < Formula
    desc "qcadmin is an open-source lightweight cli tool for managing quickon."
    homepage "https://github.com/easysoft/quickon_cli"
    version "3.0.20"

    on_macos do
      if Hardware::CPU.arm?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_arm64"
        sha256 "657794601a3b0afb2a32ba980f93f3293fc0ed3c5d8d79358b432c91558646fe"

        def install
            bin.install "qcadmin_darwin_arm64" => "qcadmin"
        end
      end

      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_amd64"
        sha256 "bade98e6b400fb3d4d3d7b1b8d4b29077dc4037310ed8c651be218027221439b"

        def install
            bin.install "qcadmin_darwin_amd64" => "qcadmin"
        end
      end
    end

    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_amd64"
        sha256 "f946f5431c0ab557760c9875678126d09a84b2334374ce5e129c2e6c660efb3d"

        def install
            bin.install "qcadmin_linux_amd64" => "qcadmin"
        end
      end

      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_arm64"
        sha256 "5e75c6530ed6333a6fa576d4a717fc67dcc60763aff024d93222e8031df0a4c6"

        def install
            bin.install "qcadmin_linux_arm64" => "qcadmin"
        end
      end
    end
end
