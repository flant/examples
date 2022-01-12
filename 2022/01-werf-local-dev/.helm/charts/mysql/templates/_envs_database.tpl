{{- define "mysql_database_envs" }}
- name: MYSQL_DATABASE
  value: "{{ pluck .Values.global.env .Values.envs.MYSQL_DATABASE | first | default .Values.envs.MYSQL_DATABASE._default }}"
- name: MYSQL_ROOT_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ .Chart.Name }}
      key: MYSQL_ROOT_PASSWORD
{{- end }}
