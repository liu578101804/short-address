package app

type ShortLinkReq struct {
	Url 					string	`json:"url"  validate:"nonzero"`
	ExpirationInMinute		int64 	`json:"expiration_in_minute"  validate:"min=0"`
}

type ShortLinkResp struct {
	ShortLink string `json:"short_link"`
}
