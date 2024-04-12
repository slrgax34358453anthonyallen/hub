package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/getAlby/nostr-wallet-connect/nip47"
	"github.com/nbd-wtf/go-nostr"
	decodepay "github.com/nbd-wtf/ln-decodepay"
	"github.com/sirupsen/logrus"
)

func (svc *Service) HandleLookupInvoiceEvent(ctx context.Context, nip47Request *Nip47Request, requestEvent *RequestEvent, app *App, publishResponse func(*Nip47Response, nostr.Tags)) {

	lookupInvoiceParams := &Nip47LookupInvoiceParams{}
	resp := svc.decodeNip47Request(nip47Request, requestEvent, app, lookupInvoiceParams)
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
		"invoice":             lookupInvoiceParams.Invoice,
		"paymentHash":         lookupInvoiceParams.PaymentHash,
	}).Info("Looking up invoice")

	paymentHash := lookupInvoiceParams.PaymentHash

	if paymentHash == "" {
		paymentRequest, err := decodepay.Decodepay(strings.ToLower(lookupInvoiceParams.Invoice))
		if err != nil {
			svc.Logger.WithFields(logrus.Fields{
				"requestEventNostrId": requestEvent.NostrId,
				"appId":               app.ID,
				"invoice":             lookupInvoiceParams.Invoice,
			}).Errorf("Failed to decode bolt11 invoice: %v", err)

			publishResponse(&Nip47Response{
				ResultType: nip47Request.Method,
				Error: &Nip47Error{
					Code:    nip47.ERROR_INTERNAL,
					Message: fmt.Sprintf("Failed to decode bolt11 invoice: %s", err.Error()),
				},
			}, nostr.Tags{})
			return
		}
		paymentHash = paymentRequest.PaymentHash
	}

	transaction, err := svc.lnClient.LookupInvoice(ctx, paymentHash)
	if err != nil {
		svc.Logger.WithFields(logrus.Fields{
			"requestEventNostrId": requestEvent.NostrId,
			"appId":               app.ID,
			"invoice":             lookupInvoiceParams.Invoice,
			"paymentHash":         lookupInvoiceParams.PaymentHash,
		}).Infof("Failed to lookup invoice: %v", err)

		publishResponse(&Nip47Response{
			ResultType: nip47Request.Method,
			Error: &Nip47Error{
				Code:    nip47.ERROR_INTERNAL,
				Message: err.Error(),
			},
		}, nostr.Tags{})
		return
	}

	responsePayload := &Nip47LookupInvoiceResponse{
		Nip47Transaction: *transaction,
	}

	publishResponse(&Nip47Response{
		ResultType: nip47Request.Method,
		Result:     responsePayload,
	}, nostr.Tags{})
}
