#!/bin/sh
#
# set -eo pipefail
# shopt -s nullglob

# config
LINUX_BINARY_DEV_AMD64="https://dl.n2x.io/binaries/dev/latest/linux/amd64/n2x-node"
LINUX_BINARY_STABLE_AMD64="https://dl.n2x.io/binaries/stable/latest/linux/amd64/n2x-node"
LINUX_BINARY_DEV_386="https://dl.n2x.io/binaries/dev/latest/linux/386/n2x-node"
LINUX_BINARY_STABLE_386="https://dl.n2x.io/binaries/stable/latest/linux/386/n2x-node"
LINUX_BINARY_DEV_ARM64="https://dl.n2x.io/binaries/dev/latest/linux/arm64/n2x-node"
LINUX_BINARY_STABLE_ARM64="https://dl.n2x.io/binaries/stable/latest/linux/arm64/n2x-node"
LINUX_BINARY_DEV_ARM="https://dl.n2x.io/binaries/dev/latest/linux/arm/n2x-node"
LINUX_BINARY_STABLE_ARM="https://dl.n2x.io/binaries/stable/latest/linux/arm/n2x-node"
DARWIN_BINARY_DEV_AMD64="https://dl.n2x.io/binaries/dev/latest/darwin/amd64/n2x-node"
DARWIN_BINARY_STABLE_AMD64="https://dl.n2x.io/binaries/stable/latest/darwin/amd64/n2x-node"
DARWIN_BINARY_DEV_ARM64="https://dl.n2x.io/binaries/dev/latest/darwin/arm64/n2x-node"
DARWIN_BINARY_STABLE_ARM64="https://dl.n2x.io/binaries/stable/latest/darwin/arm64/n2x-node"

LINUX_PKG_REPO_DEV="repo.n2x.dev"
LINUX_PKG_REPO_STABLE="repo.n2x.io"

DOCKER=0
DOCKER_IMAGE_DEV="n2xdev/n2x-node:dev"
DOCKER_IMAGE_STABLE="n2xdev/n2x-node:latest"

DEV_MODE=0

# functions
check_root() {
  if [ "$(id -u)" -ne "0" ]; then
    printf >&2 "Root privileges required\n"
    exit 1
  fi
}

parse_args() {
  while [ -n "${1}" ]; do
    case "${1}" in
      --help|-h)
        usage
        exit 0
        ;;
      --dev)
        DEV_MODE=1
        LINUX_PKG_REPO="${LINUX_PKG_REPO_DEV}"
        DOCKER_IMAGE="${DOCKER_IMAGE_DEV}"
        ;;
      --docker)
        if command -v docker >/dev/null; then
          DOCKER=1
        else
          printf >&2 "Docker binary not found\n\n"
          exit 1
        fi
        ;;
      --token)
        if [ -n "${2}" ]; then
          TOKEN="${2}"
          shift 1
        else
          printf >&2 "A token must be specified with the --token option\n\n"
          exit 1
        fi
        ;;
      --port)
        if [ -n "${2}" ]; then
          PORT="${2}"
          shift 1
        else
          printf >&2 "A port must be specified with the --port option\n\n"
          exit 1
        fi
        ;;
      --dns-port)
        if [ -n "${2}" ]; then
          DNS_PORT="${2}"
          shift 1
        else
          printf >&2 "A DNS-port must be specified with the --dns-port option\n\n"
          exit 1
        fi
        ;;
      *)
       printf >&2 "Unrecognized option '%s'\n\n" "${1}"
       usage
       exit 1
       ;;
    esac
    shift 1
  done

  check_args
}

check_args() {
  if [ -z "${LINUX_PKG_REPO}" ]; then
    LINUX_PKG_REPO="${LINUX_PKG_REPO_STABLE}"
  fi

  if [ -z "${DOCKER_IMAGE}" ]; then
    DOCKER_IMAGE="${DOCKER_IMAGE_STABLE}"
  fi

  if [ -z "${TOKEN}" ]; then
    printf >&2 "A token must be specified with the --token option\n\n"
    usage
    exit 1
  fi

  if [ -z "${PORT}" ]; then
    PORT="57775"
  fi

  if [ -z "${DNS_PORT}" ]; then
    DNS_PORT="53535"
  fi
}

usage() {
  cat << EOF
Usage:
  $(basename "$0") --token <auth_token> [options]
  $(basename "$0") -h|--help

Install n2x-node.

Options:
  --token <auth_token>   node authorization token to join the network
  --port <port>          port used for P2P network traffic (default 57775)
  --dns-port <dns-port>  port used for DNS integration (default 53535)
  --docker               install n2x-node as a docker container
  --dev                  dev mode

  -h, --help             display this help

EOF
}

apt_deps_install() {
  echo "Checking dependencies..."

  if ! command -v curl >/dev/null; then
    echo "Installing missing dependency: curl"
    apt-get -y update
    apt-get -y install curl
  fi

  if ! command -v modprobe >/dev/null; then
    echo "Installing missing dependency: kmod"
    apt-get -y update
    apt-get -y install kmod
  fi
}

apt_node_install() {
  echo "deb [trusted=yes] https://${LINUX_PKG_REPO}/apt/ /" | tee /etc/apt/sources.list.d/n2x.list
  apt-get -y update
  apt-get -y install n2x-node
}

apt_node_uninstall() {
  if command -v dpkg >/dev/null; then
    if dpkg -l | grep -q "n2x-node"; then
      echo "Removing n2x-node deb package..."
      dpkg --remove n2x-node
    fi
  fi
}

yum_deps_install() {
  echo "Checking dependencies..."

  if ! command -v curl >/dev/null; then
    echo "Installing missing dependency: curl"
    yum -y install curl
  fi

  if ! command -v modprobe >/dev/null; then
    echo "Installing missing dependency: kmod"
    yum -y install kmod
  fi
}

yum_node_install() {
  cat <<EOF | tee /etc/yum.repos.d/n2x.repo
[n2x]
name=n2x repository
baseurl=https://${LINUX_PKG_REPO}/yum
enabled=1
gpgcheck=0
EOF
  yum -y install n2x-node
}

yum_node_uninstall() {
  if command -v yum >/dev/null; then
    if yum list installed | grep -q "n2x-node"; then
      echo "Removing n2x-node rpm package..."
      yum -y remove n2x-node
    fi
  fi
}

apk_deps_install() {
  echo "Checking dependencies..."

  if ! command -v curl >/dev/null; then
    echo "Installing missing dependency: curl"
    apk update
    apk add curl
  fi

  if ! command -v modprobe >/dev/null; then
    echo "Installing missing dependency: kmod"
    apk update
    apk add kmod
  fi
}

linux_binary_install() {
  mkdir -p /usr/local/bin
  rm -f /usr/local/bin/n2x-node

  if command -v curl >/dev/null; then
    if ! curl -s -o /usr/local/bin/n2x-node "${NODE_BINARY}"; then
      echo "Unable to download binary"
      exit 1
    fi
  else
    echo "Missing required command: curl"
    echo "Please, install required dependencies: curl"
    exit 1
  fi

  chmod 0750 /usr/local/bin/n2x-node

  cat << EOF > /etc/systemd/system/n2x-node.service
[Unit]
Description=n2x-node service
After=network.target

[Service]
Type=simple
# Another Type: forking
#User=
WorkingDirectory=/var/lib/n2x
ExecStart=/usr/local/bin/n2x-node start
Restart=always
# Other restart options: always, on-abort, etc

# The install section is needed to use
# 'systemctl enable' to start on boot
# For a user service that you want to enable
# and start automatically, use 'default.target'
# For system level services, use 'multi-user.target'
[Install]
WantedBy=multi-user.target
EOF

  mkdir -p /var/lib/n2x
  mkdir -p /var/cache/n2x
}

linux_service_stop() {
  if [ -s /etc/systemd/system/n2x-node.service ]; then
    echo "Stopping n2x-node service..."
    systemctl stop n2x-node
    systemctl disable n2x-node
    systemctl daemon-reload
  fi
}

linux_binary_uninstall() {
  if [ -x /usr/local/bin/n2x-node ]; then
    echo "Removing existing n2x-node binary..."
    rm -f /usr/local/bin/n2x-node
    echo "Removing existing n2x-node service configuration..."
    rm -f /etc/systemd/system/n2x-node.service
  fi
}

darwin_binary_install() {
  mkdir -p /opt/n2x/libexec
  rm -f /opt/n2x/libexec/n2x-node

  if command -v curl >/dev/null; then
    if ! curl -s -o /opt/n2x/libexec/n2x-node "${NODE_BINARY}"; then
      echo "Unable to download binary"
      exit 1
    fi
  else
    echo "Missing required command: curl"
    echo "Please, install required dependencies: curl"
    exit 1
  fi

  chmod 0750 /opt/n2x/libexec/n2x-node
  chown root: /opt/n2x/libexec/n2x-node

  mkdir -p /opt/n2x/var/lib
  mkdir -p /opt/n2x/var/cache
  mkdir -p /opt/n2x/var/log
  mkdir -p /opt/n2x/var/tmp
}

get_system_info() {
  SYSARCH="$(uname -m)"

  case "$(uname -s)" in
    Linux)
      SYSTYPE="Linux"

      os_release_file=
      if [ -s "/etc/os-release" ] && [ -r "/etc/os-release" ]; then
        os_release_file="/etc/os-release"
      elif [ -s "/usr/lib/os-release" ] && [ -r "/usr/lib/os-release" ]; then
        os_release_file="/usr/lib/os-release"
      else
        printf >&2 "Unable to find usable OS release info\n"
      fi

      if [ -n "${os_release_file}" ]; then
        # DISTRO="$(awk -F= '/^ID=/ {print $2}' "${os_release_file}" | tr -d \")"
        . "${os_release_file}"
        DISTRO="${ID}"
        echo "Detected Linux distribution: ${DISTRO}"
      else
        DISTRO="unknown"
      fi
      ;;
    Darwin)
      SYSTYPE="Darwin"
      # SYSVERSION="$(sw_vers -buildVersion)"
      ;;
    FreeBSD)
      SYSTYPE="FreeBSD"
      # SYSVERSION="$(uname -K)"
      ;;
    *)
      printf >&2 "Unsupported system type detected\n"
      exit 1
      ;;
  esac
}

pkg_install() {
  echo "Installing n2x-node.."

  case "${SYSTYPE}" in
    Linux)
      case "${SYSARCH}" in
        x86_64)
          if [ "${DEV_MODE}" -eq 1 ]; then
            NODE_BINARY="${LINUX_BINARY_DEV_AMD64}"
          else
            NODE_BINARY="${LINUX_BINARY_STABLE_AMD64}"
          fi
          ;;
        i386|i686)
          if [ "${DEV_MODE}" -eq 1 ]; then
            NODE_BINARY="${LINUX_BINARY_DEV_386}"
          else
            NODE_BINARY="${LINUX_BINARY_STABLE_386}"
          fi
          ;;
        arm64|armv8b|armv8l|aarch64|aarch64_be)
          if [ "${DEV_MODE}" -eq 1 ]; then
            NODE_BINARY="${LINUX_BINARY_DEV_ARM64}"
          else
            NODE_BINARY="${LINUX_BINARY_STABLE_ARM64}"
          fi
          ;;
        arm|armv6l|armv7l)
          if [ "${DEV_MODE}" -eq 1 ]; then
            NODE_BINARY="${LINUX_BINARY_DEV_ARM}"
          else
            NODE_BINARY="${LINUX_BINARY_STABLE_ARM}"
          fi
          ;;
        *)
          printf >&2 "Unsupported system arch detected\n"
          exit 1
          ;;
      esac

      linux_pkg_install
      ;;
    Darwin)
      case "${SYSARCH}" in
        x86_64)
          if [ "${DEV_MODE}" -eq 1 ]; then
            NODE_BINARY="${DARWIN_BINARY_DEV_AMD64}"
          else
            NODE_BINARY="${DARWIN_BINARY_STABLE_AMD64}"
          fi
          ;;
        arm64|armv8b|armv8l|aarch64|aarch64_be)
          if [ "${DEV_MODE}" -eq 1 ]; then
            NODE_BINARY="${DARWIN_BINARY_DEV_ARM64}"
          else
            NODE_BINARY="${DARWIN_BINARY_STABLE_ARM64}"
          fi
          ;;
        *)
          printf >&2 "Unsupported system arch detected\n"
          exit 1
          ;;
      esac

      darwin_pkg_install
      ;;
    FreeBSD)
      # freebsd_pkg_install
      ;;
  esac
}

linux_pkg_install() {
  docker_uninstall
  linux_service_stop
  linux_binary_uninstall

  case "${DISTRO}" in
    debian|ubuntu)
      echo "Detected Linux apt-based distribution, checking deps.."
      apt_deps_install
      apt_node_install
      ;;
    cloudlinux|almalinux|rocky|rhel|centos|fedora)
      echo "Detected Linux yum-based distribution, checking deps.."
      yum_deps_install
      yum_node_install
      ;;
    alpine)
      echo "Detected Linux apk-based distribution, checking deps.."
      apk_deps_install
      linux_binary_install
      ;;
    *)
      if command -v apt-get >/dev/null; then
        echo "Detected Linux apt-based distribution, checking deps.."
        apt_deps_install
        apt_node_install
      elif command -v yum >/dev/null; then
        echo "Detected Linux yum-based distribution, checking deps.."
        yum_deps_install
        yum_node_install
      elif command -v apk >/dev/null; then
        echo "Detected Linux apk-based distribution, checking deps.."
        apk_deps_install
        linux_binary_install
      else
        echo "Unable to detect Linux distribution, trying binary install..."
        linux_binary_install
      fi
      ;;
  esac

  linux_set_config
  linux_setup
}

darwin_pkg_install() {
  darwin_binary_install
  darwin_set_config
  darwin_setup
}

# freebsd_pkg_install() {
#   freebsd_set_config
#   # freebsd_setup
# }

linux_set_config() {
  echo "Setting n2x-node configuration..."

  mkdir -p /etc/n2x

  if [ -s /etc/n2x/n2x-node.yml ]; then
    echo "Saving existing n2x-node configuration..."
    cp /etc/n2x/n2x-node.yml /etc/n2x/n2x-node.yml_old
  fi

  echo "Creating n2x-node configuration..."

  cat << EOF > /etc/n2x/n2x-node.yml
# n2x-node configuration

token: "${TOKEN}"

# loglevel: INFO
port: ${PORT}
# dnsPort: ${DNS_PORT}
EOF

  chmod 0700 /etc/n2x
  chmod 0600 /etc/n2x/n2x-node.yml
}

darwin_set_config() {
  echo "Setting n2x-node configuration..."

  mkdir -p /opt/n2x/etc

  if [ -s /opt/n2x/etc/n2x-node.yml ]; then
    echo "Saving existing n2x-node configuration..."
    cp /opt/n2x/etc/n2x-node.yml /opt/n2x/etc/n2x-node.yml_old
  fi

  echo "Creating n2x-node configuration..."

  cat << EOF > /opt/n2x/etc/n2x-node.yml
# n2x-node configuration

token: "${TOKEN}"

# loglevel: INFO
port: ${PORT}
# dnsPort: ${DNS_PORT}
EOF

  chmod 0700 /opt/n2x/etc
  chmod 0600 /opt/n2x/etc/n2x-node.yml
}

check_tun_kernel_module() {
  if command -v modprobe >/dev/null; then
    if ! modprobe tun; then
      echo "Unable to load tun kernel module"
      exit 1
    fi
  else
    echo "Unable to check tun kernel module: missing modprobe command"
    echo "Please, install required dependencies: kmod"
    exit 1
  fi
}

linux_setup() {
  echo "Starting n2x-node setup on Linux..."

  check_tun_kernel_module

  echo "Enabling systemd n2x-node service..."

  if command -v systemctl >/dev/null; then
    systemctl daemon-reload
    systemctl enable n2x-node
    systemctl restart n2x-node
  else
    echo "Unable to start service: systemctl NOT detected"
    exit 1
  fi
}

darwin_setup() {
  echo "Starting n2x-node setup on Darwin..."

  /opt/n2x/libexec/n2x-node service-install
  launchctl print system/com.n2x.n2x-node
}

# freebsd_setup() {
#   echo "Starting n2x-node setup on FreeBSD..."
# }

docker_install() {
  linux_service_stop
  linux_binary_uninstall
  apt_node_uninstall
  yum_node_uninstall

  linux_set_config
  docker_setup
}

docker_uninstall() {
  if command -v docker >/dev/null; then
    if docker ps -a | grep -q "n2x-node"; then
      docker rm -f n2x-node
    fi
  fi
}

docker_setup() {
  echo "Starting n2x-node as Docker container..."

  mkdir -p /var/lib/n2x
  # mkdir -p /var/cache/n2x

  check_tun_kernel_module

  docker_uninstall

  docker pull "${DOCKER_IMAGE}"

  docker run -d --restart=always \
    --net=host \
    --cap-add=net_admin \
    --device=/dev/net/tun \
    --name n2x-node \
    -v /etc/n2x:/etc/n2x:ro \
    -v /var/lib/n2x:/var/lib/n2x \
    "${DOCKER_IMAGE}" start
}

# main program

check_root
parse_args "$@"

if [ "${DOCKER}" -eq 1 ]; then
  docker_install
else
  get_system_info
  pkg_install
fi

echo "Done."
exit 0
