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

// TestHdfWalletExtCmds tests all of the hdfwallet extended commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestHdfWalletExtCmds(t *testing.T) {
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
			name: "createnewaccount",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("createnewaccount", "acct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewCreateNewAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"createnewaccount","params":["acct"],"id":1}`,
			unmarshalled: &hdfjson.CreateNewAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "dumpwallet",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("dumpwallet", "filename")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewDumpWalletCmd("filename")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpwallet","params":["filename"],"id":1}`,
			unmarshalled: &hdfjson.DumpWalletCmd{
				Filename: "filename",
			},
		},
		{
			name: "importaddress",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("importaddress", "1Address", "")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewImportAddressCmd("1Address", "", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importaddress","params":["1Address",""],"id":1}`,
			unmarshalled: &hdfjson.ImportAddressCmd{
				Address: "1Address",
				Rescan:  hdfjson.Bool(true),
			},
		},
		{
			name: "importaddress optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("importaddress", "1Address", "acct", false)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewImportAddressCmd("1Address", "acct", hdfjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importaddress","params":["1Address","acct",false],"id":1}`,
			unmarshalled: &hdfjson.ImportAddressCmd{
				Address: "1Address",
				Account: "acct",
				Rescan:  hdfjson.Bool(false),
			},
		},
		{
			name: "importpubkey",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("importpubkey", "031234")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewImportPubKeyCmd("031234", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importpubkey","params":["031234"],"id":1}`,
			unmarshalled: &hdfjson.ImportPubKeyCmd{
				PubKey: "031234",
				Rescan: hdfjson.Bool(true),
			},
		},
		{
			name: "importpubkey optional",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("importpubkey", "031234", false)
			},
			staticCmd: func() interface{} {
				return hdfjson.NewImportPubKeyCmd("031234", hdfjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importpubkey","params":["031234",false],"id":1}`,
			unmarshalled: &hdfjson.ImportPubKeyCmd{
				PubKey: "031234",
				Rescan: hdfjson.Bool(false),
			},
		},
		{
			name: "importwallet",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("importwallet", "filename")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewImportWalletCmd("filename")
			},
			marshalled: `{"jsonrpc":"1.0","method":"importwallet","params":["filename"],"id":1}`,
			unmarshalled: &hdfjson.ImportWalletCmd{
				Filename: "filename",
			},
		},
		{
			name: "renameaccount",
			newCmd: func() (interface{}, error) {
				return hdfjson.NewCmd("renameaccount", "oldacct", "newacct")
			},
			staticCmd: func() interface{} {
				return hdfjson.NewRenameAccountCmd("oldacct", "newacct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"renameaccount","params":["oldacct","newacct"],"id":1}`,
			unmarshalled: &hdfjson.RenameAccountCmd{
				OldAccount: "oldacct",
				NewAccount: "newacct",
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
