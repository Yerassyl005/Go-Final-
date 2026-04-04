package service

import "github.com/skip2/go-qrcode"

type QRService struct{}

func NewQRService() *QRService {
	return &QRService{}
}

func (s *QRService) GenerateQR(text string) ([]byte, error) {
	return qrcode.Encode(text, qrcode.Medium, 256)
}
