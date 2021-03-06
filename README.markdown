# Dupes

[![Build Status](https://img.shields.io/travis/yesmar/dupes.svg)](https://travis-ci.org/yesmar/dupes)
&nbsp;
![Platforms](https://img.shields.io/badge/platform-macOS%20|%20Linux%20|%20Windows-lightgrey.svg)

The `dupes` command identifies duplicate files within the user-specified set of target files and directories. It does this by walking the specified file systems, computing and comparing SHA-256 hashes of the files it encounters. Any hash that has been seen previously is considered a duplicate.

Important considerations:

- The `dupes` command does not perform any destructive operations.
- Zero-length files serve as a trigger mechanism on many operating systems. Removing such files could have negative effects, so `dupes` completely ignores empty files.
- Symbolic links *will* be followed, which could lead to trouble should a cycle be encountered. `dupes` assumes that you know what you are doing and understand the file systems you intend to run it on.

The author uses `dupes` to identify and delete duplicate images downloaded while satisfying his [imgur](https://imgur.com) addiction.

## Install

```bash
go get github.com/yesmar/dupes
```

## Usage

```bash
Usage: dupes [flags] <pathname> […]
```

Currently, `dupes` understands the `-verbose` and `-version` flags. Use the `-verbose` flag to enable additional output. Pass the `-version` flag to display version information.

You can pass any combination of files and directory pathnames to `dupes`. The pathnames of duplicate files are output, one to a line. To do something further with these pathnames, pipe `stdout` to `xargs` or a similar program.

## Examples

By default, `dupes` prints the pathnames of duplicate files to `stdout`.

```bash
dupes /tmp
```

Any errors encountered will be output in alarming red on `stderr`. If you do not wish to see error output, redirect `stderr` to `/dev/null`, e.g.:

```bash
dupes /var 2>/dev/null
```

Output from `dupes` can be piped to `xargs` to perform operations on duplicate files, _e.g._:

```bash
dupes /tmp | xargs ls -l
```

Pass the `-verbose` flag if you wish to see additional information about duplicate matches. The pathname of the duplicate file is output to `stdout`. All additional information is output to `stderr` on the following line. This enables verbose output to remain compatible with `xargs` and similar programs.

Duplicate pathnames containing embedded spaces causes problems for `xargs` and friends. If you had been using `find`, you could have passed it the `-print0` switch and invoked `xargs` with `-0`. Although `dupes` does not have a `-print0` flag, it can be simulated using `sed` to quote the strings:

```bash
dupes /var 2>/dev/null | sed 's/.*/"&"/' | xargs ls -l
```

Alternately, use the simpler form:

```bash
dupes /var 2>/dev/null | xargs -I{} ls -l {}
```

## Notes for Windows users

All of the above examples make use of `xargs`, which is not a native Windows program. However, `xargs` implementations for Windows do exist. Alternately, you can run `dupes.exe` from a Git Bash window and pipe the output to `xargs`.

Building from source is not currently as easy/fun as it is on Linux and macOS. Windows uses `nmake`, and `Makefile` and `config.mk` are written for GNU Make, but don't worry. It is easy to approximate things. Just `type config.mk`, remember the version number, and type:

```bash
go build -ldflags "-s -w -X main.release=VERSION" -o dupes.exe
```

Replace `VERSION` with the version number you saw in `config.mk`, _e.g._, `1.2.0`.

## Credits

- [Ramsey Dow](https://github.com/yesmar)

## License

The BSD 2-clause "Simplified" License. See the [`LICENSE`](https://github.com/yesmar/dupes/blob/master/LICENSE) file for more details.
