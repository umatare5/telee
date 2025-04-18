# telee

![GitHub Tag](https://img.shields.io/github/v/tag/umatare5/telee?label=Latest%20version)
[![Go Reference](https://pkg.go.dev/badge/umatare5/telee.svg)](https://pkg.go.dev/github.com/umatare5/telee)
[![Go Report Card](https://goreportcard.com/badge/github.com/umatare5/telee?style=flat-square)](https://goreportcard.com/report/github.com/umatare5/telee)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/umatare5/telee/blob/main/LICENSE)

telee [téli] is a CLI works on **TE**rmina**L** to **E**x**E**cute a command on networking device through the user authentication.

It has following advantages compared to use standard telnet and SSH.

- Reduces login and logout operations.

  No longer have to enter username and password every time! 🎉

- Realizes centralized operation from single host.

  Able to get and compare the status, configuration and others easily! 🎉

For those who use many "expect" scripts and "TeraTerm Macro", telee may be a simple alternative.

In additional, the execution performance of telee is 6 to 72 times faster than [napalm](https://napalm.readthedocs.io/en/latest/cli.html)! 🚀

![](https://github.com/umatare5/telee/blob/images/promo.gif)

## Installation

```bash
docker run ghcr.io/umatare5/telee
```

> [!Tip]
> If you prefer using binaries, download them from the [release page](https://github.com/umatare5/telee/releases).
>
> - Supported Platforms: `linux_amd64`, `linux_arm64`, `darwin_amd64` and `darwin_arm64`

## Syntax

```bash
NAME:
   telee - One-line command executor

USAGE:
   telee -H HOSTNAME -C COMMAND [options...]

VERSION:
   1.7.6

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --hostname value, -H value          Set hostname or IP address. [$TELEE_HOSTNAME]
   --port value, -P value              Set port number. (default: 0)
   --timeout value, -t value           Set timeout seconds. (default: 5)
   --command value, -C value           Set a command. [$TELEE_COMMAND]
   --exec-platform value, -x value     Set exec-platform. Refer to README.md what to be set. (default: "ios")
   --enable-mode, -e, --ena, --enable  Raise to privileged EXEC mode. (default: false)
   --redundant-mode, -r, --redundant   Use redundant prompt mode. (default: false)
   --secure-mode, -s, --sec, --secure  Use ssh mode. (default: false)
   --default-privilege-mode, -d        Use default privileged mode assinged by RADIUS attribute. (default: false)
   --username value, -u value          Set username. (default: "admin") [$TELEE_USERNAME]
   --password value, -p value          Set password. (default: "cisco") [$TELEE_PASSWORD]
   --priv-password value, --pp value   Set password to raise to privileged EXEC mode. (default: "enable") [$TELEE_PRIVPASSWORD]
   --help, -h                          show help (default: false)
   --version, -v                       print the version (default: false)
```

## Usage

- Set credentials as environment variable.

```bash
export TELEE_USERNAME=telee
export TELEE_PASSWORD=Teleedev!
export TELEE_PRIVPASSWORD=Teleedev!!
```

- Run the command with hostname.

```console
$ telee --hostname lab-cat29l-02f99-01 --command "show int descr"
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
lab-cat29l-02f99-01>
```

- Be able to grep the stdout.

```console
$ telee --hostname lab-cat29l-02f99-01 --command "show int descr" | grep "Interface\|down"
Interface                      Status         Protocol Description
Vl1                            admin down     down
Gi0/1                          down           down     CLIENT_DEVICE_LONG_DESCR
Gi0/6                          down           down     CLIENT_DEVICE
Gi0/7                          down           down     CLIENT_DEVICE
Gi0/9                          admin down     down
Gi0/10                         admin down     down
```

- Also able to redirect to file.

```console
$ telee --hostname lab-cat29l-02f99-01 --command "show run" --enable > telee.log
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

- When use on other than IOS, need to set `--exec-platform` (`-x`) option.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ telee -H 192.168.0.250 -C "show sysinfo" -x aireos
  show sysinfo

  Manufacturer's Name.............................. Cisco Systems Inc.
  Product Name..................................... Cisco Controller
  Product Version.................................. 8.5.120.0
  Bootloader Version............................... 1.0.20
  Field Recovery Image Version..................... 7.6.101.1
  Firmware Version................................. PIC 19.0

  OUI File Last Update Time........................ Sun Sep 07 10:44:07 IST 2014

  Build Type....................................... DATA + WPS

  System Name...................................... lab-wlc-01f01-01a
  System Location..................................
  System Contact...................................
  System ObjectID.................................. 1.3.6.1.4.1.9.1.1279
  IP Address....................................... 192.168.0.250
  <snip>
  ```

  </p></details>

- When use ASA, need to set `--enable-mode` option. It doesn't support `ter pag 0` in user-level.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ telee -H lab-asa5505-02f01-01 -C "show version" -x asa --enable-mode --pp Pswd1234#
  show version

  Cisco Adaptive Security Appliance Software Version 9.0(4)
  Device Manager Version 7.1(5)100

  Compiled on Wed 04-Dec-13 08:33 by builders
  System image file is "disk0:/asa904-k8.bin"
  Config file at boot was "startup-config"

  lab-asa5505-02f01-01 up 70 days 2 hours

  Hardware:   ASA5505, 512 MB RAM, CPU Geode 500 MHz,
  Internal ATA Compact Flash, 128MB
  BIOS Flash M50FW016 @ 0xfff00000, 2048KB

  Encryption hardware device : Cisco ASA-5505 on-board accelerator (revision 0x0)
                               Boot microcode        : CN1000-MC-BOOT-2.00
                               SSL/IKE microcode     : CNLite-MC-SSLm-PLUS-2.03
  <snip>
  ```

  </p></details>

- When use SSH, need to set `--secure` option.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ telee -H lab-cat29l-02f99-01 -C "show run" --enable --secure
  show run
  Load for five secs: 8%/0%; one minute: 2%; five minutes: 1%
  Time source is NTP, 02:25:22.496 JST Fri May 14 2021

  Building configuration...

  Current configuration : 18716 bytes
  !
  ! Last configuration change at 01:46:41 JST Fri May 14 2021 by raciadev
  !
  version 15.2
  no service pad
  service tcp-keepalives-in
  service timestamps debug datetime msec localtime show-timezone
  service timestamps log datetime msec localtime show-timezone
  service password-encryption
  !
  hostname lab-cat29l-02f99-01
  <snip>
  ```

  </p></details>

- When use RADIUS to raise the privilege, need to set `--default-privilege-mode` option.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ telee -H lab-nx70-02f01-01 -C "show version" -x nxos --default-privilege-mode
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

- When face the timeout, be able to extend the time using `--timeout` option.

  <details><summary><u>Click to show example</u></summary><p>

  ```console
  $ telee -H lab-fs909-02f01-01 -C "show system" -x allied -u manager --timeout 10
  show system
  Switch System Status                     Date 2021-05-09 Time 01:04:54
  Board     Bay      Board Name
  ----------------------------------------------------------------------
  Base      -        FS909M
  ----------------------------------------------------------------------
  Memory -  DRAM : 32768 kB  FLASH : 8192 kB   MAC : 00-1A-EB-93-1C-95
  ----------------------------------------------------------------------
  SysDescription  : CentreCOM FS909M Ver 1.6.14 B02
  SysContact      :
  SysLocation     : LAB
  SysName         : lab-fs909-02f01-01
  SysUpTime       : 1267989237(146days, 18:11:32)
  Release Version : 1.6.14
  Release built   : B02 (Nov 23 2010 at 14:29:56)
  Flash PROM      : Good
  RAM             : Good
  SW chip         : Good
  <snip>
  ```

  </p></details>

## Usecase

- [umatare5/my-infra-network](https://github.com/umatare5/my-infra-network)

## Exec Platform

- telee works for several operating systems. These are called exec-platform.
- The following table shows each exec-platform was verified on which OS version.

### Matrix

| Name (`-x`) | Description              | Enable Mode (`-e`) | Redundant Mode (`-r`) |
| :---------- | :----------------------- | ------------------ | --------------------- |
| aireos      | Cisco AireOS             | Optional           | Not Available         |
| allied      | AlliedTelesis AlliedWare | Not Available      | Not Available         |
| asa         | Cisco ASA Software       | **REQUIRED**       | Optional              |
| foundry     | Brocade IronWare         | Optional           | Not Available         |
| ios         | Cisco IOS, IOS-XE        | Optional           | Not Available         |
| nxos        | Cisco NX-OS              | Optional           | Not Available         |
| srx         | JuniperNetworks JunOS    | Not Available      | Not Available         |
| ssg         | JuniperNetworks ScreenOS | Not Available      | Optional              |
| yamaha      | YAMAHA RT OS             | Optional           | Not Available         |

### Verified On

- "⚠ Not Verified" means "implemented but not checked". I'm waiting your report! 💓

| Name (`-x`)          | Telnet          | SSH (--secure)   | Default PrivMode (`-d`) |
| :------------------- | :-------------- | :--------------- | ----------------------- |
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

## Development

- build

```bash
make build
```

- release

```bash
make release
```
