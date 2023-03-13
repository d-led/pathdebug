# pathdebug

a simple interactive [tui](https://en.wikipedia.org/wiki/Text-based_user_interface) to debug path list environment variables

status: *untested, spike quality, use at your own peril*

## Install

local: `go install .``

from github: `go install github.com/d-led/pathdebug@main`

## Usage

```bash
pathdebug {nameOfEnvironmentVariable}
```

### Example

```bash
export PATH=/sbin:/a:/b:/a:/c:/d:/e:/f:/g
pathdebug PATH
```

&darr;

```text
tap Esc/q/Ctrl-C to quit, <-/-> to paginate
+-------+-----+-------+
| DUP # | BAD | PATH  |
+-------+-----+-------+
|       |     | /sbin |
| 2     | X   | /a    |
|       | X   | /b    |
| 2     | X   | /a    |
+-------+-----+-------+
  •○○
```
