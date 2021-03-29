package wallet

type Expiration string

func (e Expiration) ToPtr() *Expiration {
	return &e
}

const (
	ExpirationExpired Expiration = "expired"
	ExpirationNever Expiration = "never"
	Expiration24h Expiration = "24h"

	ExpirationAuto Expiration = "auto" // introduced with 6.0 API
)

type AddressType string

func (at AddressType) ToPtr() *AddressType {
	return &at
}

// special address types (offline, max privacy, public offline) are described in this document:
// https://github.com/BeamMW/beam/wiki/Lelantus-CLI
// choice described here:
// https://github.com/BeamMW/beam/wiki/Beam-wallet-protocol-API#create_address
const (
	AddressTypeRegular AddressType = "regular" // old regular addresses, overrideable with UseNewRegularStyleOnly
	AddressTypeNewRegular AddressType = "regular_new" // 6.0 base58 based addresses
	AddressTypeMaxPrivacy AddressType = "max_privacy"
	AddressTypeOffline AddressType = "offline"
	AddressTypePublicOffline AddressType = "public_offline"
)
