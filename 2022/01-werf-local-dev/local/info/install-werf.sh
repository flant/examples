#!/usr/bin/env bash
echo ""
echo "#############################################################################################################"
echo ""
echo "Installation manual: https://werf.io/installation.html?version=1.2&channel=ea&os=linux&method=trdl&arch=amd64"
echo '
# Prepare user
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker

# Add ~/bin to the PATH.
echo "export PATH=$HOME/bin:$PATH" >> ~/.bash_profile
export PATH="$HOME/bin:$PATH"

# Install trdl.
curl -L "https://tuf.trdl.dev/targets/releases/0.1.3/linux-amd64/bin/trdl" -o /tmp/trdl
mkdir -p ~/bin
install /tmp/trdl ~/bin/trdl

# Add werf repo
trdl add werf https://tuf.werf.io 1 b7ff6bcbe598e072a86d595a3621924c8612c7e6dc6a82e919abe89707d7e3f468e616b5635630680dd1e98fc362ae5051728406700e6274c5ed1ad92bea52a2

# Enable werf
source $(trdl use werf 1.2 ea)
'

printf "# Auto-enable werf on session start
echo source \$(trdl use werf 1.2 ea) >> ~/.bashrc

# Check evverything works
werf version
"
