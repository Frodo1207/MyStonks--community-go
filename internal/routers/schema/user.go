package schema

import "MyStonks-go/internal/common/response"

type UserResp struct {
	response.Response
	Data struct {
		ID         uint   `json:"id"`
		SolAddress string `json:"sol_address"`
		Username   string `json:"username"`
	} `json:"data"`
}
