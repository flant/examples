TAGS=$(git ls-remote --refs --tags https://github.com/flant/werf.git "v*" | awk '{print $2}' | cut -d / -f 3 | xargs)

RELEASES=$(for releases in $(echo $TAGS | xargs -n 1 | grep -oEi '^v([0-9]+\.[0-9])' | cut -dv -f2 | sort -nr | uniq) ; do echo "$releases%$(multiwerf available-releases $releases | head -n 1 | awk '{print $1}')"; done | xargs)
echo "export RELEASES='$RELEASES'"

CHANNELS=$(for releases in $(echo $TAGS | xargs -n 1 | grep -oEi '^v([0-9]+\.[0-9])' | cut -dv -f2 | sort -nr | uniq) ; do echo $(multiwerf available-releases $releases -o text | grep -E '^v.+ \[(.+) (alpha|beta|ea|rc)+\]' | sed -E 's#^(v.+) \[([0-9\.]+) (alpha|beta|ea|rc).*\]$#\2-\3\%\1#i'); done | xargs)
echo "export CHANNELS='$CHANNELS'"

ROOT_VERSION=$(multiwerf available-releases ${ROOT_RELEASE} | head -n 1 | awk '{print $1}')
echo "export ROOT_VERSION='$ROOT_VERSION'"
