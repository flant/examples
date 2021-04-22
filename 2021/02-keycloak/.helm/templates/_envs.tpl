{{- define "envs" }}
- name: POD_NAME
  valueFrom:
    fieldRef:
      fieldPath: metadata.name
- name: POD_NAMESPACE
  valueFrom:
    fieldRef:
      fieldPath: metadata.namespace
- name: KEYCLOAK_USER
  value: "admin"
- name: KEYCLOAK_PASSWORD
  value: "{{ pluck .Values.global.env .Values.keycloak.password | first | default .Values.keycloak.password._default }}"
- name: DB_VENDOR
  value: "postgres"
- name: DB_USER
  value: "keycloak"
- name: DB_ADDR
  value: "{{ pluck .Values.global.env .Values.db.host | first | default (printf .Values.db.host._default .Values.global.env .Values.global.env) }}"
- name: DB_PORT
  value: "26257"
- name: DB_DATABASE
  value: "keycloak"
- name: PROXY_ADDRESS_FORWARDING
  value: "true"
- name: JDBC_PARAMS
  value: "sslmode=verify-ca&sslcert=/cockroach-certs/client.keycloak.crt&sslkey=/cockroach-certs/client.keycloak.pk8&sslrootcert=/cockroach-certs/ca.crt"
- name: JDBC_PARAMS_IS
  value: "?sslmode=verify-ca&sslcert=/cockroach-certs/client.keycloak.crt&sslkey=/cockroach-certs/client.keycloak.pk8&sslrootcert=/cockroach-certs/ca.crt"
- name: INFINISPAN_SERVER
  value: "{{ pluck .Values.global.env .Values.infinispan.host | first | default (printf .Values.infinispan.host._default .Values.global.env) }}"
- name: JAVA_OPTS
  value: "{{ pluck .Values.global.env .Values.java | first | default .Values.java._default }}"
{{- end }}
