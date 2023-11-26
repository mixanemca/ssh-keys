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

	"golang.org/x/crypto/ssh"
)

func isPrivateKey(p []byte) bool {
	_, err := ssh.ParsePrivateKey(p)
	if err != nil {
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			return true
		}
		return false
	}
	return true
}

func LoadPrivateKeys(root string) ([]string, error) {
	var keys []string
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if path == root {
			return nil
		}
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		privateBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read key file: %v", err)
		}

		if ok := isPrivateKey(privateBytes); ok {
			keys = append(keys, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", root, err)
		return nil, err
	}

	return keys, nil
}
