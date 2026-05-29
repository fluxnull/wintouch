# WinTouch

Windows-native `touch.exe` inspired by BSD `touch`.

WinTouch creates files and edits Windows file timestamps from the command line. Unlike a POSIX-only clone, it treats the three Windows-visible timestamps as first-class fields:

- Created
- Modified
- Accessed

## Command

```cmd
touch [options] file [file ...]
```

## Features

- Updates Created, Modified, and Accessed timestamps.
- Default behavior updates all three Windows timestamps.
- Select individual timestamps with `-a`, `-m`, and `-c`.
- Supports grouped boolean flags such as `-amc`, `-np`, and `-rm`.
- Sets readable dates with `-d <DATE>`.
- Sets BSD compact timestamps with `-t <TIMESTAMP>`.
- Synchronizes timestamps from another file with `-s / --sync <FILE>`.
- Applies relative timestamp offsets with `-e / --edit <OFFSET>`.
- Supports no-create mode with `-n / --no-create`.
- Creates missing parent directories with `-p / --parents`.
- Supports wildcard targets.
- Supports files and directories.
- Supports Unicode paths.
- Supports long Windows paths.
- Supports reparse controls with `-f / --follow-reparse` and `-i / --ignore-reparse`.
- Supports UAC relaunch with `-r / --runas`.

## Examples

Create a file if missing, or update Created, Modified, and Accessed:

```cmd
touch file.txt
```

Set Modified only:

```cmd
touch -m -d "2026-05-28 21:30:00" file.txt
```

Set Created, Modified, and Accessed:

```cmd
touch -amc -d "2026-05-28 21:30:00" file.txt
```

Synchronize timestamps from another file:

```cmd
touch -s reference.txt target.txt
```

Add one minute and thirty seconds to all selected timestamps:

```cmd
touch -e +0130 file.txt
```

Create parent directories and the file:

```cmd
touch -p logs\build\output.txt
```

Run elevated through UAC:

```cmd
touch -r -m -d "2027-06-29 22:31:01" protected.txt
```

## Help

Full command help is in [`HELP.txt`](HELP.txt).

```cmd
touch /?
```

## Build

```cmd
build.cmd
```

The Windows AMD64 executable is written to:

```text
bin\touch.exe
```

## Test

Run the CMD-only behavior script:

```cmd
test\test-wintouch-cmd.cmd bin\touch.exe
```

The test runner creates:

```text
test-output\wintouch_cmd_report_<timestamp>.txt
test-output\wintouch_cmd_cases_<timestamp>\
```

See [`TESTING.md`](TESTING.md).

## Project scope

WinTouch is Windows-first. It is not intended to be a POSIX compatibility layer.

## Provenance

WinTouch is inspired by BSD `touch`, but this repository does not redistribute FreeBSD source files.

See [`PROVENANCE.md`](PROVENANCE.md) for upstream reference links.

## License

WinTouch original code and documentation are released under the Zero-Clause BSD license (`0BSD`).

See [`LICENSE`](LICENSE).
