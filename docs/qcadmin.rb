class Qcadmin < Formula
    desc "qcadmin is an open-source lightweight cli tool for managing quickon."
    homepage "https://github.com/easysoft/quickon_cli"
    version "2.0.9"

    on_macos do
      if Hardware::CPU.arm?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_arm64"
        sha256 "bb5ddd24e293b3fafe26b6baf82a8cebbde3a7223a949717a5e00ebdf09a340e"

        def install
            bin.install "qcadmin_darwin_arm64" => "qcadmin"
        end
      end

      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_amd64"
        sha256 "d57d179a4a800937c1fe09d478229943691aaf2a9112369f45822bc89cfa634c"

        def install
            bin.install "qcadmin_darwin_amd64" => "qcadmin"
        end
      end
    end

    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_amd64"
        sha256 "71da2e7d5e2dd3b88e4b92ae8f61fbc5b187645e03e556519c4275d59d93d202"

        def install
            bin.install "qcadmin_linux_amd64" => "qcadmin"
        end
      end

      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_arm64"
        sha256 "0828391c5a07b8c1592bb6c5d2d6b5b4919e41a6fe9470a41640aa8985d7291c"

        def install
            bin.install "qcadmin_linux_arm64" => "qcadmin"
        end
      end
    end
end
