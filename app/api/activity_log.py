from datetime import datetime
from typing import Any, Optional

import httpx


async def user_played_items(
    user_id: str,
    base_url: str,
    api_key: str,
    from_user_last_played_date: Optional[datetime] = None,
    **extra_headers: dict[str, Any],
) -> list[dict]:
    params = {
        "userId": user_id,
        "isPlayed": True,
        "recursive": True,
        "IncludeItemTypes": "Movie,Episode",
    }
    endpoint = f"{base_url}/Items"
    headers = {"X-Emby-Token": api_key, **extra_headers}

    async with httpx.AsyncClient() as client:
        response = await client.get(endpoint, headers=headers, params=params)
        response.raise_for_status()
        activity = response.json()


    if from_user_last_played_date:
        threshold_date = f"{from_user_last_played_date.isoformat()}Z"
        filtered = []
        for act in activity["Items"]:
            if act["UserData"].get("LastPlayedDate") and act["UserData"]["LastPlayedDate"] > threshold_date:
                filtered.append(act)
        return filtered
    return activity
