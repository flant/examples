echo ""
echo "####################################################################################"
echo ""
echo '
Add new key to the file /etc/docker/daemon.json (default location). Create file if missing:

{
   "insecure-registries": ["registry.local.dev:80"]
}

sudo systemctl restart docker

'
