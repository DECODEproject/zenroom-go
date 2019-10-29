package zenroom_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/DECODEproject/zenroom-go"
)

func TestMissingScript(t *testing.T) {
	_, err := zenroom.Exec(nil)
	if err == nil {
		t.Error("Expected error for nil script, got nil")
	}

	_, err = zenroom.Exec([]byte{})
	if err == nil {
		t.Errorf("Expected error for empty script, got nil")
	}
}

func TestBasicCall(t *testing.T) {
	script := []byte(`print(1)`)

	res, err := zenroom.Exec(script)
	if err != nil {
		t.Error(err)
	}

	if string(res) != "1" {
		t.Errorf("unexpected response: expected 'hello world', got '%v'", res)
	}
}

func TestCallStrings(t *testing.T) {
	testcases := []struct {
		label  string
		script []byte
		data   []byte
		resp   []byte
	}{
		{
			label:  "string variable",
			script: []byte(`hello = 'Hello World!' print(hello)`),
			resp:   []byte("Hello World!"),
		},
		{
			label:  "naked string",
			script: []byte(`print('hello')`),
			resp:   []byte("hello"),
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.label, func(t *testing.T) {
			res, err := zenroom.Exec(testcase.script, zenroom.WithData(testcase.data))
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(res, testcase.resp) {
				t.Errorf("calling [%s] got %s of len %d", testcase.script, res, len(res))
			}
		})
	}
}

func TestData(t *testing.T) {
	script := []byte(`print(DATA)`)
	data := []byte(`Hello data`)

	res, err := zenroom.Exec(script, zenroom.WithData(data))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if string(res) != "Hello data" {
		t.Errorf("Unexpected output, expected 'Hello data', got '%s'", res)
	}
}

func TestEmptyData(t *testing.T) {
	script := []byte(`print(DATA)`)

	_, err := zenroom.Exec(script, zenroom.WithData([]byte{}))
	if err == nil {
		t.Error("Expected error for empty data, got nil")
	}
}

func TestKeys(t *testing.T) {
	script := []byte(`print(KEYS)`)
	keys := []byte(`Hello keys`)

	res, err := zenroom.Exec(script, zenroom.WithKeys(keys))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if string(res) != "Hello keys" {
		t.Errorf("Unexpected output, expected 'Hello keys', got '%s'", res)
	}
}

func TestEmptyKeys(t *testing.T) {
	script := []byte(`print(KEYS)`)

	_, err := zenroom.Exec(script, zenroom.WithKeys([]byte{}))
	if err == nil {
		t.Error("Expected error for empty keys, got nil")
	}
}

func TestEncodeDecode(t *testing.T) {
	encryptKeys := []byte(`
	{
 		"device_token": "abc123",
 		"community_id": "foo",
 		"community_pubkey": "u64:BA2U3SX2mNFOTFA2K05tkHlZadaDHftkKedKXqFjERZ9df6VFuZIgF20q0kjn9uy2vaYSYx6zEm1zrwvV3vwovc"
	}
	`)

	data := []byte(`{"msg": "secret"}`)

	encryptScript := []byte(`
-- Encryption script for DECODE IoT Pilot
curve = 'ed25519'

-- import and validate KEYS data
keys = JSON.decode(KEYS)

-- generate a new device keypair every time
device_key = ECDH.keygen(curve)

-- read the payload we will encrypt
payload = JSON.decode(DATA)

-- The device's public key, community_id and the curve type are tranmitted in
-- clear inside the header, which is authenticated AEAD
header = {}
header['device_pubkey'] = device_key:public():base64()
header['community_id'] = keys['community_id']

iv = O.random(16)
header['iv'] = iv:url64()

-- encrypt the data, and build our output object
local pub = ECDH.new(curve)
pub:public(url64(keys.community_pubkey))

local session = device_key:session(pub)
local head = url64(JSON.encode(header))
local out = { header = head }
out.text, out.checksum = ECDH.aead_encrypt(session, url64(JSON.encode(payload)), iv, head)

-- output = map(out, base64)
out.zenroom = VERSION
out.curve = curve

print(JSON.encode(out))
`)

	decryptKeys := []byte(`
	{
		"community_seckey": "u64:Cf88o0bEY3igf3mbnKTT7s7_huXDPvlATz7J1T7atZo"
	}
	`)

	decryptScript := []byte(`
	-- Decryption script for DECODE IoT Pilot

	-- curve used
	curve = 'ed25519'

	-- read and validate data
	keys = JSON.decode(KEYS)
	data = JSON.decode(DATA)
	header = JSON.decode(data.header)

	community_key = ECDH.new(curve)
	community_key:private(url64(keys.community_seckey))

	local pub = ECDH.new(curve)
	pub:public(base64(header.device_pubkey))
	session = community_key:session(pub)

	decode = { header = header }
	decode.text, decode.checksum = ECDH.aead_decrypt(session, url64(data.text), url64(header.iv), url64(data.header))

	print(decode.text:str())
	`)

	encryptedMessage, err := zenroom.Exec(encryptScript, zenroom.WithData(data), zenroom.WithKeys(encryptKeys))
	if err != nil {
		t.Fatalf("Error encrypting message: %v", err)
	}

	if len(encryptedMessage) == 0 {
		t.Errorf("Length of encrypted message should not be 0")
	}

	decryptedMessage, err := zenroom.Exec(decryptScript, zenroom.WithData(encryptedMessage), zenroom.WithKeys(decryptKeys))
	if err != nil {
		t.Fatalf("Error encrypting message: %v", err)
	}

	var decrypted map[string]interface{}
	err = json.Unmarshal(decryptedMessage, &decrypted)
	if err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
	}

	if decrypted["msg"] != "secret" {
		t.Errorf("Unexpected decrypted output, got %v, expected %v", decrypted["msg"], "secret")
	}
}

func BenchmarkBasicPrint(b *testing.B) {
	script := []byte(`print ('hello')`)
	for n := 0; n < b.N; n++ {
		_, _ = zenroom.Exec(script)
	}
}

func BenchmarkBasicKeyandEncrypt(b *testing.B) {
	script := []byte(`
	msg = str(DATA)
	kr = ECDH.new()
	kr:keygen()
	encrypted = ECDH.encrypt(kr, kr:public(), msg, kr:public())
	print (encrypted)
	`)
	data := []byte(`temperature:25.1`)

	for n := 0; n < b.N; n++ {
		_, _ = zenroom.Exec(script, zenroom.WithData(data))
	}
}

func ExampleExec() {
	script := []byte(`print("hello world")`)
	res, _ := zenroom.Exec(script)
	fmt.Println(string(res))
	// Output: hello world
}
