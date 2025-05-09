class Qcadmin < Formula
    desc "qcadmin is an open-source lightweight cli tool for managing quickon."
    homepage "https://github.com/easysoft/quickon_cli"
    version "4.3.1"

    on_macos do
      if Hardware::CPU.arm?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_arm64"
        sha256 "155d79b0953edc028bc5770e62bbfe256d14e1e4aa572e8fe43f9c071794c166"

        def install
            bin.install "qcadmin_darwin_arm64" => "qcadmin"
        end
      end

      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_amd64"
        sha256 "5a1a56a0cdee4b164afc31d68acbe4ee8dbcec4b69dfa49ecba6121a18fb62a9"

        def install
            bin.install "qcadmin_darwin_amd64" => "qcadmin"
        end
      end
    end

    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_amd64"
        sha256 "e6f66ec8a27fea3ec7a9aa942c5137b043f9932fb79ca3c0905609ef76ec3f91"

        def install
            bin.install "qcadmin_linux_amd64" => "qcadmin"
        end
      end

      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_arm64"
        sha256 "921c6e6c441799d4820e0cc1d0c6be5fda84edf48cf48b52e3dad33f97c7c36c"

        def install
            bin.install "qcadmin_linux_arm64" => "qcadmin"
        end
      end
    end
end
