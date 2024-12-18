from collections import defaultdict

from datetime import datetime
import os
import asyncio
from app.api.activity_log import user_played_items
from app.models.activity_log import ActivityLogModel

# Configuration
JELLYFIN_URL = os.getenv("JELLYFIN_URL")
API_KEY = os.getenv("API_KEY")
USER_ID = os.getenv("USER_ID")

# datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")


async def main():
    now = datetime.now()
    ending_day_of_current_year = now.replace(
        year=now.year - 1,
        month=12,
        day=31,
        hour=23,
        minute=59,
        second=59,
        microsecond=999999,
    )

    result = await user_played_items(
        user_id=USER_ID,  # type: ignore
        base_url=JELLYFIN_URL,  # type: ignore
        api_key=API_KEY,  # type: ignore
        from_user_last_played_date=ending_day_of_current_year,
    )
    # let's filter the ones with type "VideoPlaybackStopped"
    # filtered = [act for act in result if act.type == "VideoPlaybackStopped"]

    # # let's aggregate the items by user
    # aggregated: dict[str, list[ActivityLogModel]] = defaultdict(list)
    # for activity in filtered:
    #     aggregated[activity.user_id].append(activity)
    # return aggregated


if __name__ == "__main__":
    asyncio.run(main())
