![vncparty logo](./logo.png)
## vncparty
shareable vnc connections in the browser
### Build & Run
* Configure a fresh Postgres database
* Configure .env ([example](./.env.default))
* `go build -o vncparty && ./vncparty`

Works best with QEMU's VNC frontend, use something like this command:

    qemu-system-x86_64 -hda hdd.qcow2 -m 4G --enable-kvm -display none -vnc :1,websocket=15901 -k en-us

![vncparty in action](./sc2.png)

## Roadmap
- [x] Proxy VNC connections
- [ ] Docker image
- [ ] In-room chat
- [ ] MCP support
- [ ] Help text
- [ ] Docker image



