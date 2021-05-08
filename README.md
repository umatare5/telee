# telee

telee [téli] is a **TEL**net client to **E**x**E**cute a command on remote device via authentication.

You no longer need to log-in to run `show int status`, `show int description` and others!

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
   1.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --hostname value, -H value         Set hostname or IP address. [$TELEE_HOSTNAME]
   --port value, -P value             Set port number. (default: 23)
   --timeout value, -t value          Set timeout seconds. (default: 5)
   --command value, -C value          Set a command. [$TELEE_COMMAND]
   --platform value, -x value         Set a platform. It supports 'ios' and 'aireos'. (default: "ios")
   --enable-mode, -e                  Log in to privileged EXEC mode. (default: false)
   --username value, -u value         Set username. (default: "admin") [$TELEE_USERNAME]
   --password value, -p value         Set password. (default: "cisco") [$TELEE_PASSWORD]
   --priv-password value, --pp value  Set password to change to privileged EXEC mode. (default: "enable") [$TELEE_PRIVPASSWORD]
   --help, -h                         show help (default: false)
   --version, -v                      print the version (default: false)
```

## Usage

- Set environment variables.

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

- Use `-x aireos` to execute the command for AireOS.

```console
$ telee --hostname 192.168.0.250 --command "show sysinfo" -x aireos
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
IPv6 Address..................................... ::
Last Reset....................................... Power on reset
System Up Time................................... 0 days 0 hrs 27 mins 16 secs
System Timezone Location......................... (GMT +9:00) Tokyo, Osaka, Sapporo
System Stats Realtime Interval................... 5
System Stats Normal Interval..................... 180

Configured Country............................... J4  - Japan 4(Q)
Operating Environment............................ Commercial (0 to 40 C)
Internal Temp Alarm Limits....................... 0 to 65 C
Internal Temperature............................. +37 C
External Temperature............................. +46 C
Fan Status....................................... 5200 rpm

State of 802.11b Network......................... Enabled
State of 802.11a Network......................... Disabled
Number of WLANs.................................. 6
Number of Active Clients......................... 0

OUI Classification Failure Count................. 0

Burned-in MAC Address............................ E0:AC:F1:E1:BB:20
Maximum number of APs supported.................. 75
System Nas-Id.................................... lab-wlc-01f01-01a
WLC MIC Certificate Types........................ SHA1/SHA2

(Cisco Controller) >
```

## Development

- build

```bash
make build
```

- release

```bash
make release
```
