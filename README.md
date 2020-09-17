### public IP checker [in progress]

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

$ curl localhost:8080
<ip>
```