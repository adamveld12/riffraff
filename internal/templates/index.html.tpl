<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8" />
        <link
            rel="search"
            type="application/opensearchdescription+xml"
            title="riff-raff"
            href="/search_plugin.xml"
        />

        <Url type="application/opensearchdescription+xml"
            rel="self"
            template="/search_plugin.xml" />
    </head>
    <body style="display: flex; width: 100%; align-items: center; flex-direction: column;">
        <div style="display: flex; align-content: center; flex-direction: column; min-width: 40%;">
            <h1>Hello internet - {{ .Host }} greets you.</h1>
            <!-- Add search box form here for testing/adding shortcuts -->
            <form action="/search" method="GET" style="display:flex; margin-bottom: 10px; width: 100%;">
                <input type="text" name="q" value="" style="flex-grow: 1;"/>
                <input type="submit" />
            </form>
            <div style="background: lightgray; padding: 10px 5px; margin-bottom: 10px;">
            add &lt;shortcut&gt; &lt;url&gt;: adds a url as a shortcut<br/>
            remove &lt;shortcut&gt; : removes a url shortcut<br/>
            </div>
            <ul style="list-style-type: none; padding: 0; margin-top: 15px;">
                {{ range $key, $value := .Entries }}
                <li>
                   - {{ $key }} → <a href="{{ $value }}">{{ $value }}</a>
                </li>
                {{ end }}
            </ul>
        </div>
    </body>
</html>
