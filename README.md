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

For tests use Virtual Machine, Azure VM etc.  
App doesn't display **Remote IP**, neither on your **workstation** nor on your **workstation** in docker.

```
## VM example

#

make build

sudo SERVER_PORT=8080 ./whatismyip
OR
sudo SERVER_PORT=8080 ./whatismyip -inter <network_interface>

> raw
$ curl <vm_ip>:8080/
<source_ip>

> html
$ curl <vm_ip>:8080/ip


# on VM in docker 

make docker-run


## Azure ACI [toverify]

az login

make azure-rg
make azure-aci
```
