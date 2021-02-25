package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Taowyoo/cloud-album/util"
)

func main() {
	enableUpload := flag.Bool("f", false, "FILE0 FILE1 ...\nThe images to upload")
	enableList := flag.Bool("l", false, "List info of images now on cloud")
	enableListDetail := flag.Bool("ll", false, "List detail info of images now on cloud")
	flag.Parse()
	imgsPath := flag.Args()
	if *enableUpload == false && *enableList == false && *enableListDetail == false {
		fmt.Printf("Got error arguments: %v\n", os.Args)
		fmt.Println("Usage: cloud_album [-f IMG1 IMG2 ...] [-l] [-ll]")
		flag.PrintDefaults()
		return
	}
	client, _ := util.NewAWSClient("config/config.json")
	if *enableUpload {
		util.UploadImgs(client, imgsPath)
	} else if *enableList {
		util.ListImageOnCloud(client)
	} else if *enableListDetail {
		util.ListImagesDetailOnCloud(client)
	}

}
