go-tvdb-cli
-----------

A simple wrapper for go-tvdb so I can use thetvdb.com from the CLI

![screenshot](https://github.com/yanfali/go-tvdb-cli/raw/test/screenshot.png)

Features
========

 - curses style interface via [termbox-go](https://github.com/nsf/termbox-go)
 - vim and cursor key bindings
 - built-in help via '?' key - primitive
 - faster than the web interface for quick searchs
 - built-in highlighting for episode titles via '/' key
 - displays episode summary

Building
========

go get github.com/yanfali/go-tvdb-cli
$GOPATH/bin/go-tvdb-cli "[program title]"

Examples
========

go-tvdb-cli "The Moomins"
