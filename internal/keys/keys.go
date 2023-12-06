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
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mixanemca/ssh-keys/internal/models"
	"golang.org/x/crypto/ssh"
)

func isPrivateKey(p []byte) (ssh.Signer, bool) {
	signer, err := ssh.ParsePrivateKey(p)
	if err != nil {
		return nil, false
	}
	return signer, true
}

func LoadPrivateKeys(root string) ([]*models.Key, error) {
	keys := make([]*models.Key, 0)
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if path == root {
			return nil
		}

		if err != nil {
			return fmt.Errorf("prevent panic by handling failure accessing a path")
		}

		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		privateBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read key file: %v", err)
		}

		// Try to read the comment from public key file
		var comment string
		publicBytes, _ := os.ReadFile(path + ".pub")
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("read public key file: %v", err)
		}
		if len(publicBytes) > 0 {
			_, comment, _, _, _ = ssh.ParseAuthorizedKey(publicBytes)
		}

		if signer, ok := isPrivateKey(privateBytes); ok {
			if name, err := filepath.Rel(root, path); err == nil {
				if privKey, err := ssh.ParseRawPrivateKey(privateBytes); err == nil {
					key := &models.Key{
						Name:    name,
						Path:    path,
						Format:  signer.PublicKey().Type(),
						Comment: comment,
						Private: privKey,
						Public:  signer.PublicKey(),
					}
					keys = append(keys, key)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return keys, nil
}
