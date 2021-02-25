package util

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestListImageOnCloud(t *testing.T) {
	type args struct {
		client *s3.Client
	}
	cli, _ := NewAWSClient("../config/config.json")
	tests := []struct {
		name string
		args args
	}{
		{"Test", args{
			cli,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ListImageOnCloud(tt.args.client)
		})
	}
}

func TestNewAWSClient(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test",
			args{
				"../config/config.json",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAWSClient(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAWSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestListImagesDetailOnCloud(t *testing.T) {
	type args struct {
		client *s3.Client
	}
	cli, _ := NewAWSClient("../config/config.json")
	tests := []struct {
		name string
		args args
	}{
		{"Test", args{
			cli,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ListImagesDetailOnCloud(tt.args.client)
		})
	}
}

func TestUploadImg(t *testing.T) {
	type args struct {
		client  *s3.Client
		imgPath string
	}
	cli, _ := NewAWSClient("../config/config.json")
	tests := []struct {
		name string
		args args
	}{
		{"Test", args{
			cli,
			"../example_imgs/turkey-bodrum.jpg",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UploadImg(tt.args.client, tt.args.imgPath)
		})
	}
}

func TestUploadImgs(t *testing.T) {
	type args struct {
		client   *s3.Client
		imgPaths []string
	}
	cli, _ := NewAWSClient("../config/config.json")
	tests := []struct {
		name string
		args args
	}{
		{"Test", args{
			cli,
			[]string{"../example_imgs/england-london-bridge.jpg", "../example_imgs/germany-garching-heide.jpg"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UploadImgs(tt.args.client, tt.args.imgPaths)
		})
	}
}
