Client:
 Version:           {{ .Client.Cli.Version }}
 Go version:        {{ .Client.Go.Version }}
 Git commit:        {{ .Client.GitCommit }}
 Built:             {{ .Client.Built }}
 OS/Arch:           {{ .Client.OsArch }}
 Experimental:      true
{{- if .ServerDeployed }}
Server:
 Engine:
  Version:          {{ .Server.Engine.Version }}
 Web:
  Version:          {{ .Server.Web.Version }}
 API:
  Version:          {{ .Server.API.Version }}
{{- end}}
