[![n2x.io](https://github.com/n2x-io/assets/blob/HEAD/images/logo/n2x-logo_black_180x34.png)](https://n2x.io)

[![Discord](https://img.shields.io/discord/654291649572241408?color=%236d82cb&style=flat&logo=discord&logoColor=%23ffffff&label=Chat)](https://n2x.io/discord)
[![GitHub Discussions](https://img.shields.io/badge/GitHub_Discussions-181717?style=flat&logo=github&logoColor=white)](https://github.com/orgs/n2x-io/discussions)
[![X](https://img.shields.io/badge/Follow_on_X-000000?style=flat&logo=x&logoColor=white)](https://x.com/n2xHQ)
[![Mastodon](https://img.shields.io/badge/Follow_on_Mastodon-2f0c7a?style=flat&logo=mastodon&logoColor=white)](https://mastodon.social/@n2x)

Open source projects from [n2x.io](https://n2x.io).

# n2x-node

[![Go Report Card](https://goreportcard.com/badge/n2x.dev/x-node)](https://goreportcard.com/report/n2x.dev/x-node)
[![Release](https://img.shields.io/github/v/release/n2x-io/x-node?display_name=tag&style=flat)](https://github.com/n2x-io/x-node/releases/latest)
[![GitHub](https://img.shields.io/github/license/n2x-io/x-node?style=flat)](/LICENSE)

This repository contains the `n2x-node` agent, the component that runs on the machines you want to connect to your [n2x.io](https://n2x.io) network.

`n2x-node` is available for a variety of Linux platforms, macOS and Windows.

## Minimum Requirements

`n2x-node` has the same [minimum requirements](https://github.com/golang/go/wiki/MinimumRequirements#minimum-requirements) as Go:

- Linux kernel version 2.6.23 or later
- Windows 7 or later
- FreeBSD 11.2 or later
- MacOS 10.11 El Capitan or later

## Getting Started

The instructions in this repo assume you already have a n2x.io account and are ready to start adding nodes.

See [Quick Start](https://n2x.io/docs/platform/getting-started/quickstart/) to learn how to start building your n2x.io cloud-agnostic architecture.

See [Installation](#installation) for more details and other platforms.

## Documentation

For the complete n2x.io platform documentation visit [n2x.io/docs](https://n2x.io/docs/).

## Installation

### Binary Downloads

Linux, macOS and Windows binary downloads are available from the [Releases](https://github.com/n2x-io/x-node/releases) page.

You can download the pre-compiled binaries and install them with the appropriate tools.

### Linux Installation

#### Linux installation with one-line command

The n2x.io platform provides a **one-line command** for the simplest and quickest way to install the `n2x-node` agent in Linux. You'll find this command readily available when you [add a new node](https://n2x.io/docs/howto-guides/nodes/manage-nodes/#add-a-connected-node) through the n2x.io webUI.

Once installed you can review the configuration at `/etc/n2x/n2x-node.yml`.

> See the [n2x-node configuration reference](https://n2x.io/docs/reference/node-configuration/) to find all the configuration options.

#### Linux binary installation with curl

1. Download the latest release.

    ```shell
    curl -LO "https://dl.n2x.io/binaries/stable/latest/linux/amd64/n2x-node"
    ```

2. Validate the binary (optional).

    Download the n2x-node checksum file:

    ```shell
    curl -LO "https://dl.n2x.io/binaries/stable/latest/linux/amd64/n2x-node_checksum.sha256"
    ```

    Validate the n2x-node binary against the checksum file:

    ```bash
    sha256sum --check < n2x-node_checksum.sha256
    ```

    If valid, the output must be:

    ```console
    n2x-node: OK
    ```

    If the check fails, sha256 exits with nonzero status and prints output similar to:

    ```console
    n2x-node: FAILED
    sha256sum: WARNING: 1 computed checksum did NOT match
    ```

3. Install n2x-node and create its configuration file according to your needs.

    ```shell
    sudo install -o root -g root -m 0750 n2x-node /usr/local/bin/n2x-node
    sudo mkdir /var/lib/n2x
    sudo mkdir /var/cache/n2x
    sudo mkdir /etc/n2x
    sudo vim /etc/n2x/n2x-node.yml
    ```

    See the [n2x-node configuration reference](https://n2x.io/docs/reference/node-configuration/) to find all the configuration options.

4. Create the `n2x-node.service` for systemd.

    ```shell
    sudo cat << EOF > /etc/systemd/system/n2x-node.service
    [Unit]
    Description=n2x-node service
    Documentation=https://github.com/n2x-io/x-node
    After=network.target

    [Service]
    Type=simple

    # Another Type: forking

    # User=
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
    ```

5. Ensure the `tun` kernel module is loaded.

    ```shell
    sudo modprobe tun
    ```

6. Start the `n2x-node` service.

    ```shell
    sudo systemctl daemon-reload
    sudo systemctl enable n2x-node
    sudo systemctl restart n2x-node
    ```

#### Package Repository

n2x.io provides a package repository that contains both DEB and RPM downloads.

##### **Debian/Ubuntu**

1. Run the following to setup a new APT `sources.list` entry and install `n2x-node`:

    ```shell
    echo 'deb [trusted=yes] https://repo.n2x.io/apt/ /' | sudo tee /etc/apt/sources.list.d/n2x.list
    sudo apt update
    sudo apt install n2x-node
    ```

2. Check `n2x-node` service status:

    ```shell
    sudo systemctl status n2x-node
    ```

##### **RHEL/CentOS** 

1. Run the following to create a `n2x.repo` file and install `n2x-node`:

    ```shell
    cat <<EOF | sudo tee /etc/yum.repos.d/n2x.repo
    [n2x]
    name=n2x repository - stable
    baseurl=https://repo.n2x.io/yum
    enabled=1
    gpgcheck=0
    EOF
    sudo yum install n2x-node
    ```

2. Check `n2x-node` service status:

    ```shell
    sudo systemctl status n2x-node
    ```

### macOS Installation

#### macOS installation with one-line command

The n2x.io platform provides a **one-line command** for the simplest and quickest way to install the `n2x-node` agent in macOS. You'll find this command readily available when you [add a new node](https://n2x.io/docs/howto-guides/nodes/manage-nodes/#add-a-connected-node) through the n2x.io webUI.

Once installed you can review the configuration at `/etc/n2x/n2x-node.yml`.

> See the [n2x-node configuration reference](https://n2x.io/docs/reference/node-configuration/) to find all the configuration options.

#### macOS binary installation with curl

1. Download the latest release.

    **Intel**:

    ```shell
    curl -LO "https://dl.n2x.io/binaries/stable/latest/darwin/amd64/n2x-node"
    ```

    **Apple Silicon**:

    ```shell
    curl -LO "https://dl.n2x.io/binaries/stable/latest/darwin/arm64/n2x-node"
    ```

2. Validate the binary (optional).

    Download the n2x-node checksum file:

    **Intel**:

    ```shell
    curl -LO "https://dl.n2x.io/binaries/stable/latest/darwin/amd64/n2x-node_checksum.sha256"
    ```

    **Apple Silicon**:

    ```shell
    curl -LO "https://dl.n2x.io/binaries/stable/latest/darwin/arm64/n2x-node_checksum.sha256"
    ```

    Validate the n2x-node binary against the checksum file:

    ```console
    shasum --algorithm 256 --check n2x-node_checksum.sha256
    ```

    If valid, the output must be:

    ```console
    n2x-node: OK
    ```

    If the check fails, sha256 exits with non-zero status and prints output similar to:

    ```console
    n2x-node: FAILED
    sha256sum: WARNING: 1 computed checksum did NOT match
    ```

3. Install `n2x-node` and create its configuration file according to your needs.

    ```console
    chmod +x n2x-node
    sudo mkdir -p /opt/n2x/libexec
    sudo mv n2x-node /opt/n2x/libexec/n2x-node
    sudo chown root: /opt/n2x/libexec/n2x-node
    sudo mkdir -p /opt/n2x/etc
    sudo vim /opt/n2x/etc/n2x-node.yml
    sudo chmod 600 /opt/n2x/etc/n2x-node.yml
    sudo mkdir -p /opt/n2x/var/lib
    sudo mkdir -p /opt/n2x/var/cache
    ```

    > **IMPORTANT**: In macOS, `iface` must be `utun[0-9]+` in the `n2x-node.yml`, being `utun7` usually a good choice for that setting. Use the command `ifconfig -a` before launching the `n2x-node` service and check that the interface is not in-use.

    See the [n2x-node configuration reference](https://n2x.io/docs/reference/node-configuration/) to find all the configuration options.

4. Install and start the n2x-node agent as a system service.

    ```shell
    sudo /opt/n2x/libexec/n2x-node service-install
    ```

5. Check the service status.

    ```shell
    launchctl print system/com.n2x.n2x-node
    ```

    You should get an output like this:

    ```console
    system/com.n2x.n2x-node = {
        active count = 1
        path = /Library/LaunchDaemons/com.n2x.n2x-node.plist
        state = running

        program = /opt/n2x/libexec/n2x-node
        arguments = {
            /opt/n2x/libexec/n2x-node
            service-start
        }

        working directory = /var/tmp

        stdout path = /usr/local/var/log/com.n2x.n2x-node.out.log
        stderr path = /usr/local/var/log/com.n2x.n2x-node.err.log
        default environment = {
            PATH => /usr/bin:/bin:/usr/sbin:/sbin
        }

        environment = {
            XPC_SERVICE_NAME => com.n2x.n2x-node
        }

        domain = system
        minimum runtime = 10
        exit timeout = 5
        runs = 1
        pid = 3925
        immediate reason = speculative
        forks = 28
        execs = 1
        initialized = 1
        trampolined = 1
        started suspended = 0
        proxy started suspended = 0
        last exit code = (never exited)

        spawn type = daemon (3)
        jetsam priority = 4
        jetsam memory limit (active) = (unlimited)
        jetsam memory limit (inactive) = (unlimited)
        jetsamproperties category = daemon
        submitted job. ignore execute allowed
        jetsam thread limit = 32
        cpumon = default

        properties = keepalive | runatload | inferred program
    }
    ```

### Windows Installation

#### Windows installation with one-line command

The n2x.io platform provides a **one-line command** for the simplest and quickest way to install the `n2x-node` agent in Windows. You'll find this command readily available when you [add a new node](https://n2x.io/docs/howto-guides/nodes/manage-nodes/#add-a-connected-node) through the n2x.io webUI.

Once installed you can review the configuration at `/etc/n2x/n2x-node.yml`.

> See the [n2x-node configuration reference](https://n2x.io/docs/reference/node-configuration/) to find all the configuration options.

#### Windows binary installation with curl

1. Open the Command Prompt as Administrator and create a folder for n2x.

    ```shell
    mkdir 'C:\Program Files\n2x'
    ```

2. Download the latest release into the n2x folder.

    ```shell
    curl -LO "https://dl.n2x.io/binaries/stable/latest/windows/amd64/n2x-node.exe"
    ```

3. Validate the binary (optional).

    Download the n2x-node.exe checksum file:

    ```shell
    curl -LO "https://dl.n2x.io/binaries/stable/latest/windows/amd64/n2x-node.exe_checksum.sha256"
    ```

    Validate the n2x-node.exe binary against the checksum file:

    - Using Command Prompt to manually compare CertUtil's output to the checksum file downloaded:

         ```shell
         CertUtil -hashfile n2x-node.exe SHA256
         type n2x-node.exe_checksum.sha256
         ```

    - Using PowerShell to automate the verification using the -eq operator to get a `True` or `False` result:

         ```powershell
         $($(CertUtil -hashfile .\n2x-node.exe SHA256)[1] -replace " ", "") -eq $(type .\n2x-node.exe_checksum.sha256).split(" ")[0]
         ```

4. Download the `wintun` driver from <https://wintun.net>.

5. Unzip the wintun archive and copy the AMD64 binary `wintun.dll` to `C:\Program Files\n2x`.

6. Use an editor to create the n2x-node configuration file `C:\Program Files\n2x\n2x-node.yml`.

    See the [n2x-node configuration reference](https://n2x.io/docs/reference/node-configuration/) to find all the configuration options.

7. Install the `n2x-node` agent as a Windows service.

    >**NOTE** The instructions below assume that the `wintun.dll`, `n2x-node.exe` and `n2x-node.yml` files are stored in `C:\Program Files\n2x`.

    ```shell
    'C:\Program Files\n2x\n2x-node.exe' service-install
    ```

8. Start the `n2x-node` service.

    ```shell
    start-Service n2x-node
    ```

9. Check `n2x-node` service status.

    ```shell
    get-Service n2x-node
    ```

## Running with Docker

You can also run the `n2x-node` agent as a Docker container. See examples below.

Registry:

- `ghcr.io/n2x-io/n2x-node`

### One-line command

The n2x.io platform provides a **one-line command** for the simplest and quickest way to running the `n2x-node` agent with Docker. You'll find this command readily available when you [add a new node](https://n2x.io/docs/howto-guides/nodes/manage-nodes/#add-a-connected-node) through the n2x.io webUI.

Once installed you can review the configuration at `/etc/n2x/n2x-node.yml`.

> See the [n2x-node configuration reference](https://n2x.io/docs/reference/node-configuration/) to find all the configuration options.

### Manual

Example usage:

```shell
docker run -d --restart=always \
  --net=host \
  --cap-add=net_admin \
  --device=/dev/net/tun \
  --name n2x-node \
  -e SCAN_FS=/rootfs-host \
  -v /etc/n2x:/etc/n2x:ro \
  -v /var/lib/n2x:/var/lib/n2x \
  -v /:/rootfs-host:ro \
  ghcr.io/n2x-io/n2x-node:latest start
```

## Artifacts Verification

### Binaries

All artifacts are checksummed and the checksum file is signed with [cosign](https://github.com/sigstore/cosign).

1. Download the files you want and the `checksums.txt`, `checksum.txt.pem` and `checksums.txt.sig` files from the [Releases](https://github.com/n2x-io/x-node/releases) page:

2. Verify the signature:

    ```shell
    cosign verify-blob \
        --cert checksums.txt.pem \
        --signature checksums.txt.sig \
        checksums.txt
    ```

3. If the signature is valid, you can then verify the SHA256 sums match with the downloaded binary:

    ```shell
    sha256sum --ignore-missing -c checksums.txt
    ```

### Docker Images

Our Docker images are signed with [cosign](https://github.com/sigstore/cosign).

Verify the signatures:

```console
COSIGN_EXPERIMENTAL=1 cosign verify ghcr.io/n2x-io/n2x-node
```

## Configuration

Once installed you can review the configuration at `/etc/n2x/n2x-node.yml`.

See the [n2x-node configuration reference](https://n2x.io/docs/reference/node-configuration/) to find all the configuration options.

## Uninstall

### Uninstall Linux n2x-node agent

To remove `n2x-node` from the system, use the following commands:

#### Binary

```shell
sudo systemctl stop n2x-node
sudo systemctl disable n2x-node
sudo rm /etc/systemd/system/n2x-node.service
sudo systemctl daemon-reload
sudo rm /usr/local/bin/n2x-node
sudo rm /etc/n2x/n2x-node.yml
sudo rmdir /etc/n2x
sudo rm -rf /var/lib/n2x
sudo rm -rf /var/cache/n2x
```

#### Package Repository

##### **Debian/Ubuntu**

```shell
sudo systemctl stop n2x-node
sudo apt-get -y remove n2x-node
sudo rm /etc/n2x/n2x-node.yml
sudo rmdir /etc/n2x
sudo rm -rf /var/lib/n2x
sudo rm -rf /var/cache/n2x
```

##### **RHEL/Centos**

```shell
sudo systemctl stop n2x-node
sudo yum -y remove n2x-node
sudo rm /etc/n2x/n2x-node.yml
sudo rmdir /etc/n2x
sudo rm -rf /var/lib/n2x
sudo rm -rf /var/cache/n2x
```

### Uninstall macOS n2x-node agent

To remove `n2x-node` from the system, use the following commands:

```shell
sudo /opt/n2x/libexec/n2x-node service-uninstall
sudo rm /opt/n2x/libexec/n2x-node
sudo rm /opt/n2x/etc/n2x-node.yml
sudo rm -rf /opt/n2x
```

### Uninstall Windows n2x-node agent

To remove `n2x-node` from the system, open the Command Prompt as Administrator and use the following commands:

```shell
stop-Service "n2x-node"
'C:\Program Files\n2xn2x-node.exe` service-uninstall
rm 'C:\Program Files\n2x' -r -force
```

## Community

Have questions, need support and or just want to talk about n2x?

Get in touch with the n2x community!

[![Discord](https://img.shields.io/badge/Join_us_on_Discord-5865F2?style=flat&logo=discord&logoColor=white)](https://n2x.io/discord)
[![GitHub Discussions](https://img.shields.io/badge/GitHub_Discussions-181717?style=flat&logo=github&logoColor=white)](https://github.com/orgs/n2x-io/discussions)
[![X](https://img.shields.io/badge/Follow_on_X-000000?style=flat&logo=x&logoColor=white)](https://x.com/n2xHQ)
[![Mastodon](https://img.shields.io/badge/Follow_on_Mastodon-2f0c7a?style=flat&logo=mastodon&logoColor=white)](https://mastodon.social/@n2x)

## Code of Conduct

Participation in the n2x community is governed by the Contributor Covenant [Code of Conduct](https://github.com/n2x-io/.github/blob/HEAD/CODE_OF_CONDUCT.md). Please make sure to read and observe this document.

Please make sure to read and observe this document. By participating, you are expected to uphold this code.

## License

The n2x open source projects are licensed under the [Apache 2.0 License](/LICENSE).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fn2x-io%2Fx-node.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fn2x-io%2Fx-node?ref=badge_large)
