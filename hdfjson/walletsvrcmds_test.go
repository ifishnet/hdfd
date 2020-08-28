// Copyright (c) 2014 The ifishnet developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hdfjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/ifishnet/hdfd/hdfjson"
)

// TestWalletSvrCmds tests all of the wallet server commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestWalletSvrCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return hdfjson.NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &hdfjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return hdfjson.NewAddMultisigAddressCmd(2, keys, hdfjson.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &hdfjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   hdfjson.String("test"),
			},
		},
		{
			name: "addwitnessaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("addwitnessaddress", "1address")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewAddWitnessAddressCmd("1address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"addwitnessaddress","params":["1address"],"id":1}`,
			unmarshalled: &hdfjson.AddWitnessAddressCmd{
				Address: "1address",
			},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return hdfjson.NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &hdfjson.CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &hdfjson.DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "encryptwallet",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("encryptwallet", "pass")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewEncryptWalletCmd("pass")
			},
			marshalled: `{"jsonrpc":"1.0","method":"encryptwallet","params":["pass"],"id":1}`,
			unmarshalled: &hdfjson.EncryptWalletCmd{
				Passphrase: "pass",
			},
		},
		{
			name: "estimatefee",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("estimatefee", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewEstimateFeeCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatefee","params":[6],"id":1}`,
			unmarshalled: &hdfjson.EstimateFeeCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "estimatesmartfee - no mode",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("estimatesmartfee", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewEstimateSmartFeeCmd(6, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatesmartfee","params":[6],"id":1}`,
			unmarshalled: &hdfjson.EstimateSmartFeeCmd{
				ConfTarget:   6,
				EstimateMode: &hdfjson.EstimateModeConservative,
			},
		},
		{
			name: "estimatesmartfee - economical mode",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("estimatesmartfee", 6, hdfjson.EstimateModeEconomical)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewEstimateSmartFeeCmd(6, &hdfjson.EstimateModeEconomical)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatesmartfee","params":[6,"ECONOMICAL"],"id":1}`,
			unmarshalled: &hdfjson.EstimateSmartFeeCmd{
				ConfTarget:   6,
				EstimateMode: &hdfjson.EstimateModeEconomical,
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &hdfjson.EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &hdfjson.GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &hdfjson.GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &hdfjson.GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetBalanceCmd{
				Account: nil,
				MinConf: hdfjson.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBalanceCmd(hdfjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &hdfjson.GetBalanceCmd{
				Account: hdfjson.String("acct"),
				MinConf: hdfjson.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBalanceCmd(hdfjson.String("acct"), hdfjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &hdfjson.GetBalanceCmd{
				Account: hdfjson.String("acct"),
				MinConf: hdfjson.Int(6),
			},
		},
		{
			name: "getbalances",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getbalances")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBalancesCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbalances","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetBalancesCmd{},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetNewAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetNewAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getnewaddress", "acct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetNewAddressCmd(hdfjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct"],"id":1}`,
			unmarshalled: &hdfjson.GetNewAddressCmd{
				Account: hdfjson.String("acct"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetRawChangeAddressCmd(hdfjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &hdfjson.GetRawChangeAddressCmd{
				Account: hdfjson.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &hdfjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: hdfjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetReceivedByAccountCmd("acct", hdfjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &hdfjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: hdfjson.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &hdfjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: hdfjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetReceivedByAddressCmd("1Address", hdfjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &hdfjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: hdfjson.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &hdfjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetTransactionCmd("123", hdfjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &hdfjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: hdfjson.Bool(true),
			},
		},
		{
			name: "getwalletinfo",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getwalletinfo")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetWalletInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getwalletinfo","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetWalletInfoCmd{},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewImportPrivKeyCmd("abc", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &hdfjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  hdfjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewImportPrivKeyCmd("abc", hdfjson.String("label"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &hdfjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   hdfjson.String("label"),
				Rescan:  hdfjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewImportPrivKeyCmd("abc", hdfjson.String("label"), hdfjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &hdfjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   hdfjson.String("label"),
				Rescan:  hdfjson.Bool(false),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &hdfjson.KeyPoolRefillCmd{
				NewSize: hdfjson.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewKeyPoolRefillCmd(hdfjson.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &hdfjson.KeyPoolRefillCmd{
				NewSize: hdfjson.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &hdfjson.ListAccountsCmd{
				MinConf: hdfjson.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListAccountsCmd(hdfjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &hdfjson.ListAccountsCmd{
				MinConf: hdfjson.Int(6),
			},
		},
		{
			name: "listaddressgroupings",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listaddressgroupings")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListAddressGroupingsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listaddressgroupings","params":[],"id":1}`,
			unmarshalled: &hdfjson.ListAddressGroupingsCmd{},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &hdfjson.ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &hdfjson.ListReceivedByAccountCmd{
				MinConf:          hdfjson.Int(1),
				IncludeEmpty:     hdfjson.Bool(false),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListReceivedByAccountCmd(hdfjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &hdfjson.ListReceivedByAccountCmd{
				MinConf:          hdfjson.Int(6),
				IncludeEmpty:     hdfjson.Bool(false),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListReceivedByAccountCmd(hdfjson.Int(6), hdfjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &hdfjson.ListReceivedByAccountCmd{
				MinConf:          hdfjson.Int(6),
				IncludeEmpty:     hdfjson.Bool(true),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListReceivedByAccountCmd(hdfjson.Int(6), hdfjson.Bool(true), hdfjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &hdfjson.ListReceivedByAccountCmd{
				MinConf:          hdfjson.Int(6),
				IncludeEmpty:     hdfjson.Bool(true),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &hdfjson.ListReceivedByAddressCmd{
				MinConf:          hdfjson.Int(1),
				IncludeEmpty:     hdfjson.Bool(false),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListReceivedByAddressCmd(hdfjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &hdfjson.ListReceivedByAddressCmd{
				MinConf:          hdfjson.Int(6),
				IncludeEmpty:     hdfjson.Bool(false),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListReceivedByAddressCmd(hdfjson.Int(6), hdfjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &hdfjson.ListReceivedByAddressCmd{
				MinConf:          hdfjson.Int(6),
				IncludeEmpty:     hdfjson.Bool(true),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListReceivedByAddressCmd(hdfjson.Int(6), hdfjson.Bool(true), hdfjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &hdfjson.ListReceivedByAddressCmd{
				MinConf:          hdfjson.Int(6),
				IncludeEmpty:     hdfjson.Bool(true),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &hdfjson.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: hdfjson.Int(1),
				IncludeWatchOnly:    hdfjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListSinceBlockCmd(hdfjson.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &hdfjson.ListSinceBlockCmd{
				BlockHash:           hdfjson.String("123"),
				TargetConfirmations: hdfjson.Int(1),
				IncludeWatchOnly:    hdfjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListSinceBlockCmd(hdfjson.String("123"), hdfjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &hdfjson.ListSinceBlockCmd{
				BlockHash:           hdfjson.String("123"),
				TargetConfirmations: hdfjson.Int(6),
				IncludeWatchOnly:    hdfjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListSinceBlockCmd(hdfjson.String("123"), hdfjson.Int(6), hdfjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &hdfjson.ListSinceBlockCmd{
				BlockHash:           hdfjson.String("123"),
				TargetConfirmations: hdfjson.Int(6),
				IncludeWatchOnly:    hdfjson.Bool(true),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &hdfjson.ListTransactionsCmd{
				Account:          nil,
				Count:            hdfjson.Int(10),
				From:             hdfjson.Int(0),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListTransactionsCmd(hdfjson.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &hdfjson.ListTransactionsCmd{
				Account:          hdfjson.String("acct"),
				Count:            hdfjson.Int(10),
				From:             hdfjson.Int(0),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListTransactionsCmd(hdfjson.String("acct"), hdfjson.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &hdfjson.ListTransactionsCmd{
				Account:          hdfjson.String("acct"),
				Count:            hdfjson.Int(20),
				From:             hdfjson.Int(0),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListTransactionsCmd(hdfjson.String("acct"), hdfjson.Int(20),
					hdfjson.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &hdfjson.ListTransactionsCmd{
				Account:          hdfjson.String("acct"),
				Count:            hdfjson.Int(20),
				From:             hdfjson.Int(1),
				IncludeWatchOnly: hdfjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListTransactionsCmd(hdfjson.String("acct"), hdfjson.Int(20),
					hdfjson.Int(1), hdfjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &hdfjson.ListTransactionsCmd{
				Account:          hdfjson.String("acct"),
				Count:            hdfjson.Int(20),
				From:             hdfjson.Int(1),
				IncludeWatchOnly: hdfjson.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &hdfjson.ListUnspentCmd{
				MinConf:   hdfjson.Int(1),
				MaxConf:   hdfjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListUnspentCmd(hdfjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &hdfjson.ListUnspentCmd{
				MinConf:   hdfjson.Int(6),
				MaxConf:   hdfjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListUnspentCmd(hdfjson.Int(6), hdfjson.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &hdfjson.ListUnspentCmd{
				MinConf:   hdfjson.Int(6),
				MaxConf:   hdfjson.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {
				return hdfjson.NewListUnspentCmd(hdfjson.Int(6), hdfjson.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &hdfjson.ListUnspentCmd{
				MinConf:   hdfjson.Int(6),
				MaxConf:   hdfjson.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []hdfjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return hdfjson.NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1}]],"id":1}`,
			unmarshalled: &hdfjson.LockUnspentCmd{
				Unlock: true,
				Transactions: []hdfjson.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "move",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("move", "from", "to", 0.5)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewMoveCmd("from", "to", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5],"id":1}`,
			unmarshalled: &hdfjson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     hdfjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "move optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("move", "from", "to", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewMoveCmd("from", "to", 0.5, hdfjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6],"id":1}`,
			unmarshalled: &hdfjson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     hdfjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "move optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("move", "from", "to", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewMoveCmd("from", "to", 0.5, hdfjson.Int(6), hdfjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6,"comment"],"id":1}`,
			unmarshalled: &hdfjson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     hdfjson.Int(6),
				Comment:     hdfjson.String("comment"),
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &hdfjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hdfjson.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSendFromCmd("from", "1Address", 0.5, hdfjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &hdfjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hdfjson.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSendFromCmd("from", "1Address", 0.5, hdfjson.Int(6),
					hdfjson.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &hdfjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hdfjson.Int(6),
				Comment:     hdfjson.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSendFromCmd("from", "1Address", 0.5, hdfjson.Int(6),
					hdfjson.String("comment"), hdfjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &hdfjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hdfjson.Int(6),
				Comment:     hdfjson.String("comment"),
				CommentTo:   hdfjson.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return hdfjson.NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &hdfjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     hdfjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return hdfjson.NewSendManyCmd("from", amounts, hdfjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &hdfjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     hdfjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return hdfjson.NewSendManyCmd("from", amounts, hdfjson.Int(6), hdfjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &hdfjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     hdfjson.Int(6),
				Comment:     hdfjson.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &hdfjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSendToAddressCmd("1Address", 0.5, hdfjson.String("comment"),
					hdfjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &hdfjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   hdfjson.String("comment"),
				CommentTo: hdfjson.String("commentto"),
			},
		},
		{
			name: "setaccount",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("setaccount", "1Address", "acct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSetAccountCmd("1Address", "acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"setaccount","params":["1Address","acct"],"id":1}`,
			unmarshalled: &hdfjson.SetAccountCmd{
				Address: "1Address",
				Account: "acct",
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &hdfjson.SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &hdfjson.SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &hdfjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    hdfjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []hdfjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}

				return hdfjson.NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &hdfjson.SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]hdfjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    hdfjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {
				txInputs := []hdfjson.RawTxInput{}
				privKeys := []string{"abc"}
				return hdfjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &hdfjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]hdfjson.RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    hdfjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []hdfjson.RawTxInput{}
				privKeys := []string{}
				return hdfjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					hdfjson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &hdfjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]hdfjson.RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    hdfjson.String("ALL"),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &hdfjson.WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &hdfjson.WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &hdfjson.WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := hdfjson.MarshalCmd(testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = hdfjson.MarshalCmd(testID, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request hdfjson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = hdfjson.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}
