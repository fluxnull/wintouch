# Testing

WinTouch includes a CMD-only behavior test runner.

No PowerShell is required.

## Run

From the repository root:

```cmd
test\test-wintouch-cmd.cmd bin\touch.exe
```

Or with an absolute path:

```cmd
test\test-wintouch-cmd.cmd F:\path\to\touch.exe
```

## Output

The test runner creates:

```text
test-output\wintouch_cmd_report_<timestamp>.txt
test-output\wintouch_cmd_cases_<timestamp>\
```

The report records:

- test section
- command label
- exact command
- stdout/stderr
- `ERRORLEVEL`
- file state from `dir /T:C`, `dir /T:W`, and `dir /T:A`
- attributes
- reparse query output where applicable

## Expected non-zero results

Some tests intentionally return non-zero `ERRORLEVEL` values because they validate rejection behavior:

- attached value forms such as `-e0130`
- value options inside boolean clusters such as `-amct`
- conflicting reparse options `-i -f`
- conflicting base time sources such as `-s` with `-d`

## Reparse tests

Symlink/reparse tests require permission to create symlinks.

Enable Developer Mode or run CMD as Administrator if `mklink` fails.

## Runas tests

`-r / --runas` may trigger UAC. Manual runas commands are included in the test report rather than forced by default.
