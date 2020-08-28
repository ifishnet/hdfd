// Copyright (c) 2014 The ifishnet developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hdfjson_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/ifishnet/hdfd/hdfjson"
	"github.com/ifishnet/hdfd/wire"
)

// TestChainSvrCmds tests all of the chain server commands marshal and unmarshal
// into valid results include handling of optional fields being omitted in the
// marshalled command, while optional fields with defaults have the default
// assigned on unmarshalled commands.
func TestChainSvrCmds(t *testing.T) {
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
			name: "addnode",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("addnode", "127.0.0.1", hdfjson.ANRemove)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewAddNodeCmd("127.0.0.1", hdfjson.ANRemove)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"addnode","params":["127.0.0.1","remove"],"id":1}`,
			unmarshalled: &hdfjson.AddNodeCmd{Addr: "127.0.0.1", SubCmd: hdfjson.ANRemove},
		},
		{
			name: "createrawtransaction",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				txInputs := []hdfjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return hdfjson.NewCreateRawTransactionCmd(txInputs, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123}],"id":1}`,
			unmarshalled: &hdfjson.CreateRawTransactionCmd{
				Inputs:  []hdfjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction - no inputs",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("createrawtransaction", `[]`, `{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"456": .0123}
				return hdfjson.NewCreateRawTransactionCmd(nil, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[],{"456":0.0123}],"id":1}`,
			unmarshalled: &hdfjson.CreateRawTransactionCmd{
				Inputs:  []hdfjson.TransactionInput{},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`, int64(12312333333))
			},
			staticCmd: func() interface{} {
				txInputs := []hdfjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return hdfjson.NewCreateRawTransactionCmd(txInputs, amounts, hdfjson.Int64(12312333333))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123},12312333333],"id":1}`,
			unmarshalled: &hdfjson.CreateRawTransactionCmd{
				Inputs:   []hdfjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts:  map[string]float64{"456": .0123},
				LockTime: hdfjson.Int64(12312333333),
			},
		},
		{
			name: "fundrawtransaction - empty opts",
			newCmd: func() (i interface{}, e error) {
				return hdfjson.NewCmd("fundrawtransaction", "deadbeef", "{}")
			},
			staticCmd: func() interface{} {
				deadbeef, err := hex.DecodeString("deadbeef")
				if err != nil {
					panic(err)
				}
				return hdfjson.NewFundRawTransactionCmd(deadbeef, hdfjson.FundRawTransactionOpts{}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"fundrawtransaction","params":["deadbeef",{}],"id":1}`,
			unmarshalled: &hdfjson.FundRawTransactionCmd{
				HexTx:     "deadbeef",
				Options:   hdfjson.FundRawTransactionOpts{},
				IsWitness: nil,
			},
		},
		{
			name: "fundrawtransaction - full opts",
			newCmd: func() (i interface{}, e error) {
				return hdfjson.NewCmd("fundrawtransaction", "deadbeef", `{"changeAddress":"bcrt1qeeuctq9wutlcl5zatge7rjgx0k45228cxez655","changePosition":1,"change_type":"legacy","includeWatching":true,"lockUnspents":true,"feeRate":0.7,"subtractFeeFromOutputs":[0],"replaceable":true,"conf_target":8,"estimate_mode":"ECONOMICAL"}`)
			},
			staticCmd: func() interface{} {
				deadbeef, err := hex.DecodeString("deadbeef")
				if err != nil {
					panic(err)
				}
				changeAddress := "bcrt1qeeuctq9wutlcl5zatge7rjgx0k45228cxez655"
				change := 1
				changeType := "legacy"
				watching := true
				lockUnspents := true
				feeRate := 0.7
				replaceable := true
				confTarget := 8

				return hdfjson.NewFundRawTransactionCmd(deadbeef, hdfjson.FundRawTransactionOpts{
					ChangeAddress:          &changeAddress,
					ChangePosition:         &change,
					ChangeType:             &changeType,
					IncludeWatching:        &watching,
					LockUnspents:           &lockUnspents,
					FeeRate:                &feeRate,
					SubtractFeeFromOutputs: []int{0},
					Replaceable:            &replaceable,
					ConfTarget:             &confTarget,
					EstimateMode:           &hdfjson.EstimateModeEconomical,
				}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"fundrawtransaction","params":["deadbeef",{"changeAddress":"bcrt1qeeuctq9wutlcl5zatge7rjgx0k45228cxez655","changePosition":1,"change_type":"legacy","includeWatching":true,"lockUnspents":true,"feeRate":0.7,"subtractFeeFromOutputs":[0],"replaceable":true,"conf_target":8,"estimate_mode":"ECONOMICAL"}],"id":1}`,
			unmarshalled: func() interface{} {
				changeAddress := "bcrt1qeeuctq9wutlcl5zatge7rjgx0k45228cxez655"
				change := 1
				changeType := "legacy"
				watching := true
				lockUnspents := true
				feeRate := 0.7
				replaceable := true
				confTarget := 8
				return &hdfjson.FundRawTransactionCmd{
					HexTx: "deadbeef",
					Options: hdfjson.FundRawTransactionOpts{
						ChangeAddress:          &changeAddress,
						ChangePosition:         &change,
						ChangeType:             &changeType,
						IncludeWatching:        &watching,
						LockUnspents:           &lockUnspents,
						FeeRate:                &feeRate,
						SubtractFeeFromOutputs: []int{0},
						Replaceable:            &replaceable,
						ConfTarget:             &confTarget,
						EstimateMode:           &hdfjson.EstimateModeEconomical,
					},
					IsWitness: nil,
				}
			}(),
		},
		{
			name: "fundrawtransaction - iswitness",
			newCmd: func() (i interface{}, e error) {
				return hdfjson.NewCmd("fundrawtransaction", "deadbeef", "{}", true)
			},
			staticCmd: func() interface{} {
				deadbeef, err := hex.DecodeString("deadbeef")
				if err != nil {
					panic(err)
				}
				t := true
				return hdfjson.NewFundRawTransactionCmd(deadbeef, hdfjson.FundRawTransactionOpts{}, &t)
			},
			marshalled: `{"jsonrpc":"1.0","method":"fundrawtransaction","params":["deadbeef",{},true],"id":1}`,
			unmarshalled: &hdfjson.FundRawTransactionCmd{
				HexTx:   "deadbeef",
				Options: hdfjson.FundRawTransactionOpts{},
				IsWitness: func() *bool {
					t := true
					return &t
				}(),
			},
		},
		{
			name: "decoderawtransaction",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("decoderawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewDecodeRawTransactionCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decoderawtransaction","params":["123"],"id":1}`,
			unmarshalled: &hdfjson.DecodeRawTransactionCmd{HexTx: "123"},
		},
		{
			name: "decodescript",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("decodescript", "00")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewDecodeScriptCmd("00")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decodescript","params":["00"],"id":1}`,
			unmarshalled: &hdfjson.DecodeScriptCmd{HexScript: "00"},
		},
		{
			name: "getaddednodeinfo",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getaddednodeinfo", true)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetAddedNodeInfoCmd(true, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true],"id":1}`,
			unmarshalled: &hdfjson.GetAddedNodeInfoCmd{DNS: true, Node: nil},
		},
		{
			name: "getaddednodeinfo optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getaddednodeinfo", true, "127.0.0.1")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetAddedNodeInfoCmd(true, hdfjson.String("127.0.0.1"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true,"127.0.0.1"],"id":1}`,
			unmarshalled: &hdfjson.GetAddedNodeInfoCmd{
				DNS:  true,
				Node: hdfjson.String("127.0.0.1"),
			},
		},
		{
			name: "getbestblockhash",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getbestblockhash")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBestBlockHashCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblockhash","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetBestBlockHashCmd{},
		},
		{
			name: "getblock",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblock", "123", hdfjson.Int(0))
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockCmd("123", hdfjson.Int(0))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",0],"id":1}`,
			unmarshalled: &hdfjson.GetBlockCmd{
				Hash:      "123",
				Verbosity: hdfjson.Int(0),
			},
		},
		{
			name: "getblock default verbosity",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblock", "123")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123"],"id":1}`,
			unmarshalled: &hdfjson.GetBlockCmd{
				Hash:      "123",
				Verbosity: hdfjson.Int(1),
			},
		},
		{
			name: "getblock required optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblock", "123", hdfjson.Int(1))
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockCmd("123", hdfjson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",1],"id":1}`,
			unmarshalled: &hdfjson.GetBlockCmd{
				Hash:      "123",
				Verbosity: hdfjson.Int(1),
			},
		},
		{
			name: "getblock required optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblock", "123", hdfjson.Int(2))
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockCmd("123", hdfjson.Int(2))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",2],"id":1}`,
			unmarshalled: &hdfjson.GetBlockCmd{
				Hash:      "123",
				Verbosity: hdfjson.Int(2),
			},
		},
		{
			name: "getblockchaininfo",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblockchaininfo")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockChainInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockchaininfo","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetBlockChainInfoCmd{},
		},
		{
			name: "getblockcount",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblockcount")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockcount","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetBlockCountCmd{},
		},
		{
			name: "getblockhash",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblockhash", 123)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockHashCmd(123)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockhash","params":[123],"id":1}`,
			unmarshalled: &hdfjson.GetBlockHashCmd{Index: 123},
		},
		{
			name: "getblockheader",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblockheader", "123")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockHeaderCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockheader","params":["123"],"id":1}`,
			unmarshalled: &hdfjson.GetBlockHeaderCmd{
				Hash:    "123",
				Verbose: hdfjson.Bool(true),
			},
		},
		{
			name: "getblockstats height",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblockstats", hdfjson.HashOrHeight{Value: 123})
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockStatsCmd(hdfjson.HashOrHeight{Value: 123}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockstats","params":[123],"id":1}`,
			unmarshalled: &hdfjson.GetBlockStatsCmd{
				HashOrHeight: hdfjson.HashOrHeight{Value: 123},
			},
		},
		{
			name: "getblockstats hash",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblockstats", hdfjson.HashOrHeight{Value: "deadbeef"})
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockStatsCmd(hdfjson.HashOrHeight{Value: "deadbeef"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockstats","params":["deadbeef"],"id":1}`,
			unmarshalled: &hdfjson.GetBlockStatsCmd{
				HashOrHeight: hdfjson.HashOrHeight{Value: "deadbeef"},
			},
		},
		{
			name: "getblockstats height optional stats",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblockstats", hdfjson.HashOrHeight{Value: 123}, []string{"avgfee", "maxfee"})
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockStatsCmd(hdfjson.HashOrHeight{Value: 123}, &[]string{"avgfee", "maxfee"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockstats","params":[123,["avgfee","maxfee"]],"id":1}`,
			unmarshalled: &hdfjson.GetBlockStatsCmd{
				HashOrHeight: hdfjson.HashOrHeight{Value: 123},
				Stats:        &[]string{"avgfee", "maxfee"},
			},
		},
		{
			name: "getblockstats hash optional stats",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblockstats", hdfjson.HashOrHeight{Value: "deadbeef"}, []string{"avgfee", "maxfee"})
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockStatsCmd(hdfjson.HashOrHeight{Value: "deadbeef"}, &[]string{"avgfee", "maxfee"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockstats","params":["deadbeef",["avgfee","maxfee"]],"id":1}`,
			unmarshalled: &hdfjson.GetBlockStatsCmd{
				HashOrHeight: hdfjson.HashOrHeight{Value: "deadbeef"},
				Stats:        &[]string{"avgfee", "maxfee"},
			},
		},
		{
			name: "getblocktemplate",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblocktemplate")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetBlockTemplateCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblocktemplate","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetBlockTemplateCmd{Request: nil},
		},
		{
			name: "getblocktemplate optional - template request",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"]}`)
			},
			staticCmd: func() interface{} {
				template := hdfjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				}
				return hdfjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"]}],"id":1}`,
			unmarshalled: &hdfjson.GetBlockTemplateCmd{
				Request: &hdfjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := hdfjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   500,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return hdfjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &hdfjson.GetBlockTemplateCmd{
				Request: &hdfjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   int64(500),
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks 2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := hdfjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return hdfjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &hdfjson.GetBlockTemplateCmd{
				Request: &hdfjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getcfilter",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getcfilter", "123",
					wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetCFilterCmd("123",
					wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilter","params":["123",0],"id":1}`,
			unmarshalled: &hdfjson.GetCFilterCmd{
				Hash:       "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getcfilterheader",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getcfilterheader", "123",
					wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetCFilterHeaderCmd("123",
					wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilterheader","params":["123",0],"id":1}`,
			unmarshalled: &hdfjson.GetCFilterHeaderCmd{
				Hash:       "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getchaintips",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getchaintips")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetChainTipsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintips","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetChainTipsCmd{},
		},
		{
			name: "getchaintxstats",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getchaintxstats")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetChainTxStatsCmd(nil, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintxstats","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetChainTxStatsCmd{},
		},
		{
			name: "getchaintxstats optional nblocks",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getchaintxstats", hdfjson.Int32(1000))
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetChainTxStatsCmd(hdfjson.Int32(1000), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getchaintxstats","params":[1000],"id":1}`,
			unmarshalled: &hdfjson.GetChainTxStatsCmd{
				NBlocks: hdfjson.Int32(1000),
			},
		},
		{
			name: "getchaintxstats optional nblocks and blockhash",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getchaintxstats", hdfjson.Int32(1000), hdfjson.String("0000afaf"))
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetChainTxStatsCmd(hdfjson.Int32(1000), hdfjson.String("0000afaf"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getchaintxstats","params":[1000,"0000afaf"],"id":1}`,
			unmarshalled: &hdfjson.GetChainTxStatsCmd{
				NBlocks:   hdfjson.Int32(1000),
				BlockHash: hdfjson.String("0000afaf"),
			},
		},
		{
			name: "getconnectioncount",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getconnectioncount")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetConnectionCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getconnectioncount","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetConnectionCountCmd{},
		},
		{
			name: "getdifficulty",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getdifficulty")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetDifficultyCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdifficulty","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetDifficultyCmd{},
		},
		{
			name: "getgenerate",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getgenerate")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetGenerateCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getgenerate","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetGenerateCmd{},
		},
		{
			name: "gethashespersec",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("gethashespersec")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetHashesPerSecCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gethashespersec","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetHashesPerSecCmd{},
		},
		{
			name: "getinfo",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getinfo")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetInfoCmd{},
		},
		{
			name: "getmempoolentry",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getmempoolentry", "txhash")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetMempoolEntryCmd("txhash")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getmempoolentry","params":["txhash"],"id":1}`,
			unmarshalled: &hdfjson.GetMempoolEntryCmd{
				TxID: "txhash",
			},
		},
		{
			name: "getmempoolinfo",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getmempoolinfo")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetMempoolInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmempoolinfo","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetMempoolInfoCmd{},
		},
		{
			name: "getmininginfo",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getmininginfo")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetMiningInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmininginfo","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetMiningInfoCmd{},
		},
		{
			name: "getnetworkinfo",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getnetworkinfo")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetNetworkInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnetworkinfo","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetNetworkInfoCmd{},
		},
		{
			name: "getnettotals",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getnettotals")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetNetTotalsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnettotals","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetNetTotalsCmd{},
		},
		{
			name: "getnetworkhashps",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getnetworkhashps")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetNetworkHashPSCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetNetworkHashPSCmd{
				Blocks: hdfjson.Int(120),
				Height: hdfjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getnetworkhashps", 200)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetNetworkHashPSCmd(hdfjson.Int(200), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200],"id":1}`,
			unmarshalled: &hdfjson.GetNetworkHashPSCmd{
				Blocks: hdfjson.Int(200),
				Height: hdfjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getnetworkhashps", 200, 123)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetNetworkHashPSCmd(hdfjson.Int(200), hdfjson.Int(123))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200,123],"id":1}`,
			unmarshalled: &hdfjson.GetNetworkHashPSCmd{
				Blocks: hdfjson.Int(200),
				Height: hdfjson.Int(123),
			},
		},
		{
			name: "getpeerinfo",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getpeerinfo")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetPeerInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getpeerinfo","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetPeerInfoCmd{},
		},
		{
			name: "getrawmempool",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getrawmempool")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetRawMempoolCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetRawMempoolCmd{
				Verbose: hdfjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getrawmempool", false)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetRawMempoolCmd(hdfjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false],"id":1}`,
			unmarshalled: &hdfjson.GetRawMempoolCmd{
				Verbose: hdfjson.Bool(false),
			},
		},
		{
			name: "getrawtransaction",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getrawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetRawTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123"],"id":1}`,
			unmarshalled: &hdfjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: hdfjson.Int(0),
			},
		},
		{
			name: "getrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getrawtransaction", "123", 1)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetRawTransactionCmd("123", hdfjson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123",1],"id":1}`,
			unmarshalled: &hdfjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: hdfjson.Int(1),
			},
		},
		{
			name: "gettxout",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("gettxout", "123", 1)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetTxOutCmd("123", 1, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1],"id":1}`,
			unmarshalled: &hdfjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: hdfjson.Bool(true),
			},
		},
		{
			name: "gettxout optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("gettxout", "123", 1, true)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetTxOutCmd("123", 1, hdfjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1,true],"id":1}`,
			unmarshalled: &hdfjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: hdfjson.Bool(true),
			},
		},
		{
			name: "gettxoutproof",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("gettxoutproof", []string{"123", "456"})
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetTxOutProofCmd([]string{"123", "456"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"]],"id":1}`,
			unmarshalled: &hdfjson.GetTxOutProofCmd{
				TxIDs: []string{"123", "456"},
			},
		},
		{
			name: "gettxoutproof optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("gettxoutproof", []string{"123", "456"},
					hdfjson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetTxOutProofCmd([]string{"123", "456"},
					hdfjson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"],` +
				`"000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"],"id":1}`,
			unmarshalled: &hdfjson.GetTxOutProofCmd{
				TxIDs:     []string{"123", "456"},
				BlockHash: hdfjson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"),
			},
		},
		{
			name: "gettxoutsetinfo",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("gettxoutsetinfo")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetTxOutSetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gettxoutsetinfo","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetTxOutSetInfoCmd{},
		},
		{
			name: "getwork",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getwork")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetWorkCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":[],"id":1}`,
			unmarshalled: &hdfjson.GetWorkCmd{
				Data: nil,
			},
		},
		{
			name: "getwork optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("getwork", "00112233")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewGetWorkCmd(hdfjson.String("00112233"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":["00112233"],"id":1}`,
			unmarshalled: &hdfjson.GetWorkCmd{
				Data: hdfjson.String("00112233"),
			},
		},
		{
			name: "help",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("help")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewHelpCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":[],"id":1}`,
			unmarshalled: &hdfjson.HelpCmd{
				Command: nil,
			},
		},
		{
			name: "help optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("help", "getblock")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewHelpCmd(hdfjson.String("getblock"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":["getblock"],"id":1}`,
			unmarshalled: &hdfjson.HelpCmd{
				Command: hdfjson.String("getblock"),
			},
		},
		{
			name: "invalidateblock",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("invalidateblock", "123")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewInvalidateBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"invalidateblock","params":["123"],"id":1}`,
			unmarshalled: &hdfjson.InvalidateBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "ping",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("ping")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewPingCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"ping","params":[],"id":1}`,
			unmarshalled: &hdfjson.PingCmd{},
		},
		{
			name: "preciousblock",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("preciousblock", "0123")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewPreciousBlockCmd("0123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"preciousblock","params":["0123"],"id":1}`,
			unmarshalled: &hdfjson.PreciousBlockCmd{
				BlockHash: "0123",
			},
		},
		{
			name: "reconsiderblock",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("reconsiderblock", "123")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewReconsiderBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"reconsiderblock","params":["123"],"id":1}`,
			unmarshalled: &hdfjson.ReconsiderBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("searchrawtransactions", "1Address")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSearchRawTransactionsCmd("1Address", nil, nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address"],"id":1}`,
			unmarshalled: &hdfjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hdfjson.Int(1),
				Skip:        hdfjson.Int(0),
				Count:       hdfjson.Int(100),
				VinExtra:    hdfjson.Int(0),
				Reverse:     hdfjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("searchrawtransactions", "1Address", 0)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSearchRawTransactionsCmd("1Address",
					hdfjson.Int(0), nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0],"id":1}`,
			unmarshalled: &hdfjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hdfjson.Int(0),
				Skip:        hdfjson.Int(0),
				Count:       hdfjson.Int(100),
				VinExtra:    hdfjson.Int(0),
				Reverse:     hdfjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("searchrawtransactions", "1Address", 0, 5)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSearchRawTransactionsCmd("1Address",
					hdfjson.Int(0), hdfjson.Int(5), nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5],"id":1}`,
			unmarshalled: &hdfjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hdfjson.Int(0),
				Skip:        hdfjson.Int(5),
				Count:       hdfjson.Int(100),
				VinExtra:    hdfjson.Int(0),
				Reverse:     hdfjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSearchRawTransactionsCmd("1Address",
					hdfjson.Int(0), hdfjson.Int(5), hdfjson.Int(10), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10],"id":1}`,
			unmarshalled: &hdfjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hdfjson.Int(0),
				Skip:        hdfjson.Int(5),
				Count:       hdfjson.Int(10),
				VinExtra:    hdfjson.Int(0),
				Reverse:     hdfjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSearchRawTransactionsCmd("1Address",
					hdfjson.Int(0), hdfjson.Int(5), hdfjson.Int(10), hdfjson.Int(1), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1],"id":1}`,
			unmarshalled: &hdfjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hdfjson.Int(0),
				Skip:        hdfjson.Int(5),
				Count:       hdfjson.Int(10),
				VinExtra:    hdfjson.Int(1),
				Reverse:     hdfjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSearchRawTransactionsCmd("1Address",
					hdfjson.Int(0), hdfjson.Int(5), hdfjson.Int(10), hdfjson.Int(1), hdfjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true],"id":1}`,
			unmarshalled: &hdfjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hdfjson.Int(0),
				Skip:        hdfjson.Int(5),
				Count:       hdfjson.Int(10),
				VinExtra:    hdfjson.Int(1),
				Reverse:     hdfjson.Bool(true),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSearchRawTransactionsCmd("1Address",
					hdfjson.Int(0), hdfjson.Int(5), hdfjson.Int(10), hdfjson.Int(1), hdfjson.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true,["1Address"]],"id":1}`,
			unmarshalled: &hdfjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hdfjson.Int(0),
				Skip:        hdfjson.Int(5),
				Count:       hdfjson.Int(10),
				VinExtra:    hdfjson.Int(1),
				Reverse:     hdfjson.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "sendrawtransaction",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendrawtransaction", "1122")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSendRawTransactionCmd("1122", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
			unmarshalled: &hdfjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: hdfjson.Bool(false),
			},
		},
		{
			name: "sendrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("sendrawtransaction", "1122", false)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSendRawTransactionCmd("1122", hdfjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &hdfjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: hdfjson.Bool(false),
			},
		},
		{
			name: "setgenerate",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("setgenerate", true)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSetGenerateCmd(true, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true],"id":1}`,
			unmarshalled: &hdfjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: hdfjson.Int(-1),
			},
		},
		{
			name: "setgenerate optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("setgenerate", true, 6)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSetGenerateCmd(true, hdfjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true,6],"id":1}`,
			unmarshalled: &hdfjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: hdfjson.Int(6),
			},
		},
		{
			name: "stop",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("stop")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewStopCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stop","params":[],"id":1}`,
			unmarshalled: &hdfjson.StopCmd{},
		},
		{
			name: "submitblock",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("submitblock", "112233")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewSubmitBlockCmd("112233", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233"],"id":1}`,
			unmarshalled: &hdfjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options:  nil,
			},
		},
		{
			name: "submitblock optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("submitblock", "112233", `{"workid":"12345"}`)
			},
			staticCmd: func() interface{} {
				options := hdfjson.SubmitBlockOptions{
					WorkID: "12345",
				}
				return hdfjson.NewSubmitBlockCmd("112233", &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233",{"workid":"12345"}],"id":1}`,
			unmarshalled: &hdfjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options: &hdfjson.SubmitBlockOptions{
					WorkID: "12345",
				},
			},
		},
		{
			name: "uptime",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("uptime")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewUptimeCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"uptime","params":[],"id":1}`,
			unmarshalled: &hdfjson.UptimeCmd{},
		},
		{
			name: "validateaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("validateaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewValidateAddressCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"validateaddress","params":["1Address"],"id":1}`,
			unmarshalled: &hdfjson.ValidateAddressCmd{
				Address: "1Address",
			},
		},
		{
			name: "verifychain",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("verifychain")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewVerifyChainCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[],"id":1}`,
			unmarshalled: &hdfjson.VerifyChainCmd{
				CheckLevel: hdfjson.Int32(3),
				CheckDepth: hdfjson.Int32(288),
			},
		},
		{
			name: "verifychain optional1",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("verifychain", 2)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewVerifyChainCmd(hdfjson.Int32(2), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2],"id":1}`,
			unmarshalled: &hdfjson.VerifyChainCmd{
				CheckLevel: hdfjson.Int32(2),
				CheckDepth: hdfjson.Int32(288),
			},
		},
		{
			name: "verifychain optional2",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("verifychain", 2, 500)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewVerifyChainCmd(hdfjson.Int32(2), hdfjson.Int32(500))
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2,500],"id":1}`,
			unmarshalled: &hdfjson.VerifyChainCmd{
				CheckLevel: hdfjson.Int32(2),
				CheckDepth: hdfjson.Int32(500),
			},
		},
		{
			name: "verifymessage",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("verifymessage", "1Address", "301234", "test")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewVerifyMessageCmd("1Address", "301234", "test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifymessage","params":["1Address","301234","test"],"id":1}`,
			unmarshalled: &hdfjson.VerifyMessageCmd{
				Address:   "1Address",
				Signature: "301234",
				Message:   "test",
			},
		},
		{
			name: "verifytxoutproof",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("verifytxoutproof", "test")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewVerifyTxOutProofCmd("test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifytxoutproof","params":["test"],"id":1}`,
			unmarshalled: &hdfjson.VerifyTxOutProofCmd{
				Proof: "test",
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
			t.Errorf("\n%s\n%s", marshalled, test.marshalled)
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

// TestChainSvrCmdErrors ensures any errors that occur in the command during
// custom mashal and unmarshal are as expected.
func TestChainSvrCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     interface{}
		marshalled string
		err        error
	}{
		{
			name:       "template request with invalid type",
			result:     &hdfjson.TemplateRequest{},
			marshalled: `{"mode":1}`,
			err:        &json.UnmarshalTypeError{},
		},
		{
			name:       "invalid template request sigoplimit field",
			result:     &hdfjson.TemplateRequest{},
			marshalled: `{"sigoplimit":"invalid"}`,
			err:        hdfjson.Error{ErrorCode: hdfjson.ErrInvalidType},
		},
		{
			name:       "invalid template request sizelimit field",
			result:     &hdfjson.TemplateRequest{},
			marshalled: `{"sizelimit":"invalid"}`,
			err:        hdfjson.Error{ErrorCode: hdfjson.ErrInvalidType},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.Unmarshal([]byte(test.marshalled), &test.result)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error - got %T (%v), "+
				"want %T", i, test.name, err, err, test.err)
			continue
		}

		if terr, ok := test.err.(hdfjson.Error); ok {
			gotErrorCode := err.(hdfjson.Error).ErrorCode
			if gotErrorCode != terr.ErrorCode {
				t.Errorf("Test #%d (%s) mismatched error code "+
					"- got %v (%v), want %v", i, test.name,
					gotErrorCode, terr, terr.ErrorCode)
				continue
			}
		}
	}
}
