class Qcadmin < Formula
    desc "qcadmin is an open-source lightweight cli tool for managing quickon."
    homepage "https://github.com/easysoft/quickon_cli"
    version "3.2.1"

    on_macos do
      if Hardware::CPU.arm?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_arm64"
        sha256 "c674ae9cb63a2de56dd6dae12f23ede42e36466658deae19c976ecfafc3889fc"

        def install
            bin.install "qcadmin_darwin_arm64" => "qcadmin"
        end
      end

      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_darwin_amd64"
        sha256 "7a9bf62f5d4dbbcddc540d24626848fac021127b37bbf7d38089b4b54d7cde9a"

        def install
            bin.install "qcadmin_darwin_amd64" => "qcadmin"
        end
      end
    end

    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_amd64"
        sha256 "e8f011325146df0dffa6b7a3981f16cbcede60fa85f1aeb0d79fe4e8ce8ec525"

        def install
            bin.install "qcadmin_linux_amd64" => "qcadmin"
        end
      end

      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/easysoft/quickon_cli/releases/download/v#{version}/qcadmin_linux_arm64"
        sha256 "dc4cb257dfbed9f9357670e816896463b04d38ddd704d3b0cb7c4490bf24b1b2"

        def install
            bin.install "qcadmin_linux_arm64" => "qcadmin"
        end
      end
    end
end
