<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/" xmlns:moz="http://www.mozilla.org/2006/browser/search/">
  <ShortName>Riff Raff</ShortName>
  <Description>slick omnibar search enhancer</Description>
  <InputEncoding>UTF-8</InputEncoding>
  <Image width="16" height="16" type="image/png">{{ .Host }}/favicon.png</Image>
  <Url type="text/html" template="{{ .Host }}/search">
    <Param name="q" value="{searchTerms}"/>
  </Url>
  <Url type="application/x-suggestions+json" template="{{ .Host }}/search"/>
  <moz:SearchForm>{{ .Host }}/search</moz:SearchForm>
</OpenSearchDescription>