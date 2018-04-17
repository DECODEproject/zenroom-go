package zenroom

import (
	"fmt"
	"testing"
)

func TestBasicCall(t *testing.T) {
	script := `print (1)`
	res, err := Exec(script, "", "")
	if err != nil {
		t.Error(err)
	}
	if res != "1" {
		t.Errorf("calling print (1), got:%s len:%d", res, len(res))
	}
}

func TestBasicString(t *testing.T) {
	script := `print (1)`
	res, err := Exec(script, "", "")
	if err != nil {
		t.Error(err)
	}
	if res != "1" {
		t.Errorf("calling print (1), got:%s len:%d", res, len(res))
	}
}

func TestCallStrings(t *testing.T) {
	testcases := []struct {
		script string
		data   string
		resp   string
	}{
		{
			script: `hello = 'Hello World!' print(hello)`,
			resp:   "Hello World!",
		},
		{
			script: `print('hello')`,
			resp:   "hello",
		},
		{
			script: `print(123)`,
			resp:   "123",
		},
	}
	for _, testcase := range testcases {
		res, err := Exec(testcase.script, "", testcase.data)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("here", res)
		if res != testcase.resp {
			t.Errorf("calling [%s] got %s of len %d", testcase.script, res, len(res))
		}
	}
}

func TestEncDec(t *testing.T) {
	testcases := []struct {
		script string
		data   string
		resp   string
	}{
		{
			script: `
			octet = require'octet'
			ecdh = require 'ecdh'
			msg = octet.new(#DATA)
			msg:string(DATA)

			ed25519 =ecdh.new('ec25519')
			pk, sk = ed25519:keygen()

			sess = ed25519:session(pk, sk)

			zmsg = ed25519:encrypt(sess, msg)

			decipher = ed25519:decrypt(sess, zmsg)
			print(decipher:string())
			`,
			data: "UltraSuper Message!",
			resp: "UltraSuper Message!",
		},
	}
	for _, testcase := range testcases {
		res, err := Exec(testcase.script, "", testcase.data)
		if err != nil {
			t.Error(err)
		}
		if res != testcase.resp {
			t.Errorf("calling [%s] got %s of len %d", testcase.script, res, len(res))
		}
	}
}
func TestExecToBuf(t *testing.T) {
	script := `print ('hello')`
	s, err := ExecToBuf(script, "", "")
	if err != nil {
		t.Error(err)
	}
	if s != "hello" {
		t.Errorf("results aren't the same %v!=%v", s, "hello")
	}
}

/*
func TestEncrypt(t *testing.T) {
	genKeysScript := `
	json = require'json'
	ecdh = require'ecdh'	keypairs = json.encode({
		uno=keyring:public():base64(),
		dos=keyring:private():base64()
	})
	print(keypairs)
	`
	encodeScript := `
	json = require'json'
	ecdh = require'ecdh'

	`
}

	keyring = ecdh.new()
	keyring:keygen()

	keypairs = json.encode({
		uno=keyring:public():base64(),
		dos=keyring:private():base64()
	})
	print(keypairs)
	`
	encodeScript := `
	json = require'json'
	ecdh = require'ecdh'

	`
}
*/
