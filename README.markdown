# dupes

`dupes` identifies duplicate files within the user-specified set of target files and directories. It does this by walking the specified file systems, computing and comparing SHA-256 hashes of the files it encounters. Any hash that has been seen previously is considered a duplicate.

Important considerations:

- Zero-length files serve as a trigger mechanism on many operating systems. Removing such files could have negative effects, so `dupes` completely ignores empty files.
- Symbolic links *will* be followed, which could lead to trouble should a cycle be encountered. `dupes` assumes that you know what you are doing and understand the file systems you intend to run it on.

The author uses `dupes` to identify and delete duplicate images downloaded from [imgur](https://imgur.com).

## Install

```bash
go get github.com/yesmar/dupes
```

Note that the `vendor` directory is here for stability. Remove it if you already have the dependencies in your `GOPATH`.

## Usage

```bash
Usage: dupes [options] <pathname> […]
```

| Flag | Description |
| :--- | :--- |
| `-no-color` | Disable color output |
| `-verbose` | Verbose output |
| `-version` | Display version information |

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

Output from `dupes` can be piped to `xargs` to perform operations on duplicate files, e.g.:

```bash
dupes /tmp | xargs ls -l
```

Pass the `-verbose` flag if you wish to see additional information. Verbose output is incompatible with `xargs` so it is output to `stderr` as a safeguard.

Duplicate pathnames containing embedded spaces causes problems for `xargs` and friends. If you had been using `find`, you could have passed it the `-print0` switch and invoked `xargs` with `-0`. Although `dupes` does not have a `-print0` option, it can be simulated using `sed`:

```bash
dupes /var 2>/dev/null | sed 's/.*/"&"/' | xargs ls -l
```

## Credits

- [Ramsey Dow](https://github.com/yesmar)

## License

The BSD 2-clause "Simplified" License. See the [`LICENSE`](https://github.com/yesmar/dupes/LICENSE) file for more details.