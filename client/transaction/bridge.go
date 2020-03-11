package transaction

import (
	sdk "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type TransferInResult struct {
	tx.TxCommitResult
}

func (c *client) TransferIn(sequence int64, contractAddr msg.EthereumAddress,
	senderAddr msg.EthereumAddress, receiverAddr sdk.AccAddress, amount sdk.Coin,
	relayFee sdk.Coin, expireTime int64, sync bool, options ...Option) (*TransferInResult, error) {
	fromAddr := c.keyManager.GetAddr()
	transferInMsg := msg.NewTransferInMsg(sequence, contractAddr, senderAddr, receiverAddr, amount,
		relayFee, fromAddr, expireTime)
	commit, err := c.broadcastMsg(transferInMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &TransferInResult{*commit}, nil
}

type BindResult struct {
	tx.TxCommitResult
}

func (c *client) Bind(symbol string, amount int64, contractAddress msg.EthereumAddress, contractDecimal int, sync bool, options ...Option) (*BindResult, error) {
	fromAddr := c.keyManager.GetAddr()
	bindMsg := msg.NewBindMsg(fromAddr, symbol, amount, contractAddress, contractDecimal)
	commit, err := c.broadcastMsg(bindMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &BindResult{*commit}, nil
}

type TransferOutResult struct {
	tx.TxCommitResult
}

func (c *client) TransferOut(to msg.EthereumAddress, amount sdk.Coin, expireTime int64, sync bool, options ...Option) (*TransferOutResult, error) {
	fromAddr := c.keyManager.GetAddr()
	transferOutMsg := msg.NewTransferOutMsg(fromAddr, to, amount, expireTime)
	commit, err := c.broadcastMsg(transferOutMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &TransferOutResult{*commit}, nil
}
