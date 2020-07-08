class Apigeecli < Formula
    desc "This is a tool to interact with Apigee APIs for Apigee hybrid."
    homepage "https://github.com/srinandan/apigeecli"
    version "v1.7.2"
    bottle :unneeded
  
    if OS.mac?
      url "https://github.com/srinandan/apigeecli/releases/download/v1.7.2/apigeecli_v1.7.2_Darwin_x86_64.zip"
      sha256 "38018708ab46a6ebfdef2094b5dabcb57d0323cb3ee7185f7930cbf8bdec7379"
    elsif OS.linux?
      if Hardware::CPU.intel?
        url "https://github.com/srinandan/apigeecli/releases/download/v1.7.2/apigeecli_v1.7.2_Linux_x86_64.zip"
        sha256 "d16183ed3da532e2710c66a3e72b2992cb863cbe0c4948b8b44e8f95eef257bb"
      end
    end
  
    def install
      bin.install "apigeecli"
    end
  
    test do
      system "#{bin}/apigeecli --version"
    end
  end