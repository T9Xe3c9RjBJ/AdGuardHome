package querylog

import (
	"bytes"
	"log"
	"strings"
	"testing"

	aglog "github.com/AdguardTeam/golibs/log"
	"github.com/stretchr/testify/assert"
)

func TestDecode_decodeQueryLog(t *testing.T) {
	stdWriter := log.Writer()
	stdLevel := aglog.GetLevel()
	t.Cleanup(func() {
		log.SetOutput(stdWriter)
		aglog.SetLevel(stdLevel)
	})

	logOut := &bytes.Buffer{}
	log.SetOutput(logOut)

	aglog.SetLevel(aglog.DEBUG)

	testCases := []struct {
		name string
		log  string
		want string
	}{{
		name: "back_compatibility_all_right",
		log:  `{"Question":"ULgBAAABAAAAAAAAC2FkZ3VhcmR0ZWFtBmdpdGh1YgJpbwAAHAAB","Time":"2020-11-13T12:41:25.970861+03:00"}`,
		want: "default",
	}, {
		name: "back_compatibility_bad_msg",
		log:  `{"Question":"","Time":"2020-11-13T12:41:25.970861+03:00"}`,
		want: "decodeLogEntry err: dns: overflow unpacking uint16\n",
	}, {
		name: "back_compatibility_bad_decoding",
		log:  `{"Question":"LgBAAABAAAAAAAAC2FkZ3VhcmR0ZWFtBmdpdGh1YgJpbwAAHAAB","Time":"2020-11-13T12:41:25.970861+03:00"}`,
		want: "decodeLogEntry err: illegal base64 data at input byte 48\n",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := logOut.Write([]byte("default"))
			assert.Nil(t, err)

			l := &logEntry{}
			decodeLogEntry(l, tc.log)

			assert.True(t, strings.HasSuffix(logOut.String(), tc.want), logOut.String())

			logOut.Reset()
		})
	}
}

func TestJSON(t *testing.T) {
	s := `
	{"keystr":"val","obj":{"keybool":true,"keyint":123456}}
	`
	k, v, jtype := readJSON(&s)
	assert.Equal(t, jtype, int32(jsonTStr))
	assert.Equal(t, "keystr", k)
	assert.Equal(t, "val", v)

	k, _, jtype = readJSON(&s)
	assert.Equal(t, jtype, int32(jsonTObj))
	assert.Equal(t, "obj", k)

	k, v, jtype = readJSON(&s)
	assert.Equal(t, jtype, int32(jsonTBool))
	assert.Equal(t, "keybool", k)
	assert.Equal(t, "true", v)

	k, v, jtype = readJSON(&s)
	assert.Equal(t, jtype, int32(jsonTNum))
	assert.Equal(t, "keyint", k)
	assert.Equal(t, "123456", v)

	_, _, jtype = readJSON(&s)
	assert.True(t, jtype == jsonTErr)
}
