# MCP Server on AVAIL DA
Interact with the avail DA via Claude Desktop using natural language prompts 

## Work In Progress
- Currently supports sending data
### Future Work 
- Query transactions on Avail DA
- Integrate Avail Nexus 
- Seamless cross-chain transfers 
 
## Setup 
### Running the local avail-client handler backend 
```sh
cd avail-client
cp .env.example .env # configure your .env file
go mod tidy
go run main.go
```
Starts the server locally on 8080 port

#### Configuring the MCP Server to Claude desktop
Add this to `claude_desktop_config.json` present somewhere in `Application Support`

- for mac-os it is in :` /Library/Application\ Support/Claude`
```sh
"AvailDA": {
    "command": "python",
    "args": [
        "FULL/PATH/TO/REPO/AvailMCP/avail-mcp/server.py"
    ]
}
```
