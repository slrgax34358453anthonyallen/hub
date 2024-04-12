package main

import (
	"context"

	"github.com/getAlby/nostr-wallet-connect/nip47"
	"github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
)

const (
	MSAT_PER_SAT = 1000
)

func (svc *Service) HandleGetBalanceEvent(ctx context.Context, nip47Request *Nip47Request, requestEvent *RequestEvent, app *App, publishResponse func(*Nip47Response, nostr.Tags)) {

	resp := svc.checkPermission(nip47Request, requestEvent.NostrId, app, 0)
	if resp != nil {
		publishResponse(resp, nostr.Tags{})
		return
	}

	svc.Logger.WithFields(logrus.Fields{
		"requestEventNostrId": requestEvent.NostrId,
		"appId":               app.ID,
	}).Info("Fetching balance")

	balance, err := svc.lnClient.GetBalance(ctx)
	if err != nil {
		svc.Logger.WithFields(logrus.Fields{
			"requestEventNostrId": requestEvent.NostrId,
			"appId":               app.ID,
		}).Infof("Failed to fetch balance: %v", err)
		publishResponse(&Nip47Response{
			ResultType: nip47Request.Method,
			Error: &Nip47Error{
				Code:    nip47.ERROR_INTERNAL,
				Message: err.Error(),
			},
		}, nostr.Tags{})
		return
	}

	responsePayload := &Nip47BalanceResponse{
		Balance: balance,
	}

	appPermission := AppPermission{}
	svc.db.Where("app_id = ? AND request_method = ?", app.ID, nip47.PAY_INVOICE_METHOD).First(&appPermission)

	maxAmount := appPermission.MaxAmount
	if maxAmount > 0 {
		responsePayload.MaxAmount = maxAmount * MSAT_PER_SAT
		responsePayload.BudgetRenewal = appPermission.BudgetRenewal
	}

	publishResponse(&Nip47Response{
		ResultType: nip47Request.Method,
		Result:     responsePayload,
	}, nostr.Tags{})
}
