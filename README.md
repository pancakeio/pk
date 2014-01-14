# pk

`pk` is a command-line tool for interacting with the Pancake.io API.

## Usage

```sh
pk add-key [-key-path <path to ssh key>]
pk list-keys
pk remove-key

pk create-project [-static <make project static>]
pk list-projects
pk delete-project

pk help [command-name]
```

## Completion

Add to `~/.bashrc`

```
complete -W "$($GOPATH/bin/pk -w)" pk
```