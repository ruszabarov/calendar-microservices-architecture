import uuid

from sqlalchemy.orm import Session
from models import Attachment
from schemas import AttachmentCreate


def get_attachment(db: Session, attachmentsId: str):
    return db.query(Attachment).filter(Attachment.attachmentsId == attachmentsId).first()


def get_attachments(db: Session, skip: int = 0, limit: int = 100):
    return db.query(Attachment).offset(skip).limit(limit).all()


def create_attachment(db: Session, attachment: AttachmentCreate):
    db_attachment = Attachment(
        attachmentsId=attachment.attachmentsId or str(uuid.uuid4()),
        meetingId=attachment.meetingId,
        attachmentUrl=str(attachment.attachmentUrl),
    )
    db.add(db_attachment)
    db.commit()
    db.refresh(db_attachment)
    return db_attachment


def delete_attachment(db: Session, attachmentsId: str):
    db_attachment = get_attachment(db, attachmentsId)
    if db_attachment:
        db.delete(db_attachment)
        db.commit()
        return True
    return False


def update_attachment(db: Session, attachmentsId: str, attachment: AttachmentCreate):
    db_attachment = get_attachment(db, attachmentsId)
    if db_attachment:
        db_attachment.meetingId = attachment.meetingId
        db_attachment.attachmentUrl = str(attachment.attachmentUrl)
        db.commit()
        db.refresh(db_attachment)
        return db_attachment
    return None
