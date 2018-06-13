package zenroom

import (
	"fmt"
	"log"
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

func TestEncDecWithFixedKeys(t *testing.T) {

	data := []byte("secret string")

	keys := []byte(fmt.Sprintf(`{"public": "%s", "private": "%s" }`, `BBaUFUb+HOi7dlssY9ZWWOSlTqOg\/x1r7ceebT\/WpXhlJj+XlaCkNzWp2emaZ9Cdonn7aNriwB5NUGigmvctnEo=`, `BGVCmpVnnG4hor9niXvoVx6OKytyTwfjxPH3dbyezys=`))

	encryptScript := []byte(`
		octet = require 'octet'
		ecdh = require 'ecdh'
		json = require 'json'

		msg = octet.new(#DATA)
		msg:string(DATA)
		
		keys = json.decode(KEYS)
		keyring = ecdh.new('ec25519')
		
		public = octet.new()
		public:base64(keys.public)
		
		private = octet.new()
		private:base64(keys.private)
		keyring:public(public)
		keyring:private(private)
		
		sess = keyring:session(public)
		zmsg = keyring:encrypt(sess, msg):base64()
		print(zmsg)
	`)
	encryptedMsg, err := Exec(encryptScript, keys, data)
	if err != nil {
		log.Fatal(err)
	}

	decryptScript := []byte(`
		octet = require 'octet'
		ecdh = require 'ecdh'
		json = require 'json'
	
		zmsg = octet.new(#DATA)
		zmsg:base64(DATA)
	
		keys = json.decode(KEYS)
	
		keyring = ecdh.new('ec25519')
	
		public = octet.new()
		public:base64(keys.public)
	
		private = octet.new()
		private:base64(keys.private)
	
		keyring:public(public)
		keyring:private(private)
	
		sess = keyring:session(public)
		msg = keyring:decrypt(sess, zmsg)
		print(msg)
		`)
	decryptedMsg, err := Exec(decryptScript, keys, encryptedMsg)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(decryptedMsg, data) {
		t.Error()
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
