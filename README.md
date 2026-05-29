# WinTouch

WinTouch is a Windows-native `touch.exe` for creating files and controlling file timestamps from CMD.

It is inspired by BSD `touch`, but it is not a POSIX compatibility shim. WinTouch is built around Windows filesystem behavior and treats the three Windows-visible timestamps as first-class fields:

* Created
* Modified
* Accessed

Classic `touch` is mainly used to create files and update file times. WinTouch keeps that core behavior, then extends it for Windows-native scripting, testing, metadata repair, and batch timestamp control.

## What it does

* Creates empty files when they do not exist.
* Updates timestamps on existing files.
* Defaults to updating Created, Modified, and Accessed.
* Selects individual timestamps with `-a`, `-m`, and `-c`.
* Sets exact readable dates with `-d <DATE>`.
* Sets BSD compact timestamps with `-t <TIMESTAMP>`.
* Synchronizes timestamps from another file with `-s / --sync <FILE>`.
* Applies relative timestamp offsets with `-e / --edit <OFFSET>`.

## Windows-native enhancements

- **Created timestamp support**  
  Windows exposes file creation time as a normal visible timestamp. WinTouch makes Created as scriptable as Modified and Accessed.

- **All-three timestamp default**  
  By default, WinTouch updates Created, Modified, and Accessed together instead of only Modified and Accessed.

- **Grouped boolean flags**  
  Short options can be grouped for fast CMD usage, such as `-amc`, `-np`, and `-rm`.

- **Long path support**  
  Handles Windows long paths beyond legacy `MAX_PATH` limits using native Windows path handling.

- **Unicode path support**  
  Works with Unicode filenames and paths through Windows wide-character APIs.

- **Directory timestamp support**  
  Updates timestamps on directories as well as files.

- **Wildcard targets**  
  Expands CMD-style wildcard targets such as `*.txt`, `*.go`, and `logs\*.json`.

- **Parent directory creation**  
  `-p / --parents` creates missing parent folders before touching the target file.

- **No-create mode**  
  `-n / --no-create` updates only existing files and skips missing paths.

- **Reparse-point control**  
  `-f / --follow-reparse` follows reparse points to their targets, while `-i / --ignore-reparse` skips reparse points entirely.

- **UAC relaunch**  
  `-r / --runas` relaunches through Windows UAC for protected paths.

- **CMD-first behavior**  
  Designed for direct use from CMD and batch files, with Windows-style help aliases such as `/?`, `/H`, and `/HELP`.

## Project scope

WinTouch is Windows-first. It is not intended to be a POSIX compatibility layer.

## Provenance

WinTouch is inspired by BSD `touch`, but this repository does not redistribute FreeBSD source files.

See [`PROVENANCE.md`](PROVENANCE.md) for upstream reference links.

## License

WinTouch original code and documentation are released under the Zero-Clause BSD license (`0BSD`).

See [`LICENSE`](LICENSE).
