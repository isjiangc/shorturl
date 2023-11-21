package v1

import "nunu_ginblog/internal/model"

type GenerateShortUrlRequest struct {
	DestUrl  string `json:"destUrl"`
	Memo     string `json:"memo"`
	OpenType int    `json:"openType"`
}
type GenerateShortUrlData struct {
	ShortUrl string `json:"shortUrl"`
}
type GenerateShortUrlResponse struct {
	Response
	Date GenerateShortUrlData
}

type UpdateShortUrlStateRequest struct {
	Enable bool `json:"enable"`
}
type UpdateShortUrlStateData struct {
	Result bool `json:"result"`
}
type UpdateShortUrlStateResponse struct {
	Response
	Date UpdateShortUrlStateData
}

type DeleteShortUrlStateData struct {
	Result bool `json:"result"`
}
type DeleteShortUrlStateResponse struct {
	Response
	Data DeleteShortUrlStateData
}

type GetShortUrlListRequest struct {
	PageNum  int
	PageSize int
}

type GetShortUrlListResponse struct {
	Response
	Data []model.ShortUrl
}
