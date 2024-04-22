// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package faucet

import (
	"net/http"

	"github.com/AssetMantle/modules/utilities/rest/queuing"
	"github.com/cosmos/cosmos-sdk/client"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gorilla/mux"
)

func RegisterRESTRoutes(context client.Context, router *mux.Router) {
	handler := func(context client.Context) http.HandlerFunc {
		return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
			if !faucetEnabled {
				rest.WriteErrorResponse(responseWriter, http.StatusForbidden, "faucet is disabled")
				return
			}

			toAddress, err := sdkTypes.AccAddressFromBech32(mux.Vars(httpRequest)[toAddress])
			if rest.CheckBadRequestError(responseWriter, err) {
				return
			}

			fromKeyInfo, err := context.Keyring.Key(faucetKeyName)
			if rest.CheckBadRequestError(responseWriter, err) {
				return
			}

			context = context.WithFromName(fromKeyInfo.GetName()).WithFromAddress(fromKeyInfo.GetAddress()).WithChainID(chainID).WithSkipConfirmation(true)

			if rest.CheckInternalServerError(responseWriter, queuing.QueueOrBroadcastTransaction(context.WithOutput(responseWriter), rest.NewBaseReq(faucetKeyName, "faucet", context.ChainID, gas, gasAdjustment, 0, 0, sdkTypes.NewCoins(), sdkTypes.NewDecCoins(), false), types.NewMsgSend(context.FromAddress, toAddress, sdkTypes.NewCoins(sdkTypes.NewCoin(faucetDenom, sdkTypes.NewInt(faucetAmount)))))) {
				return
			}

		}
	}

	router.HandleFunc("/faucet/{"+toAddress+"}", handler(context)).Methods("GET")
}
