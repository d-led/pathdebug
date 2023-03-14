# pathdebug

a simple interactive [tui](https://en.wikipedia.org/wiki/Text-based_user_interface) to debug path list environment variables

## Install

from source:

```bash
go install
```

from github:

```bash
go install github.com/d-led/pathdebug@latest
```

replace `latest` with the desired/latest commit hash if you had the tool installed already.

## Usage

```bash
pathdebug {nameOfEnvironmentVariable}
```

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

### Direct Output

```bash
export SOME_PATH='/sbin:~/.bashrc:/a:/b:/a:/c:/d:/e:/f:/g'
pathdebug SOME_PATH -o table
```

see help for other formats
