from mcp.server.fastmcp import FastMCP
import requests


mcp = FastMCP("AvailDA")

AVAIL_API_BASE = "http://localhost:8080"

@mcp.tool()
async def send_data_to_avail(data: str, appId: int= 89) -> str:
    """
    Send data to the AvailDA API.
    """
    url = f"{AVAIL_API_BASE}/send-data"
    data = {"AppId": appId, "Message": data}
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, json=data, headers=headers)
    
    if response.status_code == 200:
        return response.json()
    else:
        return f"Failed to send data: {response.status_code} - {response.text}"

if __name__ == "__main__":
    mcp.run(transport='stdio')