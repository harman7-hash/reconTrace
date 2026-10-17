from uuid import UUID
import requests
from fastapi import FastAPI,HTTPException
from model import train_model as trainer
from fastapi.responses import StreamingResponse
from dotenv import load_dotenv
from celery import Celery
from pydantic import BaseModel
from pathlib import Path
import os
from aws import upload_file
load_dotenv()
class TrainingRequest(BaseModel):
    user_id:str
CALLBACK_URL = os.getenv("CALLBACK_URL")
celery = Celery(
    "myapp",
    broker="redis://127.0.0.1:6379/0",
    backend="redis://127.0.0.1:6379/1"
)
@celery.task
def aws_queue(user_id,temp_file_path,filename,threshold):
    upload_file(temp_file_path,filename)
    webhook_queue.delay(user_id,threshold,filename)
@celery.task
def webhook_queue(user_id,threshold,filename):
    requests.post(
        CALLBACK_URL,
        json={
            "user_id":user_id,
            "threshold":float(threshold),
            "filename":filename,
            "status":"COMPLETED"
        }
    )
@celery.task
def train_model(user_id:str):
    try:
        print("TRAINING STARTED")
        THRESHOLD,model,temp_file_path,filename=trainer(user_id)
        aws_queue.delay(user_id,temp_file_path,filename,float(THRESHOLD))
        print("THRESHOLD:",float(THRESHOLD))
        print("TRAINING KHATAM")
    except Exception as e:
        print("TRAINING FAILED")
        print(e)
        requests.post(
        CALLBACK_URL,
        json={
            "user_id":user_id,
            "threshold":0,
            "filename":"NA",
            "status":"FAILED"
        }
    )
        
app=FastAPI()
@app.get('/')
def home():
    return {'message':'Welcome to my FastAPI'}
@app.post('/train')
def train_m(user_id:TrainingRequest):
    task=train_model.delay(user_id.user_id)
    print("User ID:",user_id.user_id)
    return{
        "user_id":user_id.user_id,
        "job_id":task.id,
        "status":"queued"
    }
