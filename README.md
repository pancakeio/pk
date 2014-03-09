# pk

`pk` is a command-line tool for interacting with Pancake.

## Usage

```sh
pk login

pk add-key [-key-path <path to ssh key>]
pk list-keys
pk remove-key

pk create-project [ -dropbox=true (default: false) ]
pk list-projects
pk delete-project

pk help [command-name]
```

## Completion

Add to `~/.bashrc`

```
complete -W "$($GOPATH/bin/pk -w)" pk
```