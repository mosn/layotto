# check if you already have an etcd
if test -e etcd; then
  exit 0
elif [ $(which etcd | wc -l) -gt 0 ]; then
  cp $(which etcd) ./etcd
  exit 0
fi

# configuration
ETCD_VER=v3.4.18

# choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
GITHUB_URL=https://github.com/etcd-io/etcd/releases/download
DOWNLOAD_URL=${GITHUB_URL}

download_etcd_linux() {
  rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
  rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test

  curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
  tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1
  rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz

  mv /tmp/etcd-download-test/etcd etcd
  mv /tmp/etcd-download-test/etcdctl etcdctl

  ./etcd --version
  ./etcdctl version
}

download_etcd_mac() {
  rm -f /tmp/etcd-${ETCD_VER}-darwin-amd64.zip
  rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test

  curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-darwin-amd64.zip -o /tmp/etcd-${ETCD_VER}-darwin-amd64.zip
  unzip /tmp/etcd-${ETCD_VER}-darwin-amd64.zip -d /tmp && rm -f /tmp/etcd-${ETCD_VER}-darwin-amd64.zip
  mv /tmp/etcd-${ETCD_VER}-darwin-amd64/* /tmp/etcd-download-test && rm -rf mv /tmp/etcd-${ETCD_VER}-darwin-amd64

  mv /tmp/etcd-download-test/etcd etcd
  mv /tmp/etcd-download-test/etcdctl etcdctl

  ./etcd --version
  ./etcdctl version
}

# download etcd
if test "$(uname)" = "Darwin"; then
  # Mac OS X
  download_etcd_mac
elif test "$(expr substr $(uname -s) 1 5)" = "Linux"; then
  # GNU/Linux
  download_etcd_linux
else
  # Windows or other OS
  echo "Your OS is not supported!"
  exit 1
fi
