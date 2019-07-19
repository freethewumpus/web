package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"strings"
)

func View(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	Host := strings.ToLower(r.Host)
	domain, success := GetDomain(Host)
	if !success {
		w.WriteHeader(400)
		_, err := w.Write([]byte(`
			<!DOCTYPE HTML>
			<html>
				<head>
					<title>Domain Not In Database</title>
				</head>
				<body>
					<h1>Domain Not In Database</h1>
					<p>Please contact support to fix.</p>
				</body>
			</html>
		`))
		if err != nil {
			panic(err)
		}
	} else {
		var bucket S3Bucket
		if domain.Bucket != nil {
			bucket = *domain.Bucket
		} else {
			bucket = S3Bucket{
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
				Bucket: os.Getenv("S3_BUCKET"),
				AccessKeyId: os.Getenv("AWS_ACCESS_KEY_ID"),
				Endpoint: os.Getenv("S3_ENDPOINT"),
				Region: os.Getenv("S3_REGION"),
			}
		}

		StaticCredential := credentials.NewStaticCredentials(bucket.AccessKeyId, bucket.SecretAccessKey, "")
		s3sess := session.Must(session.NewSession(&aws.Config{
			Endpoint: &bucket.Endpoint,
			Credentials: StaticCredential,
			Region: &bucket.Region,
		}))

		ImageName := strings.Split(p.ByName("image"), ".")[0]
		Key := fmt.Sprintf("%s/%s", Host, ImageName)
		svc := s3.New(s3sess)

		GetParams := &s3.GetObjectInput{
			Bucket: &bucket.Bucket,
			Key: &Key,
		}

		result, err := svc.GetObject(GetParams)

		if err == nil {
			_, err := io.Copy(w, result.Body)
			if err != nil {
				panic(err)
			}
		} else if result != nil && result.Body == nil {
			w.WriteHeader(404)
			_, err := w.Write([]byte(`
				<!DOCTYPE HTML>
				<html>
					<head>
						<title>Not Found</title>
					</head>
					<body>
						<h1>Not Found</h1>
						<p>The content you are trying to access cannot be found.</p>
					</body>
				</html>
			`))
			if err != nil {
				panic(err)
			}
		} else {
			w.WriteHeader(500)
			_, err := w.Write([]byte(`
				<!DOCTYPE HTML>
				<html>
					<head>
						<title>Could Not Access S3 Bucket</title>
					</head>
					<body>
						<h1>Could Not Access S3 Bucket</h1>
						<p>Please contact support to fix.</p>
					</body>
				</html>
			`))
			if err != nil {
				panic(err)
			}
		}
	}
}
