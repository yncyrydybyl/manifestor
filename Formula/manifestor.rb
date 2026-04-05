class Manifestor < Formula
  desc "Grab the latest file from ~/Downloads with a clean name"
  homepage "https://github.com/yncyrydybyl/manifestor"
  url "https://github.com/yncyrydybyl/manifestor/archive/refs/tags/v0.1.1.0.tar.gz"
  sha256 "35472166abc15d4048f5890a26dccd3d6101f3aac67ce35be9770046993d7b62"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build",
           "-ldflags", "-s -w -X main.version=#{version}",
           "-o", bin/"m",
           "./cmd/m"
    bin.install_symlink "m" => "mm"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/m --version")
  end
end
