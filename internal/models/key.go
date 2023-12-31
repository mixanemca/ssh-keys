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

package models

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// Key represents a SSH key with additional info like a path.
type Key struct {
	Name          string
	Path          string
	Format        string
	Comment       string
	Private       any
	Public        ssh.PublicKey
	LoadedToAgent bool
}

// String implements fmt.Stringer interface
func (k *Key) String() string {
	return fmt.Sprintf("%s %s", k.Name, k.Comment)
}
