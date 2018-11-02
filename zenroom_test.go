package zenroom_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/thingful/zenroom-go"
)

func TestBasicCall(t *testing.T) {
	script := []byte(`print (1)`)
	res, err := zenroom.Exec(script)
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
		res, err := zenroom.Exec(testcase.script, zenroom.WithData(testcase.data))
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(res, testcase.resp) {
			t.Errorf("calling [%s] got %s of len %d", testcase.script, res, len(res))
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	encryptKeys := []byte(`
	{
 		"device_id": "anonymous",
 		"community_id": "smartcitizens",
 		"community_pubkey": "BBLewg4VqLR38b38daE7Fj\/uhr543uGrEpyoPFgmFZK6EZ9g2XdK\/i65RrSJ6sJ96aXD3DJHY3Me2GJQO9\/ifjE="
	}
	`)

	data := []byte(`secret message`)

	encryptScript := []byte(`
	curve = 'ed25519'

	keys_schema = SCHEMA.Record {
		device_id        = SCHEMA.String,
		community_id     = SCHEMA.String,
		community_pubkey = SCHEMA.String
	}

	payload_schema = SCHEMA.Record {
		device_id = SCHEMA.String,
		data      = SCHEMA.String
	}

	output_schema = SCHEMA.Record {
		device_pubkey = SCHEMA.String,
		community_id  = SCHEMA.String,
		payload       = SCHEMA.String
	}

	keys = read_json(KEYS, keys_schema)

	devkey = ECDH.keygen(curve)

	payload = {}
	payload['device_id'] = keys['device_id']
	payload['data']      = DATA
	validate(payload, payload_schema)

	header = {}
	header['device_pubkey'] = devkey:public():base64()
	header['community_id'] = keys['community_id']

	output = ECDH.encrypt(
		devkey,
		base64(keys.community_pubkey),
		MSG.pack(payload),
		MSG.pack(header)
	)

	output = map(output, O.to_base64)
	output.zenroom = VERSION
	output.encoding = 'base64'
	output.curve = curve

	print(JSON.encode(output))
	`)

	decryptKeys := []byte(`
	{
		"community_seckey": "D19GsDTGjLBX23J281SNpXWUdu+oL6hdAJ0Zh6IrRHA="
	}
	`)

	decryptScript := []byte(`
	keys_schema = SCHEMA.Record { community_seckey = SCHEMA.String }

	data_schema = SCHEMA.Record {
   	text     = SCHEMA.string,
   	iv       = SCHEMA.string,
   	header   = SCHEMA.string,
   	checksum = SCHEMA.string
	}

	payload_schema = SCHEMA.Record {
  	device_id   = SCHEMA.String,
  	data        = SCHEMA.String
	}

	data = read_json(DATA) -- TODO: data_schema validation
	keys = read_json(KEYS, keys_schema)
	head = OCTET.msgunpack( base64(data.header) )

	dashkey = ECDH.new()
	dashkey:private( base64(keys.community_seckey) )

	payload,ck = ECDH.decrypt(dashkey,
  	base64( head.device_pubkey ),
   	map(data, base64))

	validate(payload, payload_schema)

-- print("Header:")
-- content(msgunpack(payload.header) )
	print(JSON.encode(OCTET.msgunpack(payload.text) ))
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

	if decrypted["data"] != "secret message" {
		t.Errorf("Unexpected decrypted output, got %s, expected %s", decrypted["data"], "secret message")
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
	script := []byte(`print("hello")`)
	res, _ := zenroom.Exec(script)
	fmt.Println(string(res))
	// Output: hello
}
