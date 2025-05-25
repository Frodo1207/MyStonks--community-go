## BUILD
CGOENABLED=0 GOOS=linux GOARCH=amd64 go build -o MyStonksDao

## RUN 
./MyStonksDao start --config=./config.yaml

## GEN DOCS
rm -rf docs

swag init

## DOCS 
http://localhost:8000/swagger/index.html


## Deploy
CGOENABLED=0 GOOS=linux GOARCH=amd64 go build -o MyStonksDao

cp MyStonksDao $HOME/deploy

cp ./config/config.yaml $HOME/config.yaml

replace ENV in ./scripts/mystonksdao.service.tmpl

sudo cp ./scripts/mystonksdao.service.tmpl /etc/systemd/system/mystonksdao.service

sudo systemctl start mystonksdao.service