from fastapi import FastAPI, Depends, HTTPException
from sqlalchemy.orm import Session
import requests

import crud
import models
import schemas
from database import SessionLocal, engine

# Creating the database tables
models.Base.metadata.create_all(bind=engine)

app = FastAPI(redirect_slashes=True)


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


MEETINGS_MICROSERVICE_URL = "http://meetings-service:8080"


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

    # Update the Meeting's List of Attachment Ids
    try:
        response = requests.post(
            f"{MEETINGS_MICROSERVICE_URL}/meetings/{attachment.meetingId}/attachments",
            json={"attachmentsId": db_attachment.attachmentsId},
        )
        response.raise_for_status()
    except requests.RequestException as e:
        # Rollback if the meeting update fails
        db.delete(db_attachment)
        db.commit()
        raise HTTPException(status_code=400, detail="Failed to update Meetings service") from e
    return db_attachment


# Get Attachment by ID
@app.get("/attachments/{attachmentsId}", response_model=schemas.Attachment)
def read_attachment(attachmentsId: str, db: Session = Depends(get_db)):
    db_attachment = crud.get_attachment(db, attachmentsId)
    if db_attachment is None:
        raise HTTPException(status_code=404, detail="Attachment not found")
    return db_attachment


# Get All Attachments
@app.get("/attachments", response_model=list[schemas.Attachment])
def read_attachments(skip: int = 0, limit: int = 100, db: Session = Depends(get_db)):
    attachments = crud.get_attachments(db, skip=skip, limit=limit)
    return attachments


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
def delete_attachment(attachmentsId: str, db: Session = Depends(get_db)):
    db_attachment = crud.get_attachment(db, attachmentsId)
    if db_attachment is None:
        raise HTTPException(status_code=404, detail="Attachment not found")
    # Update the Meeting's List of Attachment Ids
    try:
        response = requests.delete(
            f"{MEETINGS_MICROSERVICE_URL}/meetings/{db_attachment.meetingId}/attachments/{attachmentsId}"
        )
        response.raise_for_status()
    except requests.RequestException as e:
        raise HTTPException(status_code=400, detail="Failed to update Meetings service") from e
    crud.delete_attachment(db, attachmentsId)
    return {"detail": "Attachment deleted"}
