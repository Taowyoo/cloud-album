package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var bucketName = "csye6225a2test01"
var urlPrefix = fmt.Sprintf("https://%s.s3-us-west-1.amazonaws.com/", bucketName)
var readPre = "Everyone"
var fileType = "image/jpeg"

type cfg struct {
	Regin     string
	AccessKey string
	SecretKey string
}

func readFile(path string) (data []byte, err error) {
	f, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Printf("Open %s error:\n%s\n", path, err)
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer f.Close()
	data, err = ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Read %s error:\n%s\n", path, err)
		return nil, err
	}
	return
}

func readConfig(configPath string) cfg {
	var myCfg cfg
	cfgData, err := readFile(configPath)
	if err != nil {
		fmt.Printf("Fail to read %s error:\n%s\n", configPath, err)
		return myCfg
	}
	err = json.Unmarshal(cfgData, &myCfg)
	if err != nil {
		fmt.Printf("Parse %s error:\n%s\n", configPath, err)
		return myCfg
	}
	fmt.Println("Loaded config from", configPath)
	return myCfg
}

// NewAWSClient construct a new NewAWSClient object
func NewAWSClient(configPath string) (*s3.Client, error) {
	myCfg := readConfig(configPath)
	os.Setenv("AWS_ACCESS_KEY_ID", myCfg.AccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", myCfg.SecretKey)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(myCfg.Regin))
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := s3.NewFromConfig(cfg)
	return client, nil
}

// UploadImg upload one image to s3 by given path string of the file
func UploadImg(client *s3.Client, imgPath string) {
	// Open image
	file, err := os.Open(imgPath)
	defer file.Close()
	if err != nil {
		fmt.Println("Unable to open file " + imgPath)
		return
	}

	// Process image into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	// Get metadata of the image
	file.Seek(0, io.SeekStart)
	meta, err := GetExifData(file)
	if err != nil {
		log.Panic(err)
	}
	// Generate the api input
	fileName := filepath.Base(imgPath)
	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileName),
		Body:          fileBytes,
		Metadata:      meta,
		ACL:           types.ObjectCannedACLPublicRead,
		ContentLength: *aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	_, err = PutFile(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error uploading file:")
		fmt.Println(err)
		return
	}
	fmt.Println("Uploaded:", imgPath)
}

// UploadImgs upload multi images to s3 by given path string array of the file
func UploadImgs(client *s3.Client, imgPaths []string) {
	for _, path := range imgPaths {
		UploadImg(client, path)
	}
}

// ListImageOnCloud list  brief info of all files on s3
func ListImageOnCloud(client *s3.Client) {
	bucket := &bucketName
	input := &s3.ListObjectsV2Input{
		Bucket: bucket,
	}

	resp, err := GetObjects(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error retrieving list of Images:")
		fmt.Println(err)
		return
	}

	fmt.Printf("Images:\n")

	for i, item := range resp.Contents {
		fmt.Printf("=== Image %06d Begin ===\n", i)
		fmt.Println("Name:          ", *item.Key)
		fmt.Println("Last modified: ", *item.LastModified)
		fmt.Println("Size:          ", item.Size)
		fmt.Println("Storage class: ", item.StorageClass)
		fmt.Printf("=== Image %06d End ===\n", i)

	}

	fmt.Println("Found", len(resp.Contents), "images", *bucket)
}

// ListImagesDetailOnCloud retrive metadata of all images on cloud
func ListImagesDetailOnCloud(client *s3.Client) {
	bucket := &bucketName

	input := &s3.ListObjectsV2Input{
		Bucket: bucket,
	}

	resp, err := GetObjects(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error retrieving list of Images:")
		fmt.Println(err)
		return
	}
	fmt.Println("Images metadata:")
	for i, item := range resp.Contents {
		fmt.Printf("Loading Image %06d ......\n", i)
		metaData := getObjectMetadata(client, item.Key)
		fmt.Printf("=== Image %06d Begin ===\n", i)
		fmt.Println("Name:          ", *item.Key)
		fmt.Printf("URL:            %s%s\n", urlPrefix, *item.Key)
		fmt.Println("Last modified: ", *item.LastModified)
		fmt.Println("Size:          ", item.Size)
		fmt.Println("Storage class: ", item.StorageClass)
		for k, v := range metaData {
			fmt.Printf("%-15s %s\n", strings.Title(strings.ToLower(k))+":", v)
		}
		fmt.Printf("=== Image %06d End ===\n", i)
	}

}

func getObjectMetadata(client *s3.Client, key *string) map[string]string {
	bucket := &bucketName
	input := &s3.HeadObjectInput{
		Bucket: bucket,
		Key:    key,
	}

	resp, err := HeadObject(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error retrieving list of Images:")
		fmt.Println(err)
		return nil
	}
	return resp.Metadata
}
