package msg

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	sdk "github.com/binance-chain/go-sdk/common/types"
)

const (
	RouteBridge     = "bridge"
	MaxDecimal  int = 18

	TransferInMsgType         = "crossTransferIn"
	TransferOutTimeoutMsgType = "crossTransferOutTimeout"
	BindMsgType               = "crossBind"
	TransferOutMsgType        = "crossTransferOut"
	UpdateBindMsgType         = "crossUpdateBind"
)

// EthereumAddress defines a standard ethereum address
type EthereumAddress gethCommon.Address

// NewEthereumAddress is a constructor function for EthereumAddress
func NewEthereumAddress(address string) EthereumAddress {
	return EthereumAddress(gethCommon.HexToAddress(address))
}

func (ethAddr EthereumAddress) IsEmpty() bool {
	addrValue := big.NewInt(0)
	addrValue.SetBytes(ethAddr[:])

	return addrValue.Cmp(big.NewInt(0)) == 0
}

// Route should return the name of the module
func (ethAddr EthereumAddress) String() string {
	return gethCommon.Address(ethAddr).String()
}

// MarshalJSON marshals the ethereum address to JSON
func (ethAddr EthereumAddress) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", ethAddr.String())), nil
}

// UnmarshalJSON unmarshals an ethereum address
func (ethAddr *EthereumAddress) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(gethCommon.Address{}), input, ethAddr[:])
}

type TransferInMsg struct {
	Sequence         int64           `json:"sequence"`
	ContractAddress  EthereumAddress `json:"contract_address"`
	SenderAddress    EthereumAddress `json:"sender_address"`
	ReceiverAddress  sdk.AccAddress  `json:"receiver_address"`
	Amount           sdk.Coin        `json:"amount"`
	RelayFee         sdk.Coin        `json:"relay_fee"`
	ValidatorAddress sdk.AccAddress  `json:"validator_address"`
	ExpireTime       int64           `json:"expire_time"`
}

func NewTransferInMsg(sequence int64, contractAddr EthereumAddress,
	senderAddr EthereumAddress, receiverAddr sdk.AccAddress, amount sdk.Coin,
	relayFee sdk.Coin, validatorAddr sdk.AccAddress, expireTime int64) TransferInMsg {
	return TransferInMsg{
		Sequence:         sequence,
		ContractAddress:  contractAddr,
		SenderAddress:    senderAddr,
		ReceiverAddress:  receiverAddr,
		Amount:           amount,
		RelayFee:         relayFee,
		ValidatorAddress: validatorAddr,
		ExpireTime:       expireTime,
	}
}

// nolint
func (msg TransferInMsg) Route() string { return RouteBridge }
func (msg TransferInMsg) Type() string  { return TransferInMsgType }
func (msg TransferInMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.ValidatorAddress}
}

func (msg TransferInMsg) String() string {
	return fmt.Sprintf("TransferIn{%v#%s#%s#%s#%s#%s#%s#%d}",
		msg.ValidatorAddress, msg.ContractAddress.String(), msg.SenderAddress.String(), msg.ReceiverAddress.String(),
		msg.Amount.String(), msg.RelayFee.String(), msg.ValidatorAddress.String(), msg.ExpireTime)
}

// GetSignBytes - Get the bytes for the message signer to sign on
func (msg TransferInMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg TransferInMsg) GetInvolvedAddresses() []sdk.AccAddress {
	return msg.GetSigners()
}

// ValidateBasic is used to quickly disqualify obviously invalid messages quickly
func (msg TransferInMsg) ValidateBasic() error {
	if msg.Sequence < 0 {
		return fmt.Errorf("sequence should not be less than 0")
	}
	if msg.ExpireTime <= 0 {
		return fmt.Errorf("expire time should be larger than 0")
	}
	if msg.ContractAddress.IsEmpty() {
		return fmt.Errorf("contract address should not be empty")
	}
	if msg.SenderAddress.IsEmpty() {
		return fmt.Errorf("sender address should not be empty")
	}
	if len(msg.ReceiverAddress) != sdk.AddrLen {
		return fmt.Errorf("lenghth of receiver address should be %d", sdk.AddrLen)
	}
	if len(msg.ValidatorAddress) != sdk.AddrLen {
		return fmt.Errorf("lenghth of validator address should be %d", sdk.AddrLen)
	}
	if !msg.Amount.IsPositive() {
		return fmt.Errorf("amount to send should be positive")
	}
	if !msg.RelayFee.IsPositive() {
		return fmt.Errorf("amount to send should be positive")
	}
	return nil
}

type TransferOutTimeoutMsg struct {
	SenderAddress    sdk.AccAddress `json:"sender_address"`
	Sequence         int64          `json:"sequence"`
	Amount           sdk.Coin       `json:"amount"`
	ValidatorAddress sdk.AccAddress `json:"validator_address"`
}

func NewTransferOutTimeoutMsg(senderAddr sdk.AccAddress, sequence int64, amount sdk.Coin, validatorAddr sdk.AccAddress) TransferOutTimeoutMsg {
	return TransferOutTimeoutMsg{
		SenderAddress:    senderAddr,
		Sequence:         sequence,
		Amount:           amount,
		ValidatorAddress: validatorAddr,
	}
}

// nolint
func (msg TransferOutTimeoutMsg) Route() string { return RouteBridge }
func (msg TransferOutTimeoutMsg) Type() string  { return TransferOutTimeoutMsgType }
func (msg TransferOutTimeoutMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.ValidatorAddress}
}
func (msg TransferOutTimeoutMsg) String() string {
	return fmt.Sprintf("TransferOutTimeout{%s#%d#%s#%s}",
		msg.SenderAddress.String(), msg.Sequence, msg.Amount.String(), msg.ValidatorAddress.String())
}

// GetSignBytes - Get the bytes for the message signer to sign on
func (msg TransferOutTimeoutMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg TransferOutTimeoutMsg) GetInvolvedAddresses() []sdk.AccAddress {
	return msg.GetSigners()
}

// ValidateBasic is used to quickly disqualify obviously invalid messages quickly
func (msg TransferOutTimeoutMsg) ValidateBasic() error {
	if len(msg.SenderAddress) != sdk.AddrLen {
		return fmt.Errorf("lenghth of sender address should be %d", sdk.AddrLen)
	}
	if msg.Sequence < 0 {
		return fmt.Errorf("sequence should not be less than 0")
	}
	if len(msg.ValidatorAddress) != sdk.AddrLen {
		return fmt.Errorf("lenghth of validator address should be %d", sdk.AddrLen)
	}
	if !msg.Amount.IsPositive() {
		return fmt.Errorf("amount to send should be positive")
	}
	return nil
}

type BindMsg struct {
	From             sdk.AccAddress  `json:"from"`
	Symbol           string          `json:"symbol"`
	Amount           int64           `json:"amount"`
	ContractAddress  EthereumAddress `json:"contract_address"`
	ContractDecimals int8            `json:"contract_decimals"`
	ExpireTime       int64           `json:"expire_time"`
}

func NewBindMsg(from sdk.AccAddress, symbol string, amount int64, contractAddress EthereumAddress, contractDecimals int8, expireTime int64) BindMsg {
	return BindMsg{
		From:             from,
		Amount:           amount,
		Symbol:           symbol,
		ContractAddress:  contractAddress,
		ContractDecimals: contractDecimals,
		ExpireTime:       expireTime,
	}
}

func (msg BindMsg) Route() string { return RouteBridge }
func (msg BindMsg) Type() string  { return BindMsgType }
func (msg BindMsg) String() string {
	return fmt.Sprintf("Bind{%v#%s#%d$%s#%d#%d}", msg.From, msg.Symbol, msg.Amount, msg.ContractAddress.String(), msg.ContractDecimals, msg.ExpireTime)
}
func (msg BindMsg) GetInvolvedAddresses() []sdk.AccAddress { return msg.GetSigners() }
func (msg BindMsg) GetSigners() []sdk.AccAddress           { return []sdk.AccAddress{msg.From} }

func (msg BindMsg) ValidateBasic() error {
	if len(msg.From) != sdk.AddrLen {
		return fmt.Errorf("address length should be %d", sdk.AddrLen)
	}

	if len(msg.Symbol) == 0 {
		return fmt.Errorf("symbol should not be empty")
	}

	if msg.Amount <= 0 {
		return fmt.Errorf("amount should be larger than 0")
	}

	if msg.ContractAddress.IsEmpty() {
		return fmt.Errorf("contract address should not be empty")
	}

	if msg.ContractDecimals < 0 {
		return fmt.Errorf("decimal should be no less than 0")
	}

	if msg.ExpireTime <= 0 {
		return fmt.Errorf("expire time should be larger than 0")
	}

	return nil
}

func (msg BindMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

type TransferOutMsg struct {
	From       sdk.AccAddress  `json:"from"`
	To         EthereumAddress `json:"to"`
	Amount     sdk.Coin        `json:"amount"`
	ExpireTime int64           `json:"expire_time"`
}

func NewTransferOutMsg(from sdk.AccAddress, to EthereumAddress, amount sdk.Coin, expireTime int64) TransferOutMsg {
	return TransferOutMsg{
		From:       from,
		To:         to,
		Amount:     amount,
		ExpireTime: expireTime,
	}
}

func (msg TransferOutMsg) Route() string { return RouteBridge }
func (msg TransferOutMsg) Type() string  { return TransferOutMsgType }
func (msg TransferOutMsg) String() string {
	return fmt.Sprintf("TransferOut{%v#%s#%s#%d}", msg.From, msg.To.String(), msg.Amount.String(), msg.ExpireTime)
}
func (msg TransferOutMsg) GetInvolvedAddresses() []sdk.AccAddress { return msg.GetSigners() }
func (msg TransferOutMsg) GetSigners() []sdk.AccAddress           { return []sdk.AccAddress{msg.From} }
func (msg TransferOutMsg) ValidateBasic() error {
	if len(msg.From) != sdk.AddrLen {
		return fmt.Errorf("address length should be %d", sdk.AddrLen)
	}

	if msg.To.IsEmpty() {
		return fmt.Errorf("to address should not be empty")
	}

	if !msg.Amount.IsPositive() {
		return fmt.Errorf("amount should be positive")
	}

	if msg.ExpireTime <= 0 {
		return fmt.Errorf("expire time should be larger than 0")
	}

	return nil
}
func (msg TransferOutMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

type BindStatus int8

const (
	BindStatusSuccess          BindStatus = 0
	BindStatusRejected         BindStatus = 1
	BindStatusTimeout          BindStatus = 2
	BindStatusInvalidParameter BindStatus = 3
)

type UpdateBindMsg struct {
	Sequence         int64           `json:"sequence"`
	Status           BindStatus      `json:"status"`
	Symbol           string          `json:"symbol"`
	Amount           int64           `json:"amount"`
	ContractAddress  EthereumAddress `json:"contract_address"`
	ContractDecimals int8            `json:"contract_decimals"`
	ValidatorAddress sdk.AccAddress  `json:"validator_address"`
}

func NewUpdateBindMsg(sequence int64, validatorAddress sdk.AccAddress, symbol string, amount int64, contractAddress EthereumAddress, contractDecimals int8, status BindStatus) UpdateBindMsg {
	return UpdateBindMsg{
		Sequence:         sequence,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
		Symbol:           symbol,
		ContractAddress:  contractAddress,
		ContractDecimals: contractDecimals,
		Status:           status,
	}
}

func (msg UpdateBindMsg) Route() string { return RouteBridge }
func (msg UpdateBindMsg) Type() string  { return UpdateBindMsgType }
func (msg UpdateBindMsg) String() string {
	return fmt.Sprintf("UpdateBind{%v#%s#%d$%s#%d#%d}", msg.ValidatorAddress, msg.Symbol, msg.Amount, msg.ContractAddress.String(), msg.ContractDecimals, msg.Status)
}
func (msg UpdateBindMsg) GetInvolvedAddresses() []sdk.AccAddress { return msg.GetSigners() }
func (msg UpdateBindMsg) GetSigners() []sdk.AccAddress           { return []sdk.AccAddress{msg.ValidatorAddress} }

func (msg UpdateBindMsg) ValidateBasic() error {
	if len(msg.ValidatorAddress) != sdk.AddrLen {
		return fmt.Errorf("address length should be %d", sdk.AddrLen)
	}

	if len(msg.Symbol) == 0 {
		return fmt.Errorf("symbol should not be empty")
	}

	if msg.Amount <= 0 {
		return fmt.Errorf("amount should be larger than 0")
	}

	if msg.ContractAddress.IsEmpty() {
		return fmt.Errorf("contract address should not be empty")
	}

	if msg.ContractDecimals < 0 {
		return fmt.Errorf("decimal should be no less than 0")
	}

	return nil
}
func (msg UpdateBindMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}
