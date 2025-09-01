class Icignore < Formula
  desc "Exclude iCloud folders via .nosync + symlink"
  homepage "https://github.com/mathis-lambert/icloud-ignore"
  license "MIT"
  version "1.0.0"

  # Stable release (fill in on tag):
  # url "https://github.com/mathis-lambert/icloud-ignore/archive/refs/tags/v0.1.0.tar.gz"
  # sha256 "TBD"

  head "https://github.com/mathis-lambert/icloud-ignore.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X main.version=\#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags), "./cmd/icignore"
  end

  test do
    assert_match version.to_s, shell_output("\#{bin}/icignore version")
  end
end

