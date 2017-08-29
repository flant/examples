# kubernetes-dashboard
Используемый oauth провайдер - gitlab

Регистрируем в гитлаб новое приложение. Для этого идем Admin area -> Applications -> New application

Redirect URI(Callback url) устанавливаем вида https://dashboard.example.com/oauth2/callback

```
Usage: ctl.sh [OPTION]... --gitlab-url GITLAB_URL --oauth2-id ID --oauth2-secret SECRET --dashboard-url DASHBOARD_URL
Install kubernetes-dashboard to Kubernetes cluster.
Mandatory arguments:
  -i, --install                install into 'kube-nginx-ingress' namespace
  -u, --upgrade                upgrade existing installation, will reuse password and host names
  -d, --delete                 remove everything, including the namespace
Optional arguments:
  -h, --help                   output this message
      --gitlab-url             set gitlab url with schema (https://gitlab.example.com)
      --oauth2-id              set OAUTH2_PROXY_CLIENT_ID from gitlab
      --oauth2-secret          set OAUTH2_PROXY_CLIENT_SECRET from gitlab
      --dashboard-url          set dashboard url without schema (dashboard.example.com)
```
ссылки на документацию:

* https://github.com/colemickens/oauth2_proxy/blob/master/README.md#gitlab-auth-provider
* https://docs.gitlab.com/ce/integration/oauth_provider.html
* https://github.com/kubernetes/ingress/tree/master/examples/external-auth/nginx
