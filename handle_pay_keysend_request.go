package main

import (
	"context"
	"encoding/json"

	"github.com/getAlby/nostr-wallet-connect/events"
	"github.com/sirupsen/logrus"
)

func (svc *Service) HandlePayKeysendEvent(ctx context.Context, request *Nip47Request, requestEvent *RequestEvent, app *App) (result *Nip47Response, err error) {

	payParams := &Nip47KeysendParams{}
	err = json.Unmarshal(request.Params, payParams)
	if err != nil {
		svc.Logger.WithFields(logrus.Fields{
			"eventId": requestEvent.NostrId,
			"appId":   app.ID,
		}).Errorf("Failed to decode nostr event: %v", err)
		return nil, err
	}

	// We use pay_invoice permissions for budget and max amount
	hasPermission, code, message := svc.hasPermission(app, NIP_47_PAY_INVOICE_METHOD, payParams.Amount)

	if !hasPermission {
		svc.Logger.WithFields(logrus.Fields{
			"eventId":      requestEvent.NostrId,
			"appId":        app.ID,
			"senderPubkey": payParams.Pubkey,
		}).Errorf("App does not have permission: %s %s", code, message)

		return &Nip47Response{
			ResultType: request.Method,
			Error: &Nip47Error{
				Code:    code,
				Message: message,
			}}, nil
	}

	payment := Payment{App: *app, RequestEvent: *requestEvent, Amount: uint(payParams.Amount / 1000)}
	insertPaymentResult := svc.db.Create(&payment)
	if insertPaymentResult.Error != nil {
		return nil, insertPaymentResult.Error
	}

	svc.Logger.WithFields(logrus.Fields{
		"eventId":      requestEvent.NostrId,
		"appId":        app.ID,
		"senderPubkey": payParams.Pubkey,
	}).Info("Sending payment")

	preimage, err := svc.lnClient.SendKeysend(ctx, payParams.Amount, payParams.Pubkey, payParams.Preimage, payParams.TLVRecords)
	if err != nil {
		svc.Logger.WithFields(logrus.Fields{
			"eventId":         requestEvent.NostrId,
			"appId":           app.ID,
			"recipientPubkey": payParams.Pubkey,
		}).Infof("Failed to send payment: %v", err)
		svc.EventLogger.Log(&events.Event{
			Event: "nwc_payment_failed",
			Properties: map[string]interface{}{
				// "error":   fmt.Sprintf("%v", err),
				"keysend": true,
				"amount":  payParams.Amount / 1000,
			},
		})
		return &Nip47Response{
			ResultType: request.Method,
			Error: &Nip47Error{
				Code:    NIP_47_ERROR_INTERNAL,
				Message: err.Error(),
			},
		}, nil
	}
	payment.Preimage = &preimage
	svc.db.Save(&payment)
	svc.EventLogger.Log(&events.Event{
		Event: "nwc_payment_succeeded",
		Properties: map[string]interface{}{
			"keysend": true,
			"amount":  payParams.Amount / 1000,
		},
	})
	return &Nip47Response{
		ResultType: request.Method,
		Result: Nip47PayResponse{
			Preimage: preimage,
		},
	}, nil
}
