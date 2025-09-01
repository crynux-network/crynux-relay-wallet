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

// BenefitAddressMetaData contains all meta data concerning the BenefitAddress contract.
var BenefitAddressMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"benefitAddress\",\"type\":\"address\"}],\"name\":\"BenefitAddressSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"}],\"name\":\"getBenefitAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"benefitAddress\",\"type\":\"address\"}],\"name\":\"setBenefitAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BenefitAddressABI is the input ABI used to generate the binding from.
// Deprecated: Use BenefitAddressMetaData.ABI instead.
var BenefitAddressABI = BenefitAddressMetaData.ABI

// BenefitAddress is an auto generated Go binding around an Ethereum contract.
type BenefitAddress struct {
	BenefitAddressCaller     // Read-only binding to the contract
	BenefitAddressTransactor // Write-only binding to the contract
	BenefitAddressFilterer   // Log filterer for contract events
}

// BenefitAddressCaller is an auto generated read-only Go binding around an Ethereum contract.
type BenefitAddressCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BenefitAddressTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BenefitAddressTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BenefitAddressFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BenefitAddressFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BenefitAddressSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BenefitAddressSession struct {
	Contract     *BenefitAddress   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BenefitAddressCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BenefitAddressCallerSession struct {
	Contract *BenefitAddressCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// BenefitAddressTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BenefitAddressTransactorSession struct {
	Contract     *BenefitAddressTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BenefitAddressRaw is an auto generated low-level Go binding around an Ethereum contract.
type BenefitAddressRaw struct {
	Contract *BenefitAddress // Generic contract binding to access the raw methods on
}

// BenefitAddressCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BenefitAddressCallerRaw struct {
	Contract *BenefitAddressCaller // Generic read-only contract binding to access the raw methods on
}

// BenefitAddressTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BenefitAddressTransactorRaw struct {
	Contract *BenefitAddressTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBenefitAddress creates a new instance of BenefitAddress, bound to a specific deployed contract.
func NewBenefitAddress(address common.Address, backend bind.ContractBackend) (*BenefitAddress, error) {
	contract, err := bindBenefitAddress(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BenefitAddress{BenefitAddressCaller: BenefitAddressCaller{contract: contract}, BenefitAddressTransactor: BenefitAddressTransactor{contract: contract}, BenefitAddressFilterer: BenefitAddressFilterer{contract: contract}}, nil
}

// NewBenefitAddressCaller creates a new read-only instance of BenefitAddress, bound to a specific deployed contract.
func NewBenefitAddressCaller(address common.Address, caller bind.ContractCaller) (*BenefitAddressCaller, error) {
	contract, err := bindBenefitAddress(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BenefitAddressCaller{contract: contract}, nil
}

// NewBenefitAddressTransactor creates a new write-only instance of BenefitAddress, bound to a specific deployed contract.
func NewBenefitAddressTransactor(address common.Address, transactor bind.ContractTransactor) (*BenefitAddressTransactor, error) {
	contract, err := bindBenefitAddress(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BenefitAddressTransactor{contract: contract}, nil
}

// NewBenefitAddressFilterer creates a new log filterer instance of BenefitAddress, bound to a specific deployed contract.
func NewBenefitAddressFilterer(address common.Address, filterer bind.ContractFilterer) (*BenefitAddressFilterer, error) {
	contract, err := bindBenefitAddress(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BenefitAddressFilterer{contract: contract}, nil
}

// bindBenefitAddress binds a generic wrapper to an already deployed contract.
func bindBenefitAddress(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BenefitAddressMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BenefitAddress *BenefitAddressRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BenefitAddress.Contract.BenefitAddressCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BenefitAddress *BenefitAddressRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BenefitAddress.Contract.BenefitAddressTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BenefitAddress *BenefitAddressRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BenefitAddress.Contract.BenefitAddressTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BenefitAddress *BenefitAddressCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BenefitAddress.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BenefitAddress *BenefitAddressTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BenefitAddress.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BenefitAddress *BenefitAddressTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BenefitAddress.Contract.contract.Transact(opts, method, params...)
}

// GetBenefitAddress is a free data retrieval call binding the contract method 0x3fd46606.
//
// Solidity: function getBenefitAddress(address nodeAddress) view returns(address)
func (_BenefitAddress *BenefitAddressCaller) GetBenefitAddress(opts *bind.CallOpts, nodeAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _BenefitAddress.contract.Call(opts, &out, "getBenefitAddress", nodeAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBenefitAddress is a free data retrieval call binding the contract method 0x3fd46606.
//
// Solidity: function getBenefitAddress(address nodeAddress) view returns(address)
func (_BenefitAddress *BenefitAddressSession) GetBenefitAddress(nodeAddress common.Address) (common.Address, error) {
	return _BenefitAddress.Contract.GetBenefitAddress(&_BenefitAddress.CallOpts, nodeAddress)
}

// GetBenefitAddress is a free data retrieval call binding the contract method 0x3fd46606.
//
// Solidity: function getBenefitAddress(address nodeAddress) view returns(address)
func (_BenefitAddress *BenefitAddressCallerSession) GetBenefitAddress(nodeAddress common.Address) (common.Address, error) {
	return _BenefitAddress.Contract.GetBenefitAddress(&_BenefitAddress.CallOpts, nodeAddress)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BenefitAddress *BenefitAddressCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BenefitAddress.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BenefitAddress *BenefitAddressSession) Owner() (common.Address, error) {
	return _BenefitAddress.Contract.Owner(&_BenefitAddress.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BenefitAddress *BenefitAddressCallerSession) Owner() (common.Address, error) {
	return _BenefitAddress.Contract.Owner(&_BenefitAddress.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BenefitAddress *BenefitAddressTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BenefitAddress.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BenefitAddress *BenefitAddressSession) RenounceOwnership() (*types.Transaction, error) {
	return _BenefitAddress.Contract.RenounceOwnership(&_BenefitAddress.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BenefitAddress *BenefitAddressTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BenefitAddress.Contract.RenounceOwnership(&_BenefitAddress.TransactOpts)
}

// SetBenefitAddress is a paid mutator transaction binding the contract method 0xfe94eb5c.
//
// Solidity: function setBenefitAddress(address benefitAddress) returns()
func (_BenefitAddress *BenefitAddressTransactor) SetBenefitAddress(opts *bind.TransactOpts, benefitAddress common.Address) (*types.Transaction, error) {
	return _BenefitAddress.contract.Transact(opts, "setBenefitAddress", benefitAddress)
}

// SetBenefitAddress is a paid mutator transaction binding the contract method 0xfe94eb5c.
//
// Solidity: function setBenefitAddress(address benefitAddress) returns()
func (_BenefitAddress *BenefitAddressSession) SetBenefitAddress(benefitAddress common.Address) (*types.Transaction, error) {
	return _BenefitAddress.Contract.SetBenefitAddress(&_BenefitAddress.TransactOpts, benefitAddress)
}

// SetBenefitAddress is a paid mutator transaction binding the contract method 0xfe94eb5c.
//
// Solidity: function setBenefitAddress(address benefitAddress) returns()
func (_BenefitAddress *BenefitAddressTransactorSession) SetBenefitAddress(benefitAddress common.Address) (*types.Transaction, error) {
	return _BenefitAddress.Contract.SetBenefitAddress(&_BenefitAddress.TransactOpts, benefitAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BenefitAddress *BenefitAddressTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BenefitAddress.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BenefitAddress *BenefitAddressSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BenefitAddress.Contract.TransferOwnership(&_BenefitAddress.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BenefitAddress *BenefitAddressTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BenefitAddress.Contract.TransferOwnership(&_BenefitAddress.TransactOpts, newOwner)
}

// BenefitAddressBenefitAddressSetIterator is returned from FilterBenefitAddressSet and is used to iterate over the raw logs and unpacked data for BenefitAddressSet events raised by the BenefitAddress contract.
type BenefitAddressBenefitAddressSetIterator struct {
	Event *BenefitAddressBenefitAddressSet // Event containing the contract specifics and raw log

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
func (it *BenefitAddressBenefitAddressSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BenefitAddressBenefitAddressSet)
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
		it.Event = new(BenefitAddressBenefitAddressSet)
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
func (it *BenefitAddressBenefitAddressSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BenefitAddressBenefitAddressSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BenefitAddressBenefitAddressSet represents a BenefitAddressSet event raised by the BenefitAddress contract.
type BenefitAddressBenefitAddressSet struct {
	NodeAddress    common.Address
	BenefitAddress common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBenefitAddressSet is a free log retrieval operation binding the contract event 0x46d52641dcc3ae305284d3c23db6464d02f7dc1bfffb66d7b2042b43bef7873f.
//
// Solidity: event BenefitAddressSet(address indexed nodeAddress, address indexed benefitAddress)
func (_BenefitAddress *BenefitAddressFilterer) FilterBenefitAddressSet(opts *bind.FilterOpts, nodeAddress []common.Address, benefitAddress []common.Address) (*BenefitAddressBenefitAddressSetIterator, error) {

	var nodeAddressRule []interface{}
	for _, nodeAddressItem := range nodeAddress {
		nodeAddressRule = append(nodeAddressRule, nodeAddressItem)
	}
	var benefitAddressRule []interface{}
	for _, benefitAddressItem := range benefitAddress {
		benefitAddressRule = append(benefitAddressRule, benefitAddressItem)
	}

	logs, sub, err := _BenefitAddress.contract.FilterLogs(opts, "BenefitAddressSet", nodeAddressRule, benefitAddressRule)
	if err != nil {
		return nil, err
	}
	return &BenefitAddressBenefitAddressSetIterator{contract: _BenefitAddress.contract, event: "BenefitAddressSet", logs: logs, sub: sub}, nil
}

// WatchBenefitAddressSet is a free log subscription operation binding the contract event 0x46d52641dcc3ae305284d3c23db6464d02f7dc1bfffb66d7b2042b43bef7873f.
//
// Solidity: event BenefitAddressSet(address indexed nodeAddress, address indexed benefitAddress)
func (_BenefitAddress *BenefitAddressFilterer) WatchBenefitAddressSet(opts *bind.WatchOpts, sink chan<- *BenefitAddressBenefitAddressSet, nodeAddress []common.Address, benefitAddress []common.Address) (event.Subscription, error) {

	var nodeAddressRule []interface{}
	for _, nodeAddressItem := range nodeAddress {
		nodeAddressRule = append(nodeAddressRule, nodeAddressItem)
	}
	var benefitAddressRule []interface{}
	for _, benefitAddressItem := range benefitAddress {
		benefitAddressRule = append(benefitAddressRule, benefitAddressItem)
	}

	logs, sub, err := _BenefitAddress.contract.WatchLogs(opts, "BenefitAddressSet", nodeAddressRule, benefitAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BenefitAddressBenefitAddressSet)
				if err := _BenefitAddress.contract.UnpackLog(event, "BenefitAddressSet", log); err != nil {
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

// ParseBenefitAddressSet is a log parse operation binding the contract event 0x46d52641dcc3ae305284d3c23db6464d02f7dc1bfffb66d7b2042b43bef7873f.
//
// Solidity: event BenefitAddressSet(address indexed nodeAddress, address indexed benefitAddress)
func (_BenefitAddress *BenefitAddressFilterer) ParseBenefitAddressSet(log types.Log) (*BenefitAddressBenefitAddressSet, error) {
	event := new(BenefitAddressBenefitAddressSet)
	if err := _BenefitAddress.contract.UnpackLog(event, "BenefitAddressSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BenefitAddressOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BenefitAddress contract.
type BenefitAddressOwnershipTransferredIterator struct {
	Event *BenefitAddressOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BenefitAddressOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BenefitAddressOwnershipTransferred)
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
		it.Event = new(BenefitAddressOwnershipTransferred)
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
func (it *BenefitAddressOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BenefitAddressOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BenefitAddressOwnershipTransferred represents a OwnershipTransferred event raised by the BenefitAddress contract.
type BenefitAddressOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BenefitAddress *BenefitAddressFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BenefitAddressOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BenefitAddress.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BenefitAddressOwnershipTransferredIterator{contract: _BenefitAddress.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BenefitAddress *BenefitAddressFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BenefitAddressOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BenefitAddress.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BenefitAddressOwnershipTransferred)
				if err := _BenefitAddress.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BenefitAddress *BenefitAddressFilterer) ParseOwnershipTransferred(log types.Log) (*BenefitAddressOwnershipTransferred, error) {
	event := new(BenefitAddressOwnershipTransferred)
	if err := _BenefitAddress.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
