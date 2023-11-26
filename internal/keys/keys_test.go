/*
Copyright © 2023 Michael Bruskov <mixanemca@yandex.ru>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keys

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Generate keys:
	ssh-keygen -t rsa -b 1024 -f id_rsa -N ''
	ssh-keygen -t dsa -b 1024 -f id_dsa -N ''
	ssh-keygen -t ecdsa -b 521 -f id_ecdsa -N ''
	ssh-keygen -o -a 100 -t ed25519 -f id_ed25519 -N ''
	ssh-keygen -o -a 100 -t ed25519 -f id_ed25519_with_passphrase -N 'with-passphrase'
*/
const (
	keyRSA string = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAlwAAAAdzc2gtcn
NhAAAAAwEAAQAAAIEA0HZZuiPQb9dqcINmnnONEApY56IJoJBblo9rhtEdu0rbYYK6qKPm
53XpybFbIZEwkVvwcDKWFK8qmjrvme1W1++RjSTDXXM/0e2UkVvpos43KGE55qZw8tjtUH
ptvxHpdB6Zxu/1fYCnGQF+/C1NcmgeGu/FVtzX5HjB27RY6TUAAAII06XIb9OlyG8AAAAH
c3NoLXJzYQAAAIEA0HZZuiPQb9dqcINmnnONEApY56IJoJBblo9rhtEdu0rbYYK6qKPm53
XpybFbIZEwkVvwcDKWFK8qmjrvme1W1++RjSTDXXM/0e2UkVvpos43KGE55qZw8tjtUHpt
vxHpdB6Zxu/1fYCnGQF+/C1NcmgeGu/FVtzX5HjB27RY6TUAAAADAQABAAAAgFT8SuxF+Y
3/BlfEWiuy9AlcDo6wUrhw4cXpxm56BmL6y6FfSHXEDDjEq4Ecmwh+RoycLNOw69qW5wll
iZT5W3OEsx1qYzoZuwT7Mv0eNqNiOZFsn900w16iiSPZAz6jagWYeccWcGF5zw9JIAuvaV
yH7MBHr/FyQFH/+XKBp85pAAAAQQDTYH0ccdECVOxGg0lUgnaz/kS99n2M1pFl2hiXNrcp
Z1S7RAQc0hh/hioFwmoIAZadKZndyukPDhldJKDwxbzpAAAAQQDuyMMKwoBBEbEv/0pLvY
ooLUVh/ILRUhL0jfm5T6rmX9UJZz7S0g0kA42VqAYxxjt/mMxbszkQbT1u8/O3S2BzAAAA
QQDffe/vL9QB+WdCOv9EQopHDp+UUQO+Ho34fu9n/6HHdiRbg0oK3Ka05iUlXTqhR3CxII
az9L9mLwO1Ucg5gm23AAAAD21ickB0aG9yYS5sb2NhbAECAw==
-----END OPENSSH PRIVATE KEY-----`
	keyDSA string = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABsgAAAAdzc2gtZH
NzAAAAgQDujBbIa1P11VQn42XPmM9Xo9sb3LVHh2hjg6ukSTj+pC0L0wv+dMTxNLHrXiEX
duwzwoJzA75sBtZ7vvGX+1XEfn5bvJ7kjZxULUS9RZDBJ39R8OrcpB2PgpIhyWOUUodf9X
W9hzBAMkVzIv/1ALm+v9ewlnT+MWdfu/E9OXr2KQAAABUAo1t3sgqQZXpDA/cahRDW7TgX
u2EAAACBAL8AYTEQgiEehDJwMJmHkvhatuTa2jaGYk70D32+2UFnn7jHl6RRJtL2jaUGQu
FE3jSlO0qATH+vgUcg492xKp8QM7ZBwVImbWBQtXoE2c13dyWuEG/UJg9qhkpdIkhoV2Td
xNHLycI6ZLhdta+cyTp9L6MmMC/aZ/PSwNy4w3pkAAAAgAgEG62AcjTs/FWFqud0hAmzWc
f4xjKHgEdy47AFVxpwx6Bg4is89XKv9i/bA267rbhBQPh/ekSRKI0zatv5RNz7ld+rv6h0
CD2V7fp6eIxuTcE1zOtw24Ou96+la/2KjupFzCpHaqBmZBKePkpmk08mn+wqqCjrVD+PiS
VXzi0lAAAB6KilYxSopWMUAAAAB3NzaC1kc3MAAACBAO6MFshrU/XVVCfjZc+Yz1ej2xvc
tUeHaGODq6RJOP6kLQvTC/50xPE0seteIRd27DPCgnMDvmwG1nu+8Zf7VcR+flu8nuSNnF
QtRL1FkMEnf1Hw6tykHY+CkiHJY5RSh1/1db2HMEAyRXMi//UAub6/17CWdP4xZ1+78T05
evYpAAAAFQCjW3eyCpBlekMD9xqFENbtOBe7YQAAAIEAvwBhMRCCIR6EMnAwmYeS+Fq25N
raNoZiTvQPfb7ZQWefuMeXpFEm0vaNpQZC4UTeNKU7SoBMf6+BRyDj3bEqnxAztkHBUiZt
YFC1egTZzXd3Ja4Qb9QmD2qGSl0iSGhXZN3E0cvJwjpkuF21r5zJOn0voyYwL9pn89LA3L
jDemQAAACACAQbrYByNOz8VYWq53SECbNZx/jGMoeAR3LjsAVXGnDHoGDiKzz1cq/2L9sD
brutuEFA+H96RJEojTNq2/lE3PuV36u/qHQIPZXt+np4jG5NwTXM63Dbg673r6Vr/YqO6k
XMKkdqoGZkEp4+SmaTTyaf7CqoKOtUP4+JJVfOLSUAAAAUaXjZiR9imj6dNPHw9jc8IORa
1PIAAAAPbWJyQHRob3JhLmxvY2FsAQID
-----END OPENSSH PRIVATE KEY-----`
	keyECDSA string = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAArAAAABNlY2RzYS
1zaGEyLW5pc3RwNTIxAAAACG5pc3RwNTIxAAAAhQQBANm8pMZEVOy3GZVCk/7zYCXGGnra
6NmlZdghHtnSJh/5uIx/CEo+RJ/6srOtZCz+inscULTuWzqPigKZxY4xVj8Ab/b7NpmLy8
9j7m811yBF78W03aA1zX0wCmKOBy/vq9yFjYKEu7g+NoQYWx4uAb6gIAs3edMfzmxD8n6t
4eaIL0AAAAEQvoowzb6KMM0AAAATZWNkc2Etc2hhMi1uaXN0cDUyMQAAAAhuaXN0cDUyMQ
AAAIUEAQDZvKTGRFTstxmVQpP+82Alxhp62ujZpWXYIR7Z0iYf+biMfwhKPkSf+rKzrWQs
/op7HFC07ls6j4oCmcWOMVY/AG/2+zaZi8vPY+5vNdcgRe/FtN2gNc19MApijgcv76vchY
2ChLu4PjaEGFseLgG+oCALN3nTH85sQ/J+reHmiC9AAAAAQgHZwdPlvU1ai74deprivKKV
ol9ow+Om6QQyVrTrNCIev4CWSuE5i1dZZAix6OEu4o84Hy9C6756uSGIRSciT2FdhAAAAA
9tYnJAdGhvcmEubG9jYWwBAgM=
-----END OPENSSH PRIVATE KEY-----`
	keyEd25519 string = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACCKNIzSklTj21fPxMTZGBqbsZnjL6yxq2ElEgKY7FgQegAAAJgDqr8/A6q/
PwAAAAtzc2gtZWQyNTUxOQAAACCKNIzSklTj21fPxMTZGBqbsZnjL6yxq2ElEgKY7FgQeg
AAAEB1uJSCcTimhYwSjImzRDNTpvZqkGdRKQVRu2xxxXmQxIo0jNKSVOPbV8/ExNkYGpux
meMvrLGrYSUSApjsWBB6AAAAD21ickB0aG9yYS5sb2NhbAECAwQFBg==
-----END OPENSSH PRIVATE KEY-----`
	keyEd25519WithPassphrase string = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABDQ6no/lP
U5Ux98BYl9fKPIAAAAZAAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIHX9LlEhDhEuW+dG
mTQo5g00CQQAAXpx8WHFe47eTHtmAAAAoKBEdarUsbFCj97azlqVHKpSldLi3cUB/XmZok
znn2ikuxv9oCQJ2RL+6M5TOzSBX6Sa5aoocwHnCSgIas6xBs2VviHTQvhm6CgoS5KjKdgK
vDet5DRqs3jd1HQ9eLUJdkfVuO7WnfVzPhT97hJlAdw4RULfCsk9vZRtVoIQUhbt55mq/j
lFHj0Z0aVSyj520JlURJKW3DVO/DxKr4cFZeY=
-----END OPENSSH PRIVATE KEY-----`
	keyEmpty  string = ``
	keyRandom string = `pFA?mH3y6.$snFd;o@ho5VI2.6:;X|gA87)X9CVj`
)

func TestIsPrivateKey(t *testing.T) {
	cases := []struct {
		name string
		key  []byte
		want bool
	}{
		{"Test valid RSA key", []byte(keyRSA), true},
		{"Test unsupported DSA key", []byte(keyDSA), false},
		{"Test valid ECDSA key", []byte(keyECDSA), true},
		{"Test valid Ed25519 key", []byte(keyEd25519), true},
		{"Test valid Ed25519 key with passphrase", []byte(keyEd25519WithPassphrase), true},
		{"Test empty key", []byte(keyEmpty), false},
		{"Test random key", []byte(keyRandom), false},
	}

	for _, c := range cases {
		got := isPrivateKey(c.key)
		assert.Equal(t, c.want, got, c.name)
	}
}

func TestLoadPrivateKeys(t *testing.T) {
	dir := prepareTestKeysDir(t)
	defer os.RemoveAll(dir)

	files := []string{
		filepath.Join(dir, "id_ecdsa"),
		filepath.Join(dir, "id_ed25519"),
		filepath.Join(dir, "id_ed25519_with_passphrase"),
		filepath.Join(dir, "id_rsa"),
	}

	cases := []struct {
		name string
		root string
		want []string
		// wantErr     bool
		expectedErr string
	}{
		{"Test valid keys dir", dir, files, ""},
		{"Test with noexists keys dir", "noexists", nil, ""},
	}

	for _, c := range cases {
		got, err := LoadPrivateKeys(c.root)
		assert.Nil(t, err)
		// if assert.Nil(t, err) {
		// assert.EqualError(t, err, c.expectedErr, c.name)
		// }
		assert.Equal(t, c.want, got, c.name)
	}
}

func TestLoadPrivateKeysErr(t *testing.T) {
	dir := prepareTestKeysDirNoReadable(t)
	defer os.RemoveAll(dir)

	_, err := LoadPrivateKeys(dir)
	assert.EqualError(t, err, fmt.Sprintf("read key file: open %s: permission denied", filepath.Join(dir, "noreadable")))

	sdir := prepareTestKeysDirWithNoRegular(t)
	defer os.RemoveAll(sdir)

	_, err = LoadPrivateKeys(sdir)
	assert.Nil(t, err)
}

// prepareTestKeysDir creates new temp directory and files with private keys
func prepareTestKeysDir(t *testing.T) string {
	dir, err := os.MkdirTemp(".", "ssh-keys-test-dir")
	if err != nil {
		t.Fatal("failed to create temp dir: ", err)
	}

	files := map[string]string{
		"id_rsa":                     keyRSA,
		"id_dsa":                     keyDSA,
		"id_ecdsa":                   keyECDSA,
		"id_ed25519":                 keyEd25519,
		"id_ed25519_with_passphrase": keyEd25519WithPassphrase,
		"id_empty":                   keyEmpty,
		"id_random":                  keyRandom,
	}
	for k, v := range files {
		createFile(t, dir, k, v)
	}

	return dir
}

func prepareTestKeysDirNoReadable(t *testing.T) string {
	dir, err := os.MkdirTemp(".", "ssh-keys-test-dir")
	if err != nil {
		t.Fatal("failed to create temp dir: ", err)
	}

	name := "noreadable"
	createFile(t, dir, name, "")

	err = os.Chmod(filepath.Join(dir, name), 0000)
	if err != nil {
		t.Fatal("failed to chmod file: ", err)
	}

	return dir
}

func prepareTestKeysDirWithNoRegular(t *testing.T) string {
	dir, err := os.MkdirTemp(".", "ssh-keys-test-dir")
	if err != nil {
		t.Fatal("failed to create temp dir: ", err)
	}

	socketAddr := &net.UnixAddr{
		Name: filepath.Join(dir, "socket"),
		Net:  "unix",
	}

	_, err = net.ListenUnix("unix", socketAddr)
	if err != nil {
		t.Fatal("failed to listen unix socket: ", err)
	}

	return dir
}

func createFile(t *testing.T, dir, name, data string) {
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		t.Fatal("failed to create file with private key: ", err)
	}

	_, err = f.WriteString(data)
	if err != nil {
		t.Fatal("failed to write private key to file: ", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatal("failed to close file with private key: ", err)
	}
}
