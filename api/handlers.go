package api

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	fmt "fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"

	context "golang.org/x/net/context"
)

type Server struct {
	root      string
	sizeLimit int
}

var allImageType string

func NewServer(root string, sizeLimit int) (*Server, error) {
	err := createDirIfNotExist(root)
	if err != nil {
		return nil, err
	}
	return &Server{root, sizeLimit}, nil
}

func (s *Server) PutImage(ctx context.Context, req *PutRequest) (*PutResponse, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(req.Data))
	imgData, err := ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}
	if len(imgData) > s.sizeLimit {
		return nil, fmt.Errorf("file size (%d) excess the max size (%d)", len(imgData), s.sizeLimit)
	}
	processedData, format, err := processImage(imgData, int(req.Width), int(req.Height))
	if err != nil {
		return nil, err
	}
	err = createDirIfNotExist(filepath.Join(s.root, req.Class))
	if err != nil {
		return nil, err
	}
	baseName := fmt.Sprintf("%x.%s", md5.Sum(processedData), format)
	fileName := filepath.Join(s.root, req.Class, baseName)
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Write(processedData)
	if err != nil {
		return nil, err
	}
	return &PutResponse{Name: fileName}, nil
}

func (s *Server) GetImage(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	info, err := os.Stat(req.Name)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is not valid file path", req.Name)
	}
	f, err := os.Open(req.Name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(make([]byte, 0, s.sizeLimit*2))
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	defer encoder.Close()
	_, err = encoder.Write(b)
	if err != nil {
		return nil, err
	}
	return &GetResponse{Data: buffer.Bytes()}, nil
}

func resizeImage(img image.Image, format string, width, height int) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 102400))
	resized := imaging.Resize(img, width, height, imaging.Lanczos)
	if format == "jpeg" {
		err := jpeg.Encode(buf, resized, nil)
		if err != nil {
			return nil, err
		}
	} else {
		err := png.Encode(buf, resized)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil

}

func processImage(imgData []byte, width int, height int) ([]byte, string, error) {
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, "", err
	}
	if format != "jpeg" && format != "png" {
		return nil, "", errors.New("not valid image type (require JPEG or PNG)")
	}
	if width == -1 && height == -1 {
		return imgData, format, nil
	}
	resizedBytes, err := resizeImage(img, format, width, height)
	if err != nil {
		return nil, "", err
	}
	return resizedBytes, format, nil
}

func createDirIfNotExist(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0775)
			if err != nil {
				return err
			}
		}
		return err
	}
	if !info.IsDir() {
		err = os.MkdirAll(path, 0775)
		if err != nil {
			return err
		}
	}
	return nil
}
