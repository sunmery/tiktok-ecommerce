package biz

type Request struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type DeleteAddressesRequest struct {
	AddressId uint32 `json:"address_id"`
	Owner     string `json:"owner"`
	Name      string `json:"name"`
}

type Address struct {
	Id            uint32 `json:"id"`
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
}

type Addresses struct {
	Addresses []*Address `json:"addresses"`
}

type DeleteAddressesReply struct {
	Message string `json:"message"`
	Id      uint32 `json:"id"`
	Code    uint32 `json:"code"`
}
