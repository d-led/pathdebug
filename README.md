# pathdebug

a simple interactive [tui](https://en.wikipedia.org/wiki/Text-based_user_interface) to debug path list environment variables

status: *untested, spike quality, use at your own peril*

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

### Example

```bash
SOME_PATH=/sbin:/a:/b:/a:/c:/d:/e:/f:/g \
pathdebug SOME_PATH
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
