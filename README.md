This is a simple glob pattern library written Go. This library is worth considering for two reasons. First, it includes the `**` pattern to match directories recursively, in addition to the standard wildcard patterns `?`, character classes, and `*`. More importantly, it's matching algorithm works in O(n) time, where n is the length of the pattern. Russ Cox showed in his article [Glob Matching Can Be Simple And Fast Too](https://research.swtch.com/glob) that some pattern matching algorithms take exponential time.

The library handles wildcard patterns as described in the [glob(7) Linux man page](https://man7.org/linux/man-pages/man7/glob.7.html). In summary:

- sequences of literal characters.
- escaped special characters in a sequence (`\?`, `\[`, `\*`, `\\`).
- character classes that define sets of characters to compare against a single character in a path.
- escaped special characters in a character class (`\!`, `\-`, `\]`, `\\`).
- `?`: any single character except for a path separator.
- `*`: zero-or-more characters in a sequence except for a path separator.
- `**`: zero-or-more characters in a sequence, including path separators.

The API is very simple. There are just two functions, `Match(patternString, pathString string) (bool, int)` and `Validate(pattern string) (int, error)`. `Match` accepts a glob pattern and a path. It returns a `bool` indicating if there was a match, and the number of characters in the path that matched. `Match` assumes there are no errors in the pattern. `Validate` can be called first to ensure the pattern has no syntax errors. It scans the pattern and returns (n, nil) on no error, where n is the number of characters in the pattern, or a non-nil error if there's an issue and the number of characters with no errors.

The tools folder contains `profile.sh` which generates and reports coverage data for the unit test. It also build and runs the code in `cmd/main.go`, which is a program that accepts a pattern and a path on the command line, validates the pattern and (if the pattern is valid) reports whether or not the path is matched with it.
