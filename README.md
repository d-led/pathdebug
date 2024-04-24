# pathdebug

a simple interactive & one-shot tool to debug path list environment variables

[![Go Report Card](https://goreportcard.com/badge/github.com/d-led/pathdebug)](https://goreportcard.com/report/github.com/d-led/pathdebug)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fd-led%2Fpathdebug.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fd-led%2Fpathdebug?ref=badge_shield)
[![Test](https://github.com/d-led/pathdebug/actions/workflows/test.yml/badge.svg)](https://github.com/d-led/pathdebug/actions/workflows/test.yml)
[![CodeQL](https://github.com/d-led/pathdebug/actions/workflows/codeql.yml/badge.svg)](https://github.com/d-led/pathdebug/actions/workflows/codeql.yml)

```bash
pathdebug {EnvironmentVariableName}
```

## Features

- detect duplicate expanded path entries
- detecting non-existent paths and files
- best effort of **finding out where the path is set**
  - [known startup script locations](common/known_paths.go), such as `.bashrc`
  - `/etc/path.d` path lists
  - recursive follow-ups of `source $file`-like statements in the above scripts
- interactive mode
- noninteractive formatted output

## Install

### Homebrew

```bash
brew tap d-led/d-led
brew install pathdebug
```

### Chocolatey

```shell
choco install pathdebug
```

see package at [Chocolatey](https://community.chocolatey.org/packages/pathdebug)

### From Source

```bash
go install
```

from github:

```bash
go install github.com/d-led/pathdebug@latest
```

Alternatively, download the binary from [Github releases](https://github.com/d-led/pathdebug/releases/).

On OSX you might need to run `sudo xattr -r com.apple.quarantine pathdebug` on the downloaded binary.

### Download Latest Release

OSX:

```bash
curl -s -L https://github.com/d-led/pathdebug/releases/latest/download/pathdebug_Darwin_all.tar.gz | tar xvz - pathdebug
```

Linux: same procedure, but substitute the artifact name with the appropriate one.

## Usage

help:

```bash
pathdebug --help
```

### Interactive

```bash
export SOME_PATH='/sbin:~/.bashrc:/a:/b:/a:/c:/d:/e:/f:/g'
pathdebug SOME_PATH
```

&darr;

```text
tap Esc/q/Ctrl-C to quit, <-/-> to paginate
+---+--------+-----+-----------+
| # | DUP[#] | BAD | PATH      |
+---+--------+-----+-----------+
| 1 |        |     | /sbin     |
| 2 |        | F   | ~/.bashrc |
| 3 | 5      | X   | /a        |
| 4 |        | X   | /b        |
| 5 | 3      | X   | /a        |
| 6 |        | X   | /c        |
+---+--------+-----+-----------+
  •○
```

and for the path variable with source detection:

```text
+---+--------+-----+----------------------------------+-----------------+
| # | DUP[#] | BAD | PATH                             | ±SOURCES        |
+---+--------+-----+----------------------------------+-----------------+
| 1 |        |     | /Users/d/.opam/default/bin       |                 |
| 2 |        |     | /Users/d/.local/share/ponyup/bin | ~/.zshrc        |
| 3 |        |     | /usr/local/sbin                  | ~/.zshrc        |
| 4 |        |     | /opt/homebrew/bin                | ~/.bash_profile |
| 5 |        |     | /opt/homebrew/sbin               |                 |
+---+--------+-----+----------------------------------+-----------------+
  •○○○○○
```

### Direct Output

```bash
pathdebug PATH -o table
```

see help for other formats

## Limitations

- finding the sources where the `PATH` variable entries are set is best-effort, and is done only for the `PATH` environment variable at present
- path set by executables currently not tracked as sources
- CSV output may lag in output completeness

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fd-led%2Fpathdebug.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fd-led%2Fpathdebug?ref=badge_large)
