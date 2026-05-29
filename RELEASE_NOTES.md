# WinTouch v0.3.0

Initial Windows-first GitHub publishing package.

## Release identity

- Project: WinTouch
- Binary: `touch.exe`
- Version: `v0.3.0`
- Platform: Windows
- Architecture: AMD64
- License: 0BSD

## Features

- Updates Created, Modified, and Accessed timestamps.
- Default behavior updates all three Windows timestamps.
- Timestamp selectors:
  - `-a / --accessed`
  - `-m / --modified`
  - `-c / --created`
- Boolean short option grouping:
  - `-am`
  - `-ac`
  - `-mc`
  - `-amc`
  - `-np`
  - `-rm`
- Readable date input with `-d / --date <DATE>`.
- BSD compact timestamp input with `-t <TIMESTAMP>`.
- Timestamp synchronization with `-s / --sync <FILE>`.
- Relative timestamp editing with `-e / --edit <OFFSET>`.
- Parent directory creation with `-p / --parents`.
- No-create mode with `-n / --no-create`.
- Reparse controls:
  - `-f / --follow-reparse`
  - `-i / --ignore-reparse`
- UAC relaunch with `-r / --runas`.
- Wildcard support.
- Directory timestamp support.
- Unicode path support.
- Long Windows path support.

## Validation

Validated on Windows using the included CMD-only test runner.

Covered behavior includes:

- Help aliases
- Version aliases
- Timestamp selectors
- Boolean clusters
- Separated value tokens
- Date formats
- Compact timestamp format
- Sync from referenced file
- Positive, negative, and unsigned offsets
- Conflict rejection
- Wildcards
- Directories
- Unicode paths
- Long paths
- Runas clustered form

## Provenance

WinTouch is inspired by BSD `touch` behavior. FreeBSD source files are not redistributed in this package. Upstream reference links are listed in [`PROVENANCE.md`](PROVENANCE.md).

## License

WinTouch original code and documentation are released under `0BSD`.
