## **IP checker**

Check and display the IP address of the host that sent the packet.

### # OS req
```
libc-dev 
libpcap-dev
gcc
make
```

### # local deployment
```
make build

sudo SERVER_PORT=8080 ./whatismyip
OR
sudo SERVER_PORT=8080 ./whatismyip -inter <network_interface>

curl localhost:8080/
<source_ip>
```

### # remote deployment (docker)

For remote deployment use VM, Azure ACI etc. App doesn't work properly if you run it on your workstation (in docker).

```
[vm]$ make docker-run
[PC]$ curl <vm_public_ip>:8080
<source_ip>
```