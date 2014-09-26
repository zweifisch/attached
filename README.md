# attached

send file as attachment from commandline

some configure in ~/.attachedrc

```toml
account = "username"
password = "password"

from = "username@gmail.com"
signature = "Your Name"

[smtp]

host = "smtp.gmail.com"
port = 25
```

usage

```sh
attached -t alice@gmail.com -m "documents" image.jpg chart.pdf
```

