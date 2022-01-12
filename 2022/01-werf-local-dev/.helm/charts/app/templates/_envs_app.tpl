{{- define "mysql_app_envs" }}
- name: MYSQL_HOST
  value: "{{ pluck .Values.global.env .Values.envs.MYSQL_HOST | first | default .Values.envs.MYSQL_HOST._default }}"
- name: MYSQL_DATABASE
  value: "{{ pluck .Values.global.env .Values.envs.MYSQL_DATABASE | first | default .Values.envs.MYSQL_DATABASE._default }}"
- name: MYSQL_USERNAME
  value: "{{ pluck .Values.global.env .Values.envs.MYSQL_USERNAME | first | default .Values.envs.MYSQL_USERNAME._default }}"
- name: MYSQL_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ .Chart.Name }}
      key: MYSQL_PASSWORD
{{- end }}
