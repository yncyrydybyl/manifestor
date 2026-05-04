class Manifestor < Formula
  desc "Grab the latest file from ~/Downloads with a clean name"
  homepage "https://github.com/yncyrydybyl/manifestor"
  url "https://github.com/yncyrydybyl/manifestor/archive/refs/tags/v0.1.3.0.tar.gz"
  sha256 "9a259fdad343e1d0124ccf210aa8d018c05bb30ed432fc3d90132d6cae2d1a5e"
  license "MIT"
  head "https://github.com/yncyrydybyl/manifestor.git", branch: "main"

  depends_on "go" => :build

  conflicts_with "m-cli", because: "both install a `m` binary"

  def install
    ldflags = "-s -w -X main.version=#{version}"
    system "go", "build",
           *std_go_args(output: bin/"manifestor", ldflags: ldflags),
           "./cmd/m"
    bin.install_symlink "manifestor" => "m"
    bin.install_symlink "manifestor" => "mm"

    generate_completions_from_executable(bin/"manifestor", "completion")
  end

  def caveats
    <<~EOS
      manifestor installs three commands:
        manifestor   canonical name (always available)
        m            short alias
        mm           force mode (skips the staleness check)

      If commands aren't found after install, ensure Homebrew's bin is on
      your PATH:
        # Apple Silicon
        eval "$(/opt/homebrew/bin/brew shellenv)"
        # Intel
        eval "$(/usr/local/bin/brew shellenv)"

      Then refresh your shell command cache:
        hash -r        # bash
        rehash         # zsh
    EOS
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/manifestor --version")
    assert_match version.to_s, shell_output("#{bin}/m --version")
    assert_match version.to_s, shell_output("#{bin}/mm --version")
    assert_match "rainbow-beam", shell_output("#{bin}/manifestor --list-anims")
    assert_match "COMP_WORDS", shell_output("#{bin}/manifestor completion bash")
  end
end
