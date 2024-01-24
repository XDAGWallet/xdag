package wallet_state

import "goXdagWallet/i18n"

const (
	None        = 0  // Not an account
	Registering = 10 // Creating a new wallet and waiting for response

	Idle = 20 // Not connected to network

	LoadingBlocks = 30 // Loading blocks from the local storage.
	Stopped       = 40 // Blocks loaded. Waiting for 'run' command.

	ConnectingNetwork = 50 // Trying to connect to the network.
	ConnectedNetwork  = 55 // Connected to the network

	ConnectingPool = 60 // Trying to connect to the pool.

	ConnectedPool      = 65 // Connected to the pool. No mining.
	ConnectedAndMining = 67 // Connected to the pool. Mining on. Normal operation.

	Synchronizing = 70 // Synchronizing with the network
	Synchronized  = 75 // Synchronized with the network

	TransferPending = 80 // Waiting for transfer to complete.

	ResetingEngine = 90 // The local storage is corrupted. Resetting blocks engine.
)

var stateMessageMap = map[string]int{

	"Generating keys...": Registering,
	"Synchronized with the main network. Normal operation.":       Synchronized,
	"Synchronized with the test network. Normal testing.":         Synchronized,
	"Connected to the mainnet pool. Mining on. Normal operation.": ConnectedAndMining,
	"Connected to the testnet pool. Mining on. Normal testing.":   ConnectedAndMining,
	"Connected to the mainnet pool. No mining.":                   ConnectedPool,
	"Connected to the testnet pool. No mining.":                   ConnectedPool,
	"Waiting for transfer to complete.":                           TransferPending,
	"Connected to the main network. Synchronizing.":               Synchronizing,
	"Connected to the test network. Synchronizing.":               Synchronizing,
	"Trying to connect to the mainnet pool.":                      ConnectingPool,
	"Trying to connect to the testnet pool.":                      ConnectingPool,
	"Trying to connect to the main network.":                      ConnectingNetwork,
	"Trying to connect to the test network.":                      ConnectingNetwork,
	"Blocks loaded. Waiting for 'run' command.":                   Idle,
	"Loading blocks from the local storage.":                      LoadingBlocks,
	"The local storage is corrupted. Resetting blocks engine.":    ResetingEngine,
}

func MessageToState(msg string) (int, bool) {
	state, ok := stateMessageMap[msg]
	return state, ok
}

func Localize(state int) string {
	switch state {
	case None:
		return i18n.GetString("WalletState_None")
	case Registering:
		return i18n.GetString("WalletState_Registering")
	case Idle:
		return i18n.GetString("WalletState_Idle")
	case LoadingBlocks:
		return i18n.GetString("WalletState_LoadingBlocks")
	case Stopped:
		return i18n.GetString("WalletState_Stopped")
	case ConnectingNetwork:
		return i18n.GetString("WalletState_ConnectingNetwork")
	case ConnectedNetwork:
		return i18n.GetString("WalletState_ConnectedNetwork")
	case ConnectingPool:
		return i18n.GetString("WalletState_ConnectingPool")
	case ConnectedPool:
		return i18n.GetString("WalletState_ConnectedPool")
	case ConnectedAndMining:
		return i18n.GetString("WalletState_ConnectedAndMining")
	case Synchronizing:
		return i18n.GetString("WalletState_Synchronizing")
	case Synchronized:
		return i18n.GetString("WalletState_Synchronized")
	case TransferPending:
		return i18n.GetString("WalletState_TransferPending")

	case ResetingEngine:
		return i18n.GetString("WalletState_ResetingEngine")
	default:
		return ""
	}
}
