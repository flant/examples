#!/bin/bash

# VARIABLES
. variables

# Install packages
sudo apt install ca-certificates curl gnupg lsb-release netcat
sudo rm -f /usr/share/keyrings/docker-archive-keyring.gpg 2>&1
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io

# Setup docker
sudo usermod -aG docker $USER
sudo systemctl restart docker

# Install kubectl
# https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
sudo curl -L "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" -o /usr/local/bin/kubectl
sudo chown root. /usr/local/bin/kubectl && sudo chmod 755 /usr/local/bin/kubectl

# Install werf
# https://werf.io/installation.html
echo 'export PATH=$HOME/bin:$PATH' >> ~/.bash_profile
export PATH="$HOME/bin:$PATH"
curl -L "https://tuf.trdl.dev/targets/releases/0.3.1/linux-amd64/bin/trdl" -o /tmp/trdl
mkdir -p ~/bin
install /tmp/trdl ~/bin/trdl
trdl add werf https://tuf.werf.io 1 b7ff6bcbe598e072a86d595a3621924c8612c7e6dc6a82e919abe89707d7e3f468e616b5635630680dd1e98fc362ae5051728406700e6274c5ed1ad92bea52a2
echo 'command -v trdl &>/dev/null && source $(trdl use werf 1.2 ea)' >> ~/.bashrc
source $(trdl use werf 1.2 ea)
werf version

# Install minikube
# https://minikube.sigs.k8s.io/docs/start/
# for Ubuntu:
sudo curl -L https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 -o /usr/local/bin/minikube
sudo chown root. /usr/local/bin/minikube && sudo chmod 755 /usr/local/bin/minikube
