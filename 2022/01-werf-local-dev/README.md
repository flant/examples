NB. Tested with Ubuntu 20.04 only.

```
cd local
./install.sh
./setup-infra.sh
./deploy-app.sh
kubectl -n local get po
```

Visit http://test.application.local/
