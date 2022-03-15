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

```
# docker (if eth* available)

make docker-build
make docker-run

curl localhost:8080/
curl localhost:8080/ip

make docker-stop


# VM

make build

sudo SERVER_PORT=8080 ./whatismyip          // if eth* available
OR
sudo SERVER_PORT=8080 ./whatismyip -inter <network_interface>

curl <vm_ip>:8080/
curl <vm_ip>:8080/ip
```
