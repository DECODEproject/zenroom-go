package zenroom

import (
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

func BenchmarkBasicPrint(b *testing.B) {
	script := `print ('hello')`
	for n := 0; n < b.N; n++ {
		_, _ = Exec(script, "", "")
	}
}

func BenchmarkBasicKeyandEncrypt(b *testing.B) {
	script := `
	octet = require 'octet'
	ecdh = require 'ecdh'
	msg = octet.new(#DATA)
	msg:string(DATA)
	kr = ecdh.new()
	kr:keygen()
	sess = kr:session(kr:private(), kr:public())
	encrypted = kr:(sess, msg)
	print (encrypted)
	`
	data := `temperature:25.1`

	for n := 0; n < b.N; n++ {
		_, _ = Exec(script, "", data)
	}
}
