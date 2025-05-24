## BUILD
CGOENABLED=0 GOOS=linux GOARCH=amd64 go build -o mystonks 

## RUN 
./mystonks start --config=./config.yaml

## GEN DOCS
rm -rf docs
swag init

## DOCS 
http://localhost:8000/swagger/index.html
