rm -fr node_modules
sh rebuild.sh
cd ./models && go get -u -t
cd ..
go test -cover -v ./models
yarn
