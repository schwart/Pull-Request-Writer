# Pull Request Writer

Creates a short prompt that you can paste into an LLM and have it write a PR for you.

## Installing

- Make sure you have [go](https://go.dev/doc/install) installed.
- Clone it: `git clone git@github.com:schwart/Pull-Request-Writer.git`
- `cd` into the repo.
- Run `go install` to have it build and install the binary.
- Optional: set `GOBIN` to change the install location. 
- Check `go env GOBIN` to see what the value is (if you haven't changed it). The default is `~/go/bin`.
- Make sure it's on your `$PATH` or you won't be able to run it from anywhere on your machine!

## Running

Run `pr-writer`, if you're currently inside a git repo, you'll see a form appear:

![Screenshot of pr-writer](https://github.com/schwart/Pull-Request-Writer/blob/master/images/screenshot.png?raw=true)

### Default values:

Source branch: the branch that was most recently commited to.
Target branch: either "master", "main" or blank if neither of them can be found in the list of branches.

### Suggestions

Both branch inputs support auto-complete. Start typing and you'll see a match appear, press `ctrl-e` to fill it.
