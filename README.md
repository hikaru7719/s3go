# s3go

s3go is cli tool to upload some files to AWS S3.  
s3go upload large file using multipart upload.  
So, s3go cut a large file and upload a lot of small file in parallel.  
Then, You can upload large file so quickly !!!  
Try s3go command with your local large file.

## Usage

Before you run cli , set below environment varable.  
s3go command uses AWS access key and secret for API request.  
if you didn't key, get your key and secret by accessing AWS console.

```
export AWS_ACCESS_KEY_ID = your access key id
export AWS_ACCESS_KEY_SECRET = your AWS access key secret
export AWS_DEFAULT_REGION = target AWS region
```

Then, You can use s3go command !!  
s3go command usage is below.

```
NAME:
   s3go - Upload some file to AWS S3

USAGE:
   s3go [global options]

VERSION:
   0.0.0

GLOBAL OPTIONS:
   --file File, -f File                        File to upload to S3
   --bucket S3 bucket Name, -b S3 bucket Name  S3 bucket Name to upload files
   --help, -h                                  show help
   --version, -v                               print the version
```
