from pydantic import BaseModel, HttpUrl, Field
from typing import Optional, List
import uuid

class AttachmentBase(BaseModel):
    meetingId: str
    attachmentUrl: HttpUrl

class AttachmentCreate(AttachmentBase):
    attachmentsId: Optional[str] = Field(default_factory=lambda: str(uuid.uuid4()))

class AttachmentCreateWithoutId(AttachmentBase):
    pass

class Attachment(AttachmentBase):
    attachmentsId: str

    class Config:
        orm_mode = True

class AttachmentsIds(BaseModel):
    attachmentsIds: List[str]
