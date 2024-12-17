from collections import defaultdict

import os
import asyncio
from app.api.activity_log import get_activities
from app.models.activity_log import ActivityLogModel

# Configuration
JELLYFIN_URL = os.getenv("JELLYFIN_URL")
API_KEY = os.getenv("API_KEY")


async def main():
    result = await get_activities(JELLYFIN_URL, API_KEY)  # type: ignore
    # let's filter the ones with type "VideoPlaybackStopped"
    filtered = [act for act in result if act.type == "VideoPlaybackStopped"]

    # let's aggregate the items by user
    aggregated: dict[str, list[ActivityLogModel]] = defaultdict(list)
    for activity in filtered:
        aggregated[activity.user_id].append(activity)
    return aggregated


if __name__ == "__main__":
    asyncio.run(main())