from typing import Optional
from pydantic import BaseModel, Field

class ActivityLogModel(BaseModel):
    id: int = Field(..., alias="Id")
    name: str = Field(..., alias="Name")
    type: str = Field(..., alias="Type")
    date: str = Field(..., alias="Date")
    user_id: str = Field(..., alias="UserId")
    severity: str = Field(..., alias="Severity")
    short_overview: Optional[str] = Field(None, alias="ShortOverview")
    item_id: Optional[str] = Field(None, alias="ItemId")

    class Config:
        populate_by_name = True