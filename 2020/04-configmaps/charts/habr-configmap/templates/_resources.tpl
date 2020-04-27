# vi:syntax=yaml
# vi:filetype=yaml

{{- define "resources_app" }}
resources:
  requests:
    memory: {{ .Values.resources.app.requests.memory | quote }}
    cpu: {{ .Values.resources.app.requests.cpu | quote }}
  limits:
    memory: {{ .Values.resources.app.limits.memory | quote }}
{{- end }}
