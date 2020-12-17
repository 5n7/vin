![banner](./misc/banner.png)

<p align="center">
  <a href="https://github.com/skmatz/vin/actions?query=workflow%3Aci">
    <img
      src="https://github.com/skmatz/vin/workflows/ci/badge.svg"
      alt="ci"
    />
  </a>
  <a href="https://github.com/skmatz/vin/actions?query=workflow%3Arelease">
    <img
      src="https://github.com/skmatz/vin/workflows/release/badge.svg"
      alt="release"
    />
  </a>
  <a href="https://goreportcard.com/report/github.com/skmatz/vin">
    <img
      src="https://goreportcard.com/badge/github.com/skmatz/vin"
      alt="go report card"
    />
  </a>
  <a href="./LICENSE">
    <img
      src="https://img.shields.io/github/license/skmatz/vin"
      alt="license"
    />
  </a>
  <a href="./go.mod">
    <img
      src="https://img.shields.io/github/go-mod/go-version/skmatz/vin"
      alt="go version"
    />
  </a>
  <a href="https://github.com/skmatz/vin/releases/latest">
    <img
      src="https://img.shields.io/github/v/release/skmatz/vin"
      alt="release"
    />
  </a>
</p>

<p align="center">
<b><a href="#description">Description</a></b>
|
<b><a href="#features">Features</a></b>
|
<b><a href="#installation">Installation</a></b>
|
<b><a href="#usage">Usage</a></b>
|
<b><a href="#configuration">Configuration</a></b>
|
<b><a href="#references">References</a></b>
</p>

## Description

`vin` is the CLI to install applications from GitHub Releases.

In recent years, many useful CLI tools have been made by Go, Rust, etc.  
Many of them are available as pre-built binaries in GitHub Releases, but we have to manipulate the browser to download them, decompress the compressed files, and move the executables to the appropriate path.

`vin` makes it easy to manage all at once by writing a list of applications in a TOML file.  
No longer do we have to select suitable assets for our machines from GitHub Releases, no longer do we have to look up the options for the `tar` command, and no longer do we have to be pained when we setup a new machine.

## Features

- Easy to use
- TOML-based configuration

## Installation

Download the binary from [GitHub Releases](https://github.com/skmatz/vin/releases).  
Unfortunately, the first time you install `vin`, you need to open the browser.

Or, if you have Go, you can install `vin` with the following command.

```bash
go get github.com/skmatz/vin/...
```

After installation, add `~/.vin/bin` to your `$PATH`.

```bash
export PATH="$HOME/.vin/bin:$PATH"
```

## Usage

First, put the following TOML file in `~/.config/vin/vin.toml`.

```toml
[[app]]
repo = "cli/cli"
```

You can set this with the following command.

```bash
# darwin
mkdir -p ~/Library/Application\ Support/vin
vin example > ~/Library/Application\ Support/vin/vin.toml

# linux
mkdir -p ~/.config/vin
vin example > ~/.config/vin/vin.toml

# windows (confirmed on Windows Terminal)
mkdir -p ~\AppData\Roaming\vin\vin.toml
vin example > ~\AppData\Roaming\vin\vin.toml
```

Yes, all you have to do is run the following command.

```bash
vin get
```

## Configuration

```toml
[[app]]
# repo is the GitHub repository name in "owner/repo" format.
repo = "cli/cli"

# tag is the tag on GitHub.
# If empty, it is treated as "latest".
tag = "v1.2.0"

# keywords is a list of keywords for selecting suitable assets from multiple assets.
# If empty, it is treated as [$GOOS, $GOARCH].
keywords = ["amd64", "linux"]

# name is the name of the executable file.
# If empty, it is treated as the original name.
name = "gh"

# hosts is a list of host names.
# If empty, it is treated as any hosts.
hosts = ["awesome-machine"]

# priority is the priority of the application.
priority = 3

# command is the command to run after installation.
command = """
gh completion -s zsh > ~/.zfunc/_gh
"""
```

## References

Thanks to [ghg](https://github.com/Songmu/ghg) for the idea.
