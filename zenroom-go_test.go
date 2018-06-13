package zenroom

import (
	"reflect"
	"testing"
)

func TestBasicCall(t *testing.T) {
	script := []byte(`print (1)`)
	res, err := Exec(script, nil, nil)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(res, []byte("1")) {
		t.Errorf("calling print (1), got:%s len:%d", res, len(res))
	}
}

func TestCallStrings(t *testing.T) {
	testcases := []struct {
		script []byte
		data   []byte
		resp   []byte
	}{
		{
			script: []byte(`hello = 'Hello World!' print(hello)`),
			resp:   []byte("Hello World!"),
		},
		{
			script: []byte(`print('hello')`),
			resp:   []byte("hello"),
		},
		{
			script: []byte(`print(123)`),
			resp:   []byte("123"),
		},
	}
	for _, testcase := range testcases {
		res, err := Exec(testcase.script, nil, testcase.data)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(res, testcase.resp) {
			t.Errorf("calling [%s] got %s of len %d", testcase.script, res, len(res))
		}
	}
}

func TestEncDec(t *testing.T) {
	testcases := []struct {
		script []byte
		data   []byte
		resp   []byte
	}{
		{
			script: []byte(`
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
			`),
			data: []byte("UltraSuper Message!"),
			resp: []byte("UltraSuper Message!"),
		},
	}
	for _, testcase := range testcases {
		res, err := Exec(testcase.script, nil, testcase.data)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(res, testcase.resp) {
			t.Errorf("calling [%s] got %s of len %d", testcase.script, res, len(res))
		}
	}
}

func BenchmarkBasicPrint(b *testing.B) {
	script := []byte(`print ('hello')`)
	for n := 0; n < b.N; n++ {
		_, _ = Exec(script, nil, nil)
	}
}

func BenchmarkBasicKeyandEncrypt(b *testing.B) {
	script := []byte(`
	octet = require 'octet'
	ecdh = require 'ecdh'
	msg = octet.new(#DATA)
	msg:string(DATA)
	kr = ecdh.new()
	kr:keygen()
	sess = kr:session(kr:private(), kr:public())
	encrypted = kr:(sess, msg)
	print (encrypted)
	`)
	data := []byte(`temperature:25.1`)

	for n := 0; n < b.N; n++ {
		_, _ = Exec(script, nil, data)
	}
}
