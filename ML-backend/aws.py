import os
import boto3
from dotenv import load_dotenv
load_dotenv()
AWS_REGION = os.getenv("AWS_REGION")
S3_BUCKET_NAME = os.getenv("S3_BUCKET_NAME")
s3=boto3.client(
    's3',
    region_name=AWS_REGION,
)
def upload_file(filePath,fileName):
    s3.upload_file(
    filePath,
    S3_BUCKET_NAME,
    f"modelWeights/{fileName}")
    os.remove(filePath)
