# Twitter Lists
These are my twitter lists:

{{- range .Lists}}
* [{{ .Name }}](#{{ .Name }})
{{- end}}

Instructions on how to [Run your own](#run-your-own) at the bottom.

{{range .Lists}}
## <a href="https://twitter.com/i/lists/{{ .ID }}">{{ .Name }}</a>
<table>
{{range .Members}}<tr><td><a href="https://twitter.com/{{ .ScreenName }}"><img src="{{ .ProfileImage }}"></a></td><td>
<b><a href="https://twitter.com/{{ .ScreenName }}">@{{ .ScreenName }}</a> ({{ .Name }})</b><br />
<ul>
<li>{{ if .LastTweet }}Last Tweet: {{ .LastTweet }}{{else}}<i>Protected</i>{{end}}</li>
<li>{{ .Description }}</li>
</ul>
</td></tr>
{{end}}
</table>
{{- end}}

# Run Your Own
1. Fork this repo
2. Make a twitter app at https://developer.twitter.com. It needs read-only permission
3. Add all for secrets to this repo:
  * `CONSUMER_KEY`
  * `CONSUMER_SECRET`
  * `ACCESS_TOKEN`
  * `ACCESS_SECRET`
4. Update `config.yaml` with your own list names.
5. (optional) Update `README.md.tpl` with your own flare.
