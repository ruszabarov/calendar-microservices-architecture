from sqlalchemy.orm import Session
from typing import List
import models
import schemas
import uuid

def get_attachment(db: Session, attachmentsId: str):
    return db.query(models.Attachment).filter(models.Attachment.attachmentsId == attachmentsId).first()

def get_attachments(db: Session, skip: int = 0, limit: int = 100):
    return db.query(models.Attachment).offset(skip).limit(limit).all()

def get_attachments_by_ids(db: Session, attachments_ids: List[str]):
    return db.query(models.Attachment).filter(models.Attachment.attachmentsId.in_(attachments_ids)).all()

def create_attachment(db: Session, attachment: schemas.AttachmentCreate):
    db_attachment = models.Attachment(
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

def update_attachment(db: Session, attachmentsId: str, attachment: schemas.AttachmentCreateWithoutId):
    db_attachment = get_attachment(db, attachmentsId)
    if db_attachment:
        db_attachment.meetingId = attachment.meetingId
        db_attachment.attachmentUrl = str(attachment.attachmentUrl)
        db.commit()
        db.refresh(db_attachment)
        return db_attachment
    return None
