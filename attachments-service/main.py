from fastapi import FastAPI, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from typing import List, Optional

import crud
import models
import schemas
from database import SessionLocal, engine

# Create the database tables
models.Base.metadata.create_all(bind=engine)

app = FastAPI(redirect_slashes=True)


# Dependency to get DB session
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


@app.post("/attachments/by-ids", response_model=List[schemas.Attachment])
def get_attachments_by_ids(
        request: schemas.AttachmentsIds,
        db: Session = Depends(get_db)
):
    attachments = crud.get_attachments_by_ids(db, request.attachmentsIds)
    return attachments



# Get All Attachments
@app.get("/attachments", response_model=List[schemas.Attachment])
def read_attachments(
        skip: int = 0,
        limit: int = 100,
        db: Session = Depends(get_db)
):
    attachments = crud.get_attachments(db, skip=skip, limit=limit)
    return attachments


# Create Attachment with specified attachmentsId
@app.post("/attachments/{attachmentsId}", response_model=schemas.Attachment)
def create_attachment_with_id(
        attachmentsId: str,
        attachment: schemas.AttachmentCreateWithoutId,
        db: Session = Depends(get_db)
):
    # Check if an attachment with the given ID already exists
    existing_attachment = crud.get_attachment(db, attachmentsId)
    if existing_attachment:
        raise HTTPException(status_code=400, detail="Attachment with this ID already exists")

    # Create a new attachment with the specified ID
    new_attachment = schemas.AttachmentCreate(
        attachmentsId=attachmentsId,
        meetingId=attachment.meetingId,
        attachmentUrl=attachment.attachmentUrl,
    )
    db_attachment = crud.create_attachment(db, new_attachment)

    return db_attachment


# Get Attachment by ID
@app.get("/attachments/{attachmentsId}", response_model=schemas.Attachment)
def read_attachment(
        attachmentsId: str,
        db: Session = Depends(get_db)
):
    db_attachment = crud.get_attachment(db, attachmentsId)
    if db_attachment is None:
        raise HTTPException(status_code=404, detail="Attachment not found")
    return db_attachment


# Update Attachment
@app.put("/attachments/{attachmentsId}", response_model=schemas.Attachment)
def update_attachment(
        attachmentsId: str,
        attachment: schemas.AttachmentCreateWithoutId,
        db: Session = Depends(get_db)
):
    db_attachment = crud.update_attachment(db, attachmentsId, attachment)
    if db_attachment is None:
        raise HTTPException(status_code=404, detail="Attachment not found")
    return db_attachment


# Delete Attachment
@app.delete("/attachments/{attachmentsId}")
def delete_attachment(
        attachmentsId: str,
        db: Session = Depends(get_db)
):
    db_attachment = crud.get_attachment(db, attachmentsId)
    if db_attachment is None:
        raise HTTPException(status_code=404, detail="Attachment not found")
    crud.delete_attachment(db, attachmentsId)
    return {"detail": "Attachment deleted"}
