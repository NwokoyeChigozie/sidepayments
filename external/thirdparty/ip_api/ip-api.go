package ip_api

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) ResolveIp() (external_models.ResolveIpResponse, error) {

	var (
		outBoundResponse external_models.ResolveIpResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	ip, ok := idata.(string)
	if !ok {
		logger.Info("ip-api resolve ip", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	path := "/" + ip

	logger.Info("ip-api resolve ip", ip)
	err := r.getNewSendRequestObject(nil, headers, path).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("ip-api resolve ip", outBoundResponse, err.Error())
		return outBoundResponse, err
	}

	return outBoundResponse, nil
}
