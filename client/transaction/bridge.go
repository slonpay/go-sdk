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
	refundAddresses []msg.EthereumAddress, receiverAddresses []sdk.AccAddress, amounts []int64, symbol string,
	relayFee sdk.Coin, expireTime int64, sync bool, options ...Option) (*TransferInResult, error) {
	fromAddr := c.keyManager.GetAddr()
	transferInMsg := msg.NewTransferInMsg(sequence, contractAddr, refundAddresses, receiverAddresses, amounts, symbol,
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

func (c *client) Bind(symbol string, amount int64, contractAddress msg.EthereumAddress, contractDecimals int8, expireTime int64, sync bool, options ...Option) (*BindResult, error) {
	fromAddr := c.keyManager.GetAddr()
	bindMsg := msg.NewBindMsg(fromAddr, symbol, amount, contractAddress, contractDecimals, expireTime)
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

type UpdateTransferOutResult struct {
	tx.TxCommitResult
}

func (c *client) UpdateTransferOut(sequence int64, sender sdk.AccAddress, amount sdk.Coin, status msg.TransferOutStatus, sync bool, options ...Option) (*UpdateTransferOutResult, error) {
	fromAddr := c.keyManager.GetAddr()
	transferOutTimeOutMsg := msg.NewUpdateTransferOutMsg(sender, sequence, amount, fromAddr, status)
	commit, err := c.broadcastMsg(transferOutTimeOutMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &UpdateTransferOutResult{*commit}, nil
}

type UpdateBindResult struct {
	tx.TxCommitResult
}

func (c *client) UpdateBind(sequence int64, symbol string, amount sdk.Int, contractAddress msg.EthereumAddress, contractDecimals int8, status msg.BindStatus, sync bool, options ...Option) (*UpdateBindResult, error) {
	fromAddr := c.keyManager.GetAddr()
	updateBindMsg := msg.NewUpdateBindMsg(sequence, fromAddr, symbol, amount, contractAddress, contractDecimals, status)
	commit, err := c.broadcastMsg(updateBindMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &UpdateBindResult{*commit}, nil
}
