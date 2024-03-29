# azctx

This program let's you pick an Azure subscription interactively with the output of `az login` or `az account list` commands. 

Tip: this program can be replaced with the following commands. Requires `jq`, `fzf`, `awk` and `xargs`.

```bash
az account login
az account list --query '[].{Name:name, ID:id}' --output json | jq -r '.[] | "\(.Name) \(.ID)"' | fzf | awk '{print $NF}' | xargs -I {} az account set --subscription {}
```

## Usage

```bash
az login | pick
az account list | pick
```

## Installation

```bash
go install github.com/kimfy/azctx
```

## License

GNU GPLv3

