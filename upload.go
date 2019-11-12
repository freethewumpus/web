package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/julienschmidt/httprouter"
)

func SendUnauthorized(w http.ResponseWriter) {
	headers := w.Header()
	headers.Add("Content-Type", "application/json")

	w.WriteHeader(403)
	_, err := w.Write([]byte(`{
		"success": false,
		"error": "Unauthorized."
	}`))
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	host := r.Host
	MultipartHeader := r.Header.Get("Content-Type")
	if strings.Split(MultipartHeader, ";")[0] != "multipart/form-data" {
		headers := w.Header()
		headers.Add("Content-Type", "application/json")

		w.WriteHeader(400)
		_, err := w.Write([]byte(`{
			"success": false,
			"error": "This needs to be ran as multipart/form-data."
		}`))
		if err != nil {
			log.Print(err)
			w.WriteHeader(500)
		}
		return
	}
	if host != Host {
		headers := w.Header()
		headers.Add("Content-Type", "application/json")

		w.WriteHeader(400)
		_, err := w.Write([]byte(`{
			"success": false,
			"error": "This needs to be ran from the freethewump.us domain."
		}`))
		if err != nil {
			log.Print(err)
			w.WriteHeader(500)
			return
		}
	} else {
		auth := r.Header.Get("Token")
		if auth == "" {
			SendUnauthorized(w)
		} else {
			domain, uid, NamingScheme, encryption, UserSuccess := GetUser(auth)
			if !UserSuccess {
				SendUnauthorized(w)
			} else {
				DomainInformation, success := GetDomain(domain)
				if !success {
					SendUnauthorized(w)
				} else {
					file, header, err := r.FormFile("file")
					if err == http.ErrMissingFile {
						w.WriteHeader(400)
						_, err := w.Write([]byte(`{
							"success": false,
							"error": "The file is missing."
						}`))
						if err != nil {
							log.Print(err)
							w.WriteHeader(500)
							return
						}
					} else if err != nil {
						panic(err)
					} else {
						defer file.Close()
						var bucket S3Bucket
						if DomainInformation.Bucket != nil {
							bucket = *DomainInformation.Bucket
						} else {
							bucket = S3Bucket{
								SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
								Bucket:          os.Getenv("S3_BUCKET"),
								AccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
								Endpoint:        os.Getenv("S3_ENDPOINT"),
								Region:          os.Getenv("S3_REGION"),
							}
						}
						if (!DomainInformation.Public && !StringInSlice(uid, DomainInformation.Whitelist)) || (DomainInformation.Public && StringInSlice(uid, DomainInformation.Blacklist)) {
							headers := w.Header()
							headers.Add("Content-Type", "application/json")
							w.WriteHeader(403)
							_, err := w.Write([]byte(`{
								"success": false,
								"error": "You do not have access to this bucket."
							}`))
							if err != nil {
								log.Print(err)
								w.WriteHeader(500)
								return
							}
						} else {
							filename := CreateFilename(NamingScheme)
							StaticCredential := credentials.NewStaticCredentials(bucket.AccessKeyId, bucket.SecretAccessKey, "")
							s3sess := session.Must(session.NewSession(&aws.Config{
								Endpoint:    &bucket.Endpoint,
								Credentials: StaticCredential,
								Region:      &bucket.Region,
							}))
							svc := s3.New(s3sess)
							key := fmt.Sprintf("%s/%s", domain, filename)
							FileSplit := strings.Split(header.Filename, ".")
							FileType := strings.ToLower(FileSplit[len(FileSplit)-1])
							MimeType := mime.TypeByExtension(FileType)
							DecryptionBit := ""
							var FileReader io.ReadSeeker
							FileReader = file
							if encryption {
								MimeType = "encrypted/" + MimeType
								DecryptionKey := CreateFilename("cccccccccccccccccccccccccccccccc")
								DecryptionBit = "?key=" + DecryptionKey
								c, _ := aes.NewCipher([]byte(DecryptionKey))
								gcm, _ := cipher.NewGCM(c)
								nonce := make([]byte, gcm.NonceSize())
								io.ReadFull(rand.Reader, nonce)
								b, err := ioutil.ReadAll(file)
								if err != nil {
									w.WriteHeader(500)
									_, err := w.Write([]byte(`{
										"success": false,
										"error": "Failed to download the buffer."
									}`))
									if err != nil {
										log.Print(err)
									}
									return
								}
								b = gcm.Seal(nonce, nonce, b, nil)
								FileReader = bytes.NewReader(b)
							}
							UploadParams := &s3.PutObjectInput{
								Bucket:             &bucket.Bucket,
								Key:                &key,
								ContentType:        &MimeType,
								Body:               FileReader,
								ACL:                aws.String("private"),
								ContentLength:      aws.Int64(header.Size),
								ContentDisposition: aws.String("attachment"),
							}
							_, err := svc.PutObject(UploadParams)
							if err != nil {
								headers := w.Header()
								headers.Add("Content-Type", "application/json")
								fmt.Print(err)

								w.WriteHeader(500)
								_, err := w.Write([]byte(`{
									"success": false,
									"error": "Failed to upload to the specified S3 bucket."
								}`))
								if err != nil {
									log.Print(err)
									return
								}
							} else {
								headers := w.Header()
								headers.Add("Content-Type", "application/json")
								_, err := w.Write([]byte(fmt.Sprintf(`{
									"success": true,
									"url": "https://%s/%s.%s%s"
								}`, domain, filename, FileType, DecryptionBit)))
								if err != nil {
									log.Print(err)
									w.WriteHeader(500)
									return
								}
							}
						}
					}
				}
			}
		}
	}
}
