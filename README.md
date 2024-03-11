# AiSEG HomeKit

A very fragile AiSEG HomeKit bridge. All functions are based HTML parsing.

## Usage

The IP address of AiSEG can be discovered automatically. However, the username and password are required. Usually the username is `aiseg` and the password is the HEMS device ID without hyphen.

The username and password can be set by environment variables like below:

```
export AISEG_USER="aiseg"
export AISEG_PASSWORD="deviceid"
```

The PIN code for HomeKit is `20030010` by default. It can be changed like below:

```
export AISEG_PIN=xxxxxxxx
```
## Disclaimer

**Use it at your own risk**
