# rifraf

A browser omni search bar enhancer. Side effects may include enhanced productivity, happiness and less google search tracking.


## How to install


- Install it:
```
go get git.vdhsn.com/adam/riffraff
```

- Start it:
```
rifiraff [-p 80] [-accesslog]
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
