{{ define "nodeselector" }}
{{ if eq .Values.global.env "production" }}
nodeSelector:
  node-role/production: ""
{{ else }}
{{ end }}
{{ end }}
