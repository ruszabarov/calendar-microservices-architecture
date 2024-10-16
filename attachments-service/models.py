from sqlalchemy import Column, String
from database import Base


class Attachment(Base):
    __tablename__ = "attachments"

    attachmentsId = Column(String, primary_key=True, index=True)
    meetingId = Column(String, index=True)
    attachmentUrl = Column(String(2000))
