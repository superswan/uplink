# UPLINK (Ultimate Persistent Linux Implant Network Kit)

Persistent reverse shell and connection manager. Originally designed for Linux but has cross platform support. Functionally a lightweight C2. 

[Demo Video (YouTube)](https://www.youtube.com/watch?v=qzB7AJMDosY)

## Client

Main client is currently a lightweight python client. It's essentially a python reverse shell with some added features for client management and reconnection. Client will attempt to connect to the server if it fails it will sleep for some time and attempt reconnection. If the connection succeeds it sends a randomly generated `client id` to the server and starts a thread that sends a reverse shell to the socket connection. 

## Server 

The server is a small go program that can multiplex connections from clients. It can accept connections from the python script or other valid reverse shells (untested). Clients get upserted into the database backend and can be managed through the CLI or using the API.  

### API Endpoints

`GET /status/` - Returns a JSON array of Client IDs present in `activeClients` slice

`POST /command/:id` - Sends a command to client based on ID and returns the result

## Web UI
![](webui.png)

Nuxt.js application for managing client connections. Intended to be used rather than the CLI. Future versions will include more robust client management features and payload builder/distribution support.  


## Backend (PocketBase)
PocketBase is used for backend and database management. Please read their documentation for more information.

## Requirements

- Go
- Node.js
- PocketBase


## ToDo
- [ ] Windows Client (C)
- [ ] Linux Client (C)
- [ ] Netcat Client
- [ ] Add builder to WebUI
- [ ] Remove hardcoded values (Configuration)
- [ ] Installation script