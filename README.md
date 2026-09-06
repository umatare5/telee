<div align="center">

  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/umatare5/telee/main/docs/assets/logo_dark.png" width="115px" />
    <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/umatare5/telee/main/docs/assets/logo.png" width="115px" />
    <img alt="telee" src="https://raw.githubusercontent.com/umatare5/telee/main/docs/assets/logo.png" width="115px" />
  </picture>

  <h1>telee</h1>

  <p>A command-line interface that logs in to one network device and runs one command on it.</p>

  <p>
    <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/umatare5/telee?label=Latest%20version" />
    <a href="https://github.com/umatare5/telee/actions/workflows/go-test-build.yml"><img alt="Test and Build" src="https://github.com/umatare5/telee/actions/workflows/go-test-build.yml/badge.svg?branch=main" /></a>
    <img alt="Test Coverage" src="https://raw.githubusercontent.com/umatare5/telee/main/docs/assets/coverage.svg" /><br/>
    <a href="https://www.bestpractices.dev/projects/10968"><img alt="OpenSSF Best Practices" src="https://www.bestpractices.dev/projects/10968/badge" /></a>
    <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" /></a>
    <a href="https://developer.cisco.com/codeexchange/github/repo/umatare5/telee"><img alt="Published" src="https://static.production.devnetcloud.com/codeexchange/assets/images/devnet-published.svg" /></a>
  </p>

</div>

## Overview

This CLI opens one telnet or SSH session, logs in, disables paging, runs one command and prints what came back.

- 🔑 **One login**: Credentials arrive from `TELEE_*` variables, so no prompt interrupts a loop over many devices
- 🧭 **Nine platforms**: `-x` picks the prompt, paging and escalation dialect, from Cisco IOS to YAMAHA RT
- 🚿 **Shell-friendly**: Device output is the only thing on stdout, so a pipe or a redirect needs no filtering
- ⚡ **Fast**: 6 to 72 times faster than napalm on one Catalyst 2960L, measured on telee 1.6.5

Where a fleet is driven by `expect` scripts or TeraTerm macros, telee replaces the script with a single invocation. [`docs/measurements.md`](docs/measurements.md) carries the timings behind the figure above, and [umatare5/my-infra-network](https://github.com/umatare5/my-infra-network) is one repository that uses the CLI.

![telee demonstration](https://raw.githubusercontent.com/umatare5/telee/images/promo.gif)

## Supported Environment

telee ships as a static binary and as a `scratch`-based image:

- **Binaries** — `linux_amd64`, `linux_arm64`, `darwin_amd64` and `darwin_arm64`
- **Images** — `ghcr.io/umatare5/telee` for `linux/amd64` and `linux/arm64`, running as UID 65534

The [Exec Platform](#exec-platform) matrix names the OS version each path was verified on.

## Quick Start

### 1. Install the CLI

```bash
docker run --rm ghcr.io/umatare5/telee:latest --help
```

> [!TIP]
> If you prefer using binaries, download them from the [Release](https://github.com/umatare5/telee/releases) page.

### 2. Export the credentials

```bash
export TELEE_USERNAME="operator"
read -rs TELEE_PASSWORD && export TELEE_PASSWORD
```

### 3. Run one command

```bash
telee -H sw01.example.internal -C "show interfaces description"
```

> [!NOTE]
> Every platform but AireOS builds the prompt it waits for out of the value of `-H`, so that value has to match the hostname the device prints rather than merely resolve to its address.

## Syntax

One invocation runs one command on one device, and there are no subcommands.

```bash
telee -H HOSTNAME -C COMMAND [options...]
```

| Flag                             | What it sets                                              |
| :------------------------------- | :-------------------------------------------------------- |
| `--hostname`, `-H`               | Target host, which also builds the prompt telee expects   |
| `--command`, `-C`                | The single command line sent to the device                |
| `--exec-platform`, `-x`          | Platform dialect, `ios` unless set                        |
| `--port`, `-P`                   | TCP port, completed to 22 under `-s` and to 23 otherwise  |
| `--timeout`, `-t`                | Seconds one expect step may wait, 5 unless set            |
| `--secure-mode`, `-s`            | Use SSH in place of telnet                                |
| `--enable-mode`, `-e`            | Send the escalation command and the privileged password   |
| `--default-privilege-mode`, `-d` | Expect a privileged prompt at login and escalate nothing  |
| `--redundant-mode`, `-r`         | Append the failover suffix to every expected prompt       |
| `--username`, `-u`               | Account name, `admin` unless set                          |
| `--password`, `-p`               | Account password, `cisco` unless set                      |
| `--priv-password`, `--pp`        | Privileged password, `enable` unless set                  |
| `--host-key-path`, `--hkp`       | Public key file replacing `~/.ssh/known_hosts` under `-s` |

`telee --help` prints the same flags with their aliases, [`docs/configuration.md`](docs/configuration.md) carries every default and environment variable, and [`docs/README.md`](docs/README.md) indexes the reference pages.

## Usage

- **Default platform** — `ios` over telnet needs nothing but a hostname and a command.

```console
$ telee --hostname sw01 --command "show int descr"
show int descr
Load for five secs: 2%/0%; one minute: 1%; five minutes: 1%
Time source is NTP, 23:16:54.302 JST Sat May 8 2021

Interface                      Status         Protocol Description
Vl1                            admin down     down
Vl800                          up             up       *** LAB-MGMT ***
Gi0/1                          up             up       CLIENT_DEVICE_LONG_DESCR
Gi0/2                          up             up       CLIENT_DEVICE
Gi0/3                          up             up       CLIENT_DEVICE
Gi0/4                          up             up       CLIENT_DEVICE
Gi0/5                          up             up       CLIENT_DEVICE
Gi0/6                          down           down     CLIENT_DEVICE
Gi0/7                          down           down     CLIENT_DEVICE
Gi0/8                          up             up       GATEWAY_ROUTER
Gi0/9                          admin down     down
Gi0/10                         admin down     down
sw01>
```

- **Only device output on stdout** — a pipe sees exactly what the device printed.

```console
$ telee --hostname sw01 --command "show int descr" | grep "Interface\|down"
Interface                      Status         Protocol Description
Vl1                            admin down     down
Gi0/1                          down           down     CLIENT_DEVICE_LONG_DESCR
Gi0/6                          down           down     CLIENT_DEVICE
Gi0/7                          down           down     CLIENT_DEVICE
Gi0/9                          admin down     down
Gi0/10                         admin down     down
```

- **Redirect** — the same bytes land in a file, and `-e` raises the session first once `TELEE_PRIVPASSWORD` is exported.

```console
$ telee --hostname sw01 --command "show run" --enable > telee.log
$ head -n 10 telee.log
show run
Load for five secs: 1%/0%; one minute: 1%; five minutes: 1%
Time source is NTP, 23:21:34.501 JST Sat May 8 2021

Building configuration...

Current configuration : 18687 bytes
!
! Last configuration change at 01:30:16 JST Sun Feb 14 2021
!
```

- **Other platforms** — `-x` selects the dialect for anything that is not IOS.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ telee -H 192.0.2.250 -C "show sysinfo" -x aireos
  show sysinfo

  Manufacturer's Name.............................. Cisco Systems Inc.
  Product Name..................................... Cisco Controller
  Product Version.................................. 8.5.120.0
  Bootloader Version............................... 1.0.20
  Field Recovery Image Version..................... 7.6.101.1
  Firmware Version................................. PIC 19.0

  OUI File Last Update Time........................ Sun Sep 07 10:44:07 IST 2014

  Build Type....................................... DATA + WPS

  System Name...................................... wlc01
  System Location..................................
  System Contact...................................
  System ObjectID.................................. 1.3.6.1.4.1.9.1.1279
  IP Address....................................... 192.0.2.250
  <snip>
  ```

  </p></details>

- **ASA** — `terminal pager 0` is refused from a user-level session, so an `asa` run has to start privileged through either `-e` or `-d`.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ export TELEE_PRIVPASSWORD='<enable password>'
  $ telee -H fw01 -C "show version" -x asa --enable-mode
  show version

  Cisco Adaptive Security Appliance Software Version 9.0(4)
  Device Manager Version 7.1(5)100

  Compiled on Wed 04-Dec-13 08:33 by builders
  System image file is "disk0:/asa904-k8.bin"
  Config file at boot was "startup-config"

  fw01 up 70 days 2 hours

  Hardware:   ASA5505, 512 MB RAM, CPU Geode 500 MHz,
  Internal ATA Compact Flash, 128MB
  BIOS Flash M50FW016 @ 0xfff00000, 2048KB

  Encryption hardware device : Cisco ASA-5505 on-board accelerator (revision 0x0)
                               Boot microcode        : CN1000-MC-BOOT-2.00
                               SSL/IKE microcode     : CNLite-MC-SSLm-PLUS-2.03
  <snip>
  ```

  </p></details>

- **SSH** — `-s` replaces telnet, and its aliases `--sec` and `--secure` name the same flag.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ telee -H sw01 -C "show run" --enable --secure
  show run
  Load for five secs: 8%/0%; one minute: 2%; five minutes: 1%
  Time source is NTP, 02:25:22.496 JST Fri May 14 2021

  Building configuration...

  Current configuration : 18716 bytes
  !
  ! Last configuration change at 01:46:41 JST Fri May 14 2021 by operator
  !
  version 15.2
  no service pad
  service tcp-keepalives-in
  service timestamps debug datetime msec localtime show-timezone
  service timestamps log datetime msec localtime show-timezone
  service password-encryption
  !
  hostname sw01
  <snip>
  ```

  </p></details>

- **RADIUS privilege** — `-d` expects the privileged prompt straight after login, so nothing is escalated and no privileged password is read.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ telee -H sw02 -C "show version" -x nxos --default-privilege-mode
  show version
  Cisco Nexus Operating System (NX-OS) Software
  TAC support: http://www.cisco.com/tac
  Documents: http://www.cisco.com/en/US/products/ps9372/tsd_products_support_series_home.html
  Copyright (c) 2002-2015, Cisco Systems, Inc. All rights reserved.
  The copyrights to certain works contained in this software are
  owned by other third parties and used and distributed under
  license. Certain components of this software are licensed under
  the GNU General Public License (GPL) version 2.0 or the GNU
  Lesser General Public License (LGPL) Version 2.1. A copy of each
  such license is available at
  http://www.opensource.org/licenses/gpl-2.0.php and
  http://www.opensource.org/licenses/lgpl-2.1.php

  Software
  BIOS:      version N/A
  kickstart: version 6.2(14)
  system:    version 6.2(14)
  BIOS compile time:
  kickstart image file is: bootflash:///n7000-s1-kickstart.6.2.14.bin
  <snip>
  ```

  </p></details>

- **Slow devices** — `-t` widens the window each expect step waits in.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ telee -H sw03 -C "show system" -x allied -u manager --timeout 10
  show system
  Switch System Status                     Date 2021-05-09 Time 01:04:54
  Board     Bay      Board Name
  ----------------------------------------------------------------------
  Base      -        FS909M
  ----------------------------------------------------------------------
  Memory -  DRAM : 32768 kB  FLASH : 8192 kB   MAC : 00-00-5E-00-53-01
  ----------------------------------------------------------------------
  SysDescription  : CentreCOM FS909M Ver 1.6.14 B02
  SysContact      :
  SysLocation     : LAB
  SysName         : sw03
  SysUpTime       : 1267989237(146days, 18:11:32)
  Release Version : 1.6.14
  Release built   : B02 (Nov 23 2010 at 14:29:56)
  Flash PROM      : Good
  RAM             : Good
  SW chip         : Good
  <snip>
  ```

  </p></details>

## Exec Platform

telee speaks nine platform dialects, and `-x` selects one. Each decides the prompt telee waits for and the command that disables paging, and five of the nine add an escalation path on top.

### Matrix

| Name (`-x`) | Description              | Enable Mode (`-e`)  | Redundant Mode (`-r`) |
| :---------- | :----------------------- | :------------------ | :-------------------- |
| aireos      | Cisco AireOS             | Not Available       | Not Available         |
| allied      | AlliedTelesis AlliedWare | Not Available       | Not Available         |
| asa         | Cisco ASA Software       | Either `-e` or `-d` | Optional              |
| foundry     | Brocade IronWare         | Optional            | Not Available         |
| ios         | Cisco IOS/IOS-XE         | Optional            | Not Available         |
| nxos        | Cisco NX-OS              | Optional            | Not Available         |
| srx         | Juniper JunOS            | Not Available       | Not Available         |
| ssg         | Juniper ScreenOS         | Not Available       | Optional              |
| yamaha      | YAMAHA RT                | Optional            | Not Available         |

Under Enable Mode, "Not Available" means the platform builds no privileged batch, so `-e` is taken, noticed on stderr and then ignored. Under Redundant Mode it means `-r` is refused before the session opens.

`asa` is the one platform whose paging command needs a privileged session, so a run setting neither `-e` nor `-d` is refused before the dial

`-r` appends the suffix a redundant pair prints — `/pri/act` on ASA, `(M)` on ScreenOS — to every prompt telee expects, which is why the two platforms that accept it are the two that print one.

### Verified On

Each version below is the OS that path was exercised against. "⚠ Not Verified" marks a path that is implemented and reachable but never run on hardware.

| Name (`-x`)          | Telnet          | SSH (`-s`)       | Default PrivMode (`-d`) |
| :------------------- | :-------------- | :--------------- | :---------------------- |
| aireos               | ✅ 8.5.120.0    | ✅ 8.5.120.0     | Not Supported           |
| allied               | ✅ 1.6.14B02    | Not Supported    | Not Supported           |
| asa                  | ✅ 9.0(4)       | ⚠ Not Verified   | ⚠ Not Verified          |
| asa (redundant-mode) | ✅ 9.10(1)      | ⚠ Not Verified   | ⚠ Not Verified          |
| foundry              | ✅ 07.2.02aT7e1 | Not Supported    | Not Supported           |
| ios                  | ✅ 15.2(5c)E    | ✅ 15.2(5c)E     | ✅ 15.2(5c)E            |
| nxos                 | ✅ 6.2(14)      | ⚠ Not Verified   | ✅ 6.2(14)              |
| srx                  | Not Supported   | ✅ 15.1X49-D90.7 | Not Supported           |
| ssg                  | ✅ 6.3.0r21.0   | ⚠ Not Verified   | Not Supported           |
| ssg (redundant-mode) | ✅ 6.3.0r22.0   | ⚠ Not Verified   | Not Supported           |
| yamaha               | ✅ Rev.8.03.94  | ✅ Rev.10.01.78  | Not Supported           |

## Configuration

Six environment variables reach the flags below, and nothing else in the environment is read.

| Variable             | Flag                       | Default  |
| :------------------- | :------------------------- | :------- |
| `TELEE_HOSTNAME`     | `--hostname`, `-H`         | —        |
| `TELEE_COMMAND`      | `--command`, `-C`          | —        |
| `TELEE_USERNAME`     | `--username`, `-u`         | `admin`  |
| `TELEE_PASSWORD`     | `--password`, `-p`         | `cisco`  |
| `TELEE_PRIVPASSWORD` | `--priv-password`, `--pp`  | `enable` |
| `TELEE_HOSTKEYPATH`  | `--host-key-path`, `--hkp` | —        |

Under `-s` the host key is checked against `~/.ssh/known_hosts`, and `--host-key-path` narrows that to one public key file. No flag disables the check. [`docs/configuration.md`](docs/configuration.md) carries the precedence between a flag and its variable, and which flags each platform accepts.

> [!IMPORTANT]
> The guard behind `-e` compares `--priv-password` against its default literal rather than testing it for emptiness, so an unchanged privileged password reads as unset and the run stops.

## Troubleshooting

A rejected argument prints one `ERROR failed to validate arguments error="…"` line and exits 1, and a failed session prints the transport error followed by a `[Hint]` block. Both go to stderr, so a redirect that would have captured device output holds nothing after a failure.

[`docs/troubleshooting.md`](docs/troubleshooting.md) maps each message to the condition that produced it.

## Contributing

[`CONTRIBUTING.md`](CONTRIBUTING.md) carries the `make` targets, the container build and the release process.

## License

MIT. The binary statically links MIT, BSD 3-Clause and Apache 2.0 dependencies, whose notices are reproduced in [`NOTICE`](NOTICE) and shipped alongside [`LICENSE`](LICENSE) in every release archive and container image.
