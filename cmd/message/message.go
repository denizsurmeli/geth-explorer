package message

import "github.com/fatih/color"

//General messages related to the operational errors/mismatches. Observe that
//these messages are not taking any parameters, further refactors might change
//this situation.

// TODO:Add logging.
// Setup messages:
// Success -> Color:Green
// No panics -> Color:Yellow
// Panics -> Color:Red

func SetupFailEnvironmentMessage() {
	color.Red("[SETUP] Error while setting up the environment.")
}

// Network messages:
// Success -> Color:Green
// No panics -> Color:Yellow
// Panics -> Color:Red

func NetworkFailToDialMessage() {
	color.Red("[NETWORK] Error while connecting to the node.")
}
func NetworkFailWhileRequestMessage() {
	color.Red("[NETWORK] Error while request to geth.")
}
func NetworkFailToReadMessage() {
	color.Yellow("[NETWORK] Could not read the message. Session might be terminated, or provider could not supply the correct response.")
}

func NetworkFailToGetNetworkId() {
	color.Red("[NETWORK] Could not fetch the Network ID.")
}

func NetworkInterruptOnConnectionMessage() {
	color.Red("[NETWORK] Interruption on connection.")
}

func NetworkConnectionLeftOpenMessage() {
	color.Yellow("[NETWORK] Connection left open. It might cause further errors.")
}

func NetworkFailToSendMessage() {
	color.Yellow("[NETWORK] Could not send the message to the connection.")
}

func NetworkNotSupportedMessage() {
	color.Red("[NETWORK] The type of network is not supported. Use Websockets or supported companies.")
}
func NetworkConnectionSuccessfulMessage() {
	color.Green("[NETWORK] Node connection established.")
}

func NetworkDialNodeMessage() {
	color.Green("[NETWORK] Dialing node.")
}

// User messages:
// Success -> Color:Green
// No panics -> Color:Yellow
// Panics -> Color:Red

func UserFalseParameterMessage() {
	color.Yellow("[USER] Parameter is not correct. Retry with correct parameters.")
}

// Operation messages:
// Success -> Color:Green
// No panics -> Color:Yellow
// Panics -> Color:Red

func OperationErrorMessage() {
	color.Red("[OPERATION] Critical error occurred while performing the operation.")
}

func OperationFailTransaction() {
	color.Red("[OPERATION] Error occured while fetching the transaction.")
}

// Generalized Messages:
// Success -> Color:Green
// No panics -> Color:Yellow
// Panics -> Color:Red

func NonSupportedOperation() {
	color.Red("[GENERAL] This operation is unknown or not supported.")
}
