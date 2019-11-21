# riffraff

[![GoDoc](https://godoc.org/git.vdhsn.com/adam/riffraff?status.svg)](http://godoc.org/git.vdhsn.com/adam/riffraff)
[![Go Report Card](https://goreportcard.com/badge/git.vdhsn.com/adam/riffraff)](https://goreportcard.com/report/git.vdhsn.com/adam/riffraff)
![Docker Pulls](https://img.shields.io/docker/pulls/vdhsn/riffraff?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/adamveld12/riffraff?style=flat-square)
[![Gocover](http://gocover.io/_badge/github.com/adamveld12/riffraff)](http://gocover.io/github.com/adamveld12/riffraff)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/adamveld12/riffraff)
[![Build Status](https://semaphoreci.com/api/v1/adamveld12/riffraff/branches/master/badge.svg)](https://semaphoreci.com/adamveld12/riffraff)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


A browser omni search bar enhancer. Side effects may include enhanced productivity, happiness and less google search tracking.


## How to install

- Install it:
```
go get git.vdhsn.com/adam/riffraff
```

- Start it:
```
rifiraff [-bind 127.0.0.1] [-p 80] [-accesslog] [-data ./data.json]
```

- Point your browser at it:

- Customize it:
   - add a shortcut: `add gitemoji https://www.webfx.com/tools/emoji-cheat-sheet`
   - add a search *note the `%s`*: `add so https://stackoverflow.com/search?q=%s` or `add go https://godoc.org/?q=%s`

- Use it:
    - use stackoverflow shortcut: `so what is the copy and swap idiom`
    - go to a shortcut: `gitemoji`

## FAQ

- Why rifraf?
    - It's Halloween, I finally watched Rocky Horror Picture Show, and I'm following the butler theme for search apps.

- How does this work?
    - This serves a page over HTTP that implements [Opensearch spec](https://developer.mozilla.org/en-US/docs/Web/OpenSearch). 

## License

GPL
