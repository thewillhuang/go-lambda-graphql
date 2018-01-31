# setup
```bash
brew install postgres95
brew install entr
brew install dep
sh ./init.sh
```

# watch and tests on file change
```bash
ls **/*.go **/*.gql | grep -v -E "vendor|models" | entr sh -c 'clear && date +"%T" && go test -v ./graphql/... . ./services/...'
```

# dev mode
```bash
sh dev.sh
```
