import os
import asyncio
from app.api.activity_log import get_activities

# Configuration
JELLYFIN_URL = os.getenv("JELLYFIN_URL")
API_KEY = os.getenv("API_KEY")

async def main():
    result = await get_activities(JELLYFIN_URL, API_KEY) # type: ignore
    return result

if __name__ == "__main__":
    asyncio.run(main())