package public

import pb "github.com/donglei1234/platform/services/iap/generated/grpc/go/iap/api"

type IOSResponse struct {
	Status  int    `json:"status"`
	Receipt string `json:"receipt"`
}

type IAPTradingRecord struct {
	UserId           string
	Token            string
	PurchaseState    int64
	ConsumptionState int64
	Sys              pb.SYS
}

type StatInfo struct {
	ProductType  int    `json:"productType"`
	ProductId    string `json:"productId"`
	ProductToken string `json:"productToken"`
	AppStoreId   string `json:"appStoreId"`
	UserId       string `json:"user_id"`
	VC           int    `json:"vc"`
	Gaid         string `json:"gaid"`
	Adid         string `json:"adid"`
	Country      string `json:"country"`
	IP           string `json:"ip"`
	OrderInfo    string `json:"orderInfo"`
}
