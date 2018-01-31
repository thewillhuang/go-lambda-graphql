cd ./webapp && rm -fr build && yarn && cd .. && sh rebuild.sh
gin --all run main.go & cd ./webapp && yarn start && fg
