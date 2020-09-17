### public IP checker

```
OS req:
- libc-dev 
- libpcap-dev
- gcc
- make
```

```
$ make build
$ sudo SERVER_ADDR=8080 ./whatismyip
OR
$ sudo SERVER_ADDR=8080 ./whatismyip <network_interface>
```