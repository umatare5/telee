# telee

telee [téli] is a **TEL**net client to **E**x**E**cute a command on networking device that equip user-authentication.

It has the following advantages compared to use standard telnet.

- Reduces login and logout operations.

  No longer have to enter username and password every time!🎉

- Realize centralized operation from single host.

  Can get and compare a lot of status, configuration and others easily!🎉

If you use a lot of "expect" scripts and "TeraTerm Macro", telee may be a simple alternative.

## Installation

Download from [release page](https://github.com/umatare5/telee/releases).

telee works on `linux_amd64`, `linux_arm64`, `darwin_amd64` and `darwin_arm64`.

## Syntax

```bash
NAME:
   telee - One-line telnet client

USAGE:
   telee -H HOSTNAME -C COMMAND [options...]

VERSION:
   1.2.2

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --hostname value, -H value          Set hostname or IP address. [$TELEE_HOSTNAME]
   --port value, -P value              Set port number. (default: 23)
   --timeout value, -t value           Set timeout seconds. (default: 5)
   --command value, -C value           Set a command. [$TELEE_COMMAND]
   --exec-platform value, -x value     Set exec-platform. Refer to README.md what to be set. (default: "ios")
   --enable-mode, -e, --ena, --enable  Raise to privileged EXEC mode. (default: false)
   --username value, -u value          Set username. (default: "admin") [$TELEE_USERNAME]
   --password value, -p value          Set password. (default: "cisco") [$TELEE_PASSWORD]
   --priv-password value, --pp value   Set password to raise to privileged EXEC mode. (default: "enable") [$TELEE_PRIVPASSWORD]
   --help, -h                          show help (default: false)
   --version, -v                       print the version (default: false)
```

## Usage

- Set credentials into environment variables.

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

- Also be able to redirect to file.

```console
$ telee --hostname lab-cat29l-02f99-01 --command "show int descr" > telee.log
$ head -n 10 telee.log
show int descr
Load for five secs: 1%/0%; one minute: 1%; five minutes: 1%
Time source is NTP, 23:20:51.295 JST Sat May 8 2021

Interface                      Status         Protocol Description
Vl1                            admin down     down
Vl800                          up             up       *** LAB-MGMT ***
Gi0/1                          down           down     CLIENT_DEVICE_LONG_DESCR
Gi0/2                          up             up       CLIENT_DEVICE
Gi0/3                          up             up       CLIENT_DEVICE
```

- When use on other than IOS, need to set `-x` (`--exec-platform`) option.

  <details><summary><code>
  telee -H 192.168.0.250 -C "show sysinfo" -x aireos
  </code></summary><p>

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

- When use slow response devices, recommend to extend `--timeout`.

  <details><summary><code>
  telee -H lab-fs909-02f-01 -C "show system" -x allied --username manager --timeout 10
  </code></summary><p>

  ```console
  $ telee -H lab-fs909-02f-01 -C "show system" -x allied -u manager -t 10
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
  SysName         : lab-fs909-02f-01
  SysUpTime       : 1267989237(146days, 18:11:32)
  Release Version : 1.6.14
  Release built   : B02 (Nov 23 2010 at 14:29:56)
  Flash PROM      : Good
  RAM             : Good
  SW chip         : Good
  <snip>
  ```

  </p></details>

- When use ASA, must set `--enable-mode`. ASA doesn't support "terminal pager 0" in user-level.

  <details><summary><code>
  telee -H lab-asa5505-02f-01 -C "show version" -x asa --enable-mode --priv-password Pswd1234#
  </code></summary><p>

  ```console
  $ telee -H lab-asa5505-02f-01 -C "show version" -x asa -e --pp Pswd1234#
  show version

  Cisco Adaptive Security Appliance Software Version 9.0(4)
  Device Manager Version 7.1(5)100

  Compiled on Wed 04-Dec-13 08:33 by builders
  System image file is "disk0:/asa904-k8.bin"
  Config file at boot was "startup-config"

  lab-asa5505-02f-01 up 70 days 2 hours

  Hardware:   ASA5505, 512 MB RAM, CPU Geode 500 MHz,
  Internal ATA Compact Flash, 128MB
  BIOS Flash M50FW016 @ 0xfff00000, 2048KB

  Encryption hardware device : Cisco ASA-5505 on-board accelerator (revision 0x0)
                               Boot microcode        : CN1000-MC-BOOT-2.00
                               SSL/IKE microcode     : CNLite-MC-SSLm-PLUS-2.03
  <snip>
  ```

  </p></details>

## Exec Platform

- telee works for several operating systems. These are called exec-platform.
- The following table shows each exec-platform was verified on which OS version.

  | Name (`-x`) | Description                   | Verified On     |
  | :---------- | :---------------------------- | :-------------- |
  | aireos      | Cisco AireOS                  | ✅ 8.5.120.0    |
  | allied      | AlliedTelesis AlliedWare      | ✅ 1.6.14B02    |
  | asa         | Cisco ASA Software            | ✅ 9.0(4)       |
  | asa-ha      | Cisco ASA Software (HA)       | ✅ 9.10(1)      |
  | foundry     | Brocade IronWare              | ✅ 07.2.02aT7e1 |
  | ios         | Cisco IOS                     | ✅ 15.2(5c)E    |
  | ssg         | JuniperNetworks ScreenOS      | ✅ 6.3.0r21.0   |
  | ssg-ha      | JuniperNetworks ScreenOS (HA) | ✅ 6.3.0r22.0   |

## Development

- build

```bash
make build
```

- release

```bash
make release
```
