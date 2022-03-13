# Talmud

Generates unique password depends on domain and login. You need only one master password.

## How it works

For every account you can specify domain (i. e. **google.com**), 
login (i. e. **reindeer@tarandro.ru**), version (**1** by default) and password length.

Using [hmac algorithm](https://ru.wikipedia.org/wiki/HMAC) application can generate unique password
from these values as a message and the master password as a secret key.
 
## Compilation

### MacOs

```shell
go build -ldflags "-w -s" -o bin/talmud main.go
```

### Linux with dbus SecretService support

If application is used with DBUS SecretService spec like gnome-keyring or ksecretservice.

```shell
go build -ldflags "-w -s" -tags dbus -o bin/talmud main.go
```

### Linux without dbus support

If application is used without DBUS support and reads master password from stdin.

```shell
go build -ldflags "-w -s" -o bin/talmud main.go
```
