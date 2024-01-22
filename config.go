package main

import "github.com/getAlby/nostr-wallet-connect/models/db"

const (
	LNDBackendType   = "LND"
	BreezBackendType = "BREEZ"
	CookieName       = "alby_nwc_session"
)

type Config struct {
	// These config can be set either by .env or the database config table.
	// database config always takes preference.
	db.Config
	CookieSecret            string `envconfig:"COOKIE_SECRET"`
	CookieDomain            string `envconfig:"COOKIE_DOMAIN"`
	ClientPubkey            string `envconfig:"CLIENT_NOSTR_PUBKEY"`
	Relay                   string `envconfig:"RELAY" default:"wss://relay.getalby.com/v1"`
	PublicRelay             string `envconfig:"PUBLIC_RELAY"`
	LNDCertFile             string `envconfig:"LND_CERT_FILE"`
	LNDMacaroonFile         string `envconfig:"LND_MACAROON_FILE"`
	BreezWorkdir            string `envconfig:"BREEZ_WORK_DIR" default:".breez"`
	BasicAuthUser           string `envconfig:"BASIC_AUTH_USER"`
	BasicAuthPassword       string `envconfig:"BASIC_AUTH_PASSWORD"`
	Port                    string `envconfig:"PORT" default:"8080"`
	DatabaseUri             string `envconfig:"DATABASE_URI" default:"nostr-wallet-connect.db"`
	DatabaseMaxConns        int    `envconfig:"DATABASE_MAX_CONNS" default:"10"`
	DatabaseMaxIdleConns    int    `envconfig:"DATABASE_MAX_IDLE_CONNS" default:"5"`
	DatabaseConnMaxLifetime int    `envconfig:"DATABASE_CONN_MAX_LIFETIME" default:"1800"` // 30 minutes
	IdentityPubkey          string
}
