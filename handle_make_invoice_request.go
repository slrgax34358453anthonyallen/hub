package main

import (
	"context"

	"github.com/getAlby/nostr-wallet-connect/nip47"
	"github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
)

func (svc *Service) HandleMakeInvoiceEvent(ctx context.Context, nip47Request *Nip47Request, requestEvent *RequestEvent, app *App, publishResponse func(*Nip47Response, nostr.Tags)) {

	makeInvoiceParams := &Nip47MakeInvoiceParams{}
	resp := svc.decodeNip47Request(nip47Request, requestEvent, app, makeInvoiceParams)
	if resp != nil {
		publishResponse(resp, nostr.Tags{})
		return
	}

	resp = svc.checkPermission(nip47Request, requestEvent.NostrId, app, 0)
	if resp != nil {
		publishResponse(resp, nostr.Tags{})
		return
	}

	svc.Logger.WithFields(logrus.Fields{
		"requestEventNostrId": requestEvent.NostrId,
		"appId":               app.ID,
		"amount":              makeInvoiceParams.Amount,
		"description":         makeInvoiceParams.Description,
		"descriptionHash":     makeInvoiceParams.DescriptionHash,
		"expiry":              makeInvoiceParams.Expiry,
	}).Info("Making invoice")

	expiry := makeInvoiceParams.Expiry
	if expiry == 0 {
		expiry = 86400
	}

	transaction, err := svc.lnClient.MakeInvoice(ctx, makeInvoiceParams.Amount, makeInvoiceParams.Description, makeInvoiceParams.DescriptionHash, expiry)
	if err != nil {
		svc.Logger.WithFields(logrus.Fields{
			"requestEventNostrId": requestEvent.NostrId,
			"appId":               app.ID,
			"amount":              makeInvoiceParams.Amount,
			"description":         makeInvoiceParams.Description,
			"descriptionHash":     makeInvoiceParams.DescriptionHash,
			"expiry":              makeInvoiceParams.Expiry,
		}).Infof("Failed to make invoice: %v", err)

		publishResponse(&Nip47Response{
			ResultType: nip47Request.Method,
			Error: &Nip47Error{
				Code:    nip47.ERROR_INTERNAL,
				Message: err.Error(),
			},
		}, nostr.Tags{})
		return
	}

	responsePayload := &Nip47MakeInvoiceResponse{
		Nip47Transaction: *transaction,
	}

	publishResponse(&Nip47Response{
		ResultType: nip47Request.Method,
		Result:     responsePayload,
	}, nostr.Tags{})
}
