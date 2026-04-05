class Manifestor < Formula
  desc "Grab the latest file from ~/Downloads with a clean name"
  homepage "https://github.com/yncyrydybyl/manifestor"
  url "https://github.com/yncyrydybyl/manifestor/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "PUT_REAL_SHA256_HERE"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build",
           "-ldflags", "-s -w -X main.version=#{version}",
           "-o", bin/"manifestor",
           "./cmd/manifestor"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/manifestor --version")
  end
end
