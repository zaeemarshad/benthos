package nats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/benthosdev/benthos/v4/public/service"
)

func TestOutputJetStreamConfigParse(t *testing.T) {
	spec := natsJetStreamOutputConfig()
	env := service.NewEnvironment()

	t.Run("Successful config parsing", func(t *testing.T) {
		outputConfig := `
urls: [ url1, url2 ]
subject: testsubject
headers:
  Content-Type: application/json
  Timestamp: ${!meta("Timestamp")}
auth:
  nkey_file: test auth n key file
  user_credentials_file: test auth user creds file
  user_jwt: test auth inline user JWT
  user_nkey_seed: test auth inline user NKey Seed
`

		conf, err := spec.ParseYAML(outputConfig, env)
		require.NoError(t, err)

		e, err := newJetStreamWriterFromConfig(conf, service.MockResources())
		require.NoError(t, err)

		msg := service.NewMessage((nil))
		msg.MetaSet("Timestamp", "1651485106")
		assert.Equal(t, "url1,url2", e.urls)
		assert.Equal(t, "testsubject", e.subjectStr.String(msg))
		assert.Equal(t, "application/json", e.headers["Content-Type"].String(msg))
		assert.Equal(t, "1651485106", e.headers["Timestamp"].String(msg))
		assert.Equal(t, "test auth n key file", e.authConf.NKeyFile)
		assert.Equal(t, "test auth user creds file", e.authConf.UserCredentialsFile)
		assert.Equal(t, "test auth inline user JWT", e.authConf.UserJWT)
		assert.Equal(t, "test auth inline user NKey Seed", e.authConf.UserNkeySeed)
	})

	t.Run("Missing user_nkey_seed", func(t *testing.T) {
		inputConfig := `
urls: [ url1, url2 ]
subject: testsubject
auth:
  user_jwt: test auth inline user JWT
`

		conf, err := spec.ParseYAML(inputConfig, env)
		require.NoError(t, err)

		_, err = newJetStreamReaderFromConfig(conf, service.MockResources())
		require.Error(t, err)
	})

	t.Run("Missing user_jwt", func(t *testing.T) {
		inputConfig := `
urls: [ url1, url2 ]
subject: testsubject
auth:
  user_jwt: test auth inline user JWT
`

		conf, err := spec.ParseYAML(inputConfig, env)
		require.NoError(t, err)

		_, err = newJetStreamReaderFromConfig(conf, service.MockResources())
		require.Error(t, err)
	})
}
