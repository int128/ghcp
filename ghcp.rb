class Ghcp < Formula
  desc "Commit files to a repository using GitHub API without git configuration"
  homepage "https://github.com/int128/ghcp"
  url "https://github.com/int128/ghcp/releases/download/{{ env "VERSION" }}/ghcp_darwin_amd64.zip"
  version "{{ env "VERSION" }}"
  sha256 "{{ sha256 .darwin_amd64_archive }}"

  def install
    bin.install "ghcp"
  end

  test do
    system "#{bin}/ghcp -h"
  end
end
