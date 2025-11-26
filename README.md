# ytop - fast resource monitor for macOS and Linux

# Screenshot Linux

<img
  class="centered"
  src="doc/screenshot-linux.png"
  alt="ytop screenshot on Linux"
/>

# Screenshot macOS

<img
  class="centered"
  src="doc/screenshot-macos.png"
  alt="ytop screenshot on macOS"
/>

# How to run

```text
$ make
```

# How to develop

You need a number of build tools to work on `ytop`. Below, I show how
to install this toolchain on Arch Linux. If you are on a different
platform, consult the appropriate documentation, ask Google or a
friend. The `go install` commands work on all platforms, as long as
you have `go` itself installed:

```text
# pacman -S go
# pacman -S gopls
$ go install honnef.co/go/tools/cmd/staticcheck@latest
$ go install golang.org/x/vuln/cmd/govulncheck@latest
```

# License

`ytop` is free software, GPLv3, see [LICENSE](LICENSE).
