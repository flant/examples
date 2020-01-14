Листинги для интеграции Kubernetes Dashboard и пользователей GitLab из [нашей статьи на хабре](https://habr.com/ru/company/flant/blog/452988/).

# Инструкции

Используемый OAuth-провайдер — GitLab.

Регистрируем в GitLab новое приложение: Admin area → Applications → New application.

Устанавливаем Redirect URI (callback url) вида ```https://dashboard.example.com/oauth2/callback```

Далее пользуемся Bash-скриптом:
```
Usage: ctl.sh [OPTION]... --gitlab-url GITLAB_URL --oauth2-id ID --oauth2-secret SECRET --dashboard-url DASHBOARD_URL
Install kubernetes-dashboard to Kubernetes cluster.
Mandatory arguments:
  -i, --install                install into 'kube-system' namespace
  -u, --upgrade                upgrade existing installation, will reuse password and host names
  -d, --delete                 remove everything, including the namespace
      --gitlab-url             set gitlab url with schema (https://gitlab.example.com)
      --oauth2-id              set OAUTH2_PROXY_CLIENT_ID from gitlab
      --oauth2-secret          set OAUTH2_PROXY_CLIENT_SECRET from gitlab
      --dashboard-url          set dashboard url without schema (dashboard.example.com)
Optional arguments:
  -h, --help                   output this message
```

Ссылки на дополнительную документацию:

* https://github.com/colemickens/oauth2_proxy/blob/master/README.md#gitlab-auth-provider
* https://docs.gitlab.com/ce/integration/oauth_provider.html
* https://github.com/kubernetes/ingress/tree/master/examples/external-auth/nginx

