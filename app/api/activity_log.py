from typing import Any
import httpx
from app.models.activity_log import ActivityLogModel
from urllib.parse import urlencode

async def get_activities(
    base_url: str, api_key: str, **extra_headers: dict[str, Any]
) -> list[ActivityLogModel]:
    query = urlencode({
        "limit": 100_000_000, # yes, this is an arbitrary big number
        "hasUserId": "true",
    })
    endpoint = f"{base_url}/System/ActivityLog/Entries?{query}"
    headers = {"X-Emby-Token": api_key, **extra_headers}

    async with httpx.AsyncClient() as client:
        response = await client.get(endpoint, headers=headers)
        response.raise_for_status()
        activity = response.json()
    return [ActivityLogModel(**act) for act in activity["Items"]]