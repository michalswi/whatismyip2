## **IP checker**

Check and display the IP address of the host that sent the packet.

### # OS req
```
libc-dev 
libpcap-dev
gcc
make
```

### # how to run

For tests use Virtual Machine, Azure VM etc. App doesn't display **Remote IP** neither at your **workstation** nor in docker.

```
ssh <vm>

make build

sudo SERVER_PORT=8080 ./whatismyip
OR
sudo SERVER_PORT=8080 ./whatismyip -inter <network_interface>

> html
$ curl <vm_ip>:8080/

> raw
$ curl <vm_ip>:8080/ip
<source_ip>
```
