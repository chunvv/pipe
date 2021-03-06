// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crypto

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRSAEncryptDecrypt(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/public-rsa-pem")
	require.NoError(t, err)

	encrypter, err := NewRSAEncrypter(string(data))
	require.NoError(t, err)

	text := "short-text"

	encryptedText, err := encrypter.Encrypt(text)
	require.NoError(t, err)
	assert.True(t, len(encryptedText) > 0)

	fmt.Println(encryptedText)
	decrypter, err := NewRSADecrypter("testdata/private-rsa-pem")
	require.NoError(t, err)

	decryptedText, err := decrypter.Decrypt(encryptedText)
	require.NoError(t, err)
	assert.Equal(t, text, decryptedText)
}
