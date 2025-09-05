// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// WithdrawMetaData contains all meta data concerning the Withdraw contract.
var WithdrawMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getWithdrawalFeeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setWithdrawalFeeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawalFeeAmount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// WithdrawABI is the input ABI used to generate the binding from.
// Deprecated: Use WithdrawMetaData.ABI instead.
var WithdrawABI = WithdrawMetaData.ABI

// Withdraw is an auto generated Go binding around an Ethereum contract.
type Withdraw struct {
	WithdrawCaller     // Read-only binding to the contract
	WithdrawTransactor // Write-only binding to the contract
	WithdrawFilterer   // Log filterer for contract events
}

// WithdrawCaller is an auto generated read-only Go binding around an Ethereum contract.
type WithdrawCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WithdrawTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WithdrawFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WithdrawSession struct {
	Contract     *Withdraw         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WithdrawCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WithdrawCallerSession struct {
	Contract *WithdrawCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// WithdrawTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WithdrawTransactorSession struct {
	Contract     *WithdrawTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// WithdrawRaw is an auto generated low-level Go binding around an Ethereum contract.
type WithdrawRaw struct {
	Contract *Withdraw // Generic contract binding to access the raw methods on
}

// WithdrawCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WithdrawCallerRaw struct {
	Contract *WithdrawCaller // Generic read-only contract binding to access the raw methods on
}

// WithdrawTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WithdrawTransactorRaw struct {
	Contract *WithdrawTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWithdraw creates a new instance of Withdraw, bound to a specific deployed contract.
func NewWithdraw(address common.Address, backend bind.ContractBackend) (*Withdraw, error) {
	contract, err := bindWithdraw(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Withdraw{WithdrawCaller: WithdrawCaller{contract: contract}, WithdrawTransactor: WithdrawTransactor{contract: contract}, WithdrawFilterer: WithdrawFilterer{contract: contract}}, nil
}

// NewWithdrawCaller creates a new read-only instance of Withdraw, bound to a specific deployed contract.
func NewWithdrawCaller(address common.Address, caller bind.ContractCaller) (*WithdrawCaller, error) {
	contract, err := bindWithdraw(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WithdrawCaller{contract: contract}, nil
}

// NewWithdrawTransactor creates a new write-only instance of Withdraw, bound to a specific deployed contract.
func NewWithdrawTransactor(address common.Address, transactor bind.ContractTransactor) (*WithdrawTransactor, error) {
	contract, err := bindWithdraw(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WithdrawTransactor{contract: contract}, nil
}

// NewWithdrawFilterer creates a new log filterer instance of Withdraw, bound to a specific deployed contract.
func NewWithdrawFilterer(address common.Address, filterer bind.ContractFilterer) (*WithdrawFilterer, error) {
	contract, err := bindWithdraw(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WithdrawFilterer{contract: contract}, nil
}

// bindWithdraw binds a generic wrapper to an already deployed contract.
func bindWithdraw(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WithdrawMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Withdraw *WithdrawRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Withdraw.Contract.WithdrawCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Withdraw *WithdrawRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Withdraw.Contract.WithdrawTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Withdraw *WithdrawRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Withdraw.Contract.WithdrawTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Withdraw *WithdrawCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Withdraw.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Withdraw *WithdrawTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Withdraw.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Withdraw *WithdrawTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Withdraw.Contract.contract.Transact(opts, method, params...)
}

// GetWithdrawalFeeAddress is a free data retrieval call binding the contract method 0x0cf22d9f.
//
// Solidity: function getWithdrawalFeeAddress() view returns(address)
func (_Withdraw *WithdrawCaller) GetWithdrawalFeeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Withdraw.contract.Call(opts, &out, "getWithdrawalFeeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetWithdrawalFeeAddress is a free data retrieval call binding the contract method 0x0cf22d9f.
//
// Solidity: function getWithdrawalFeeAddress() view returns(address)
func (_Withdraw *WithdrawSession) GetWithdrawalFeeAddress() (common.Address, error) {
	return _Withdraw.Contract.GetWithdrawalFeeAddress(&_Withdraw.CallOpts)
}

// GetWithdrawalFeeAddress is a free data retrieval call binding the contract method 0x0cf22d9f.
//
// Solidity: function getWithdrawalFeeAddress() view returns(address)
func (_Withdraw *WithdrawCallerSession) GetWithdrawalFeeAddress() (common.Address, error) {
	return _Withdraw.Contract.GetWithdrawalFeeAddress(&_Withdraw.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Withdraw *WithdrawCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Withdraw.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Withdraw *WithdrawSession) Owner() (common.Address, error) {
	return _Withdraw.Contract.Owner(&_Withdraw.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Withdraw *WithdrawCallerSession) Owner() (common.Address, error) {
	return _Withdraw.Contract.Owner(&_Withdraw.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Withdraw *WithdrawTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Withdraw.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Withdraw *WithdrawSession) RenounceOwnership() (*types.Transaction, error) {
	return _Withdraw.Contract.RenounceOwnership(&_Withdraw.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Withdraw *WithdrawTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Withdraw.Contract.RenounceOwnership(&_Withdraw.TransactOpts)
}

// SetWithdrawalFeeAddress is a paid mutator transaction binding the contract method 0x0e423e24.
//
// Solidity: function setWithdrawalFeeAddress(address addr) returns()
func (_Withdraw *WithdrawTransactor) SetWithdrawalFeeAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _Withdraw.contract.Transact(opts, "setWithdrawalFeeAddress", addr)
}

// SetWithdrawalFeeAddress is a paid mutator transaction binding the contract method 0x0e423e24.
//
// Solidity: function setWithdrawalFeeAddress(address addr) returns()
func (_Withdraw *WithdrawSession) SetWithdrawalFeeAddress(addr common.Address) (*types.Transaction, error) {
	return _Withdraw.Contract.SetWithdrawalFeeAddress(&_Withdraw.TransactOpts, addr)
}

// SetWithdrawalFeeAddress is a paid mutator transaction binding the contract method 0x0e423e24.
//
// Solidity: function setWithdrawalFeeAddress(address addr) returns()
func (_Withdraw *WithdrawTransactorSession) SetWithdrawalFeeAddress(addr common.Address) (*types.Transaction, error) {
	return _Withdraw.Contract.SetWithdrawalFeeAddress(&_Withdraw.TransactOpts, addr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Withdraw *WithdrawTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Withdraw.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Withdraw *WithdrawSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Withdraw.Contract.TransferOwnership(&_Withdraw.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Withdraw *WithdrawTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Withdraw.Contract.TransferOwnership(&_Withdraw.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0xb5c5f672.
//
// Solidity: function withdraw(address to, uint256 amount, uint256 withdrawalFeeAmount) payable returns()
func (_Withdraw *WithdrawTransactor) Withdraw(opts *bind.TransactOpts, to common.Address, amount *big.Int, withdrawalFeeAmount *big.Int) (*types.Transaction, error) {
	return _Withdraw.contract.Transact(opts, "withdraw", to, amount, withdrawalFeeAmount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xb5c5f672.
//
// Solidity: function withdraw(address to, uint256 amount, uint256 withdrawalFeeAmount) payable returns()
func (_Withdraw *WithdrawSession) Withdraw(to common.Address, amount *big.Int, withdrawalFeeAmount *big.Int) (*types.Transaction, error) {
	return _Withdraw.Contract.Withdraw(&_Withdraw.TransactOpts, to, amount, withdrawalFeeAmount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xb5c5f672.
//
// Solidity: function withdraw(address to, uint256 amount, uint256 withdrawalFeeAmount) payable returns()
func (_Withdraw *WithdrawTransactorSession) Withdraw(to common.Address, amount *big.Int, withdrawalFeeAmount *big.Int) (*types.Transaction, error) {
	return _Withdraw.Contract.Withdraw(&_Withdraw.TransactOpts, to, amount, withdrawalFeeAmount)
}

// WithdrawOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Withdraw contract.
type WithdrawOwnershipTransferredIterator struct {
	Event *WithdrawOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WithdrawOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WithdrawOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WithdrawOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WithdrawOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WithdrawOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WithdrawOwnershipTransferred represents a OwnershipTransferred event raised by the Withdraw contract.
type WithdrawOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Withdraw *WithdrawFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WithdrawOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Withdraw.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WithdrawOwnershipTransferredIterator{contract: _Withdraw.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Withdraw *WithdrawFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WithdrawOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Withdraw.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WithdrawOwnershipTransferred)
				if err := _Withdraw.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Withdraw *WithdrawFilterer) ParseOwnershipTransferred(log types.Log) (*WithdrawOwnershipTransferred, error) {
	event := new(WithdrawOwnershipTransferred)
	if err := _Withdraw.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
