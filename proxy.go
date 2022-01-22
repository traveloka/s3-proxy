package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var s3Bucket = os.Getenv("S3_BUCKET")
var s3Prefix = os.Getenv("S3_PREFIX")

func handleProxy(w http.ResponseWriter, r *http.Request) {
	var ifNoneMatch *string

	if oldEtag := r.Header.Get("If-None-Match"); oldEtag != "" {
		ifNoneMatch = &oldEtag
	}

	res, err := s3Client.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket:      &s3Bucket,
		Key:         aws.String(s3Prefix + r.URL.Path),
		IfNoneMatch: ifNoneMatch,
	})

	if err != nil {
		var respErr *awshttp.ResponseError
		if errors.As(err, &respErr) && respErr.Response.StatusCode == 304 {
			w.WriteHeader(304)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	w.Header().Set("etag", *res.ETag)
	w.Header().Set("cache-control", *res.CacheControl)
	w.Header().Set("content-type", *res.ContentType)
	io.Copy(w, res.Body)
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	s3Client = s3.NewFromConfig(cfg)
}
