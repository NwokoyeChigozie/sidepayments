package payment

import (
	"net/http"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
)

func ListBanksService(extReq request.ExternalRequest, db postgresql.Databases, countryCode string) ([]external_models.BanksResponse, int, error) {
	var (
		rave = Rave{ExtReq: extReq}
	)

	banks, err := rave.ListBanks(countryCode)
	if err != nil {
		return banks, http.StatusInternalServerError, err
	}

	return banks, http.StatusOK, nil
}

func ConvertCurrencyService(extReq request.ExternalRequest, db postgresql.Databases, amount float64, from, to string) (models.ConvertCurrencyResponse, int, error) {
	var (
		rave = Rave{ExtReq: extReq}
	)
	conversionData, err := rave.ConvertCurrency(amount, from, to)
	if err != nil {
		return models.ConvertCurrencyResponse{}, http.StatusBadRequest, err
	}
	return conversionData, http.StatusOK, nil
}
