class Smgr < Formula
  desc "Manage Semantic Versioning compliant versions"
  homepage "https://github.com/bluepr-nt/semver-manager"
  head "https://github.com/bluepr-nt/semver-manager.git", branch: "main"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args, "./src/cmd/smgr"
  end

  test do
    assert_match "0.0.1", shell_output("#{bin}/smgr increment --level patch")
  end
end