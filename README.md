# ssh-keys

Work with SSH keys easily!

## What is it?

`ssh-keys` is a terminal based utility designed from easily working with SSH keys.

Use it to discover existing, creates a new and deletes the SSH keys and for work with `ssh-agent(1)`.

Carefully! Right now there is an active development stage.

## Installation

### Go

```bash
go install github.com/mixanemca/ssh-keys@latest
```

### Build (requires Go 1.21+)

```bash
git clone https://github.com/mixanemca/ssh-keys.git
cd ssh-keys
make
```

## TODO

- [x] List all available ssh private keys
- [ ] List all ssh keys loaded to `ssh-agent`
- [ ] Load key to `ssh-agent`
- [ ] Unload key from `ssh-agent`
- [ ] Generate the new ssh key pair
- [ ] Remove key pair
- [ ] Search by key name, comment

## License

[Apache 2.0](https://github.com/mixanemca/ssh-keys/raw/main/LICENSE)
