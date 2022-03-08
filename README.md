# Talmud

Generates unique password depends on domain and login. You need only one master password.

## How it works

For every account you can specify domain (i. e. **google.com**), 
login (i. e. **reindeer@tarandro.ru**), version (**1** by default) and password length.

Using [hmac algorithm](https://ru.wikipedia.org/wiki/HMAC) application can generate unique password
from these values as a message and the master password as a secret key.
