package common

type Hash160 [20]byte

const XDAG_BIP44_CION_TYPE = 586

const (
	XDAG_BLOCK_FIELDS = 16
	XDAG_BLOCK_SIZE   = 512
	XDAG_FIELD_SIZE   = 32
	XDAG_HASH_SIZE    = 32
	XDAG_ADDRESS_SIZE = 32

	BIP32_WALLET_FOLDER    = "xdagj_wallet"
	BIP32_WALLET_FILE_NAME = "xdagj_wallet.bin"
)

type Hash [XDAG_HASH_SIZE]byte
type Field [XDAG_FIELD_SIZE]byte
type RawBlock [XDAG_BLOCK_SIZE]byte
type Signature [XDAG_FIELD_SIZE * 2]byte
