# gdrift-backend
Repository for the REST services of Gdrift's backend

Configuration
-------------------
```
An appropriate mysql installation must be available
A database example can be found in `db/example.sql`
Configuration file `config/config.yaml` must be prepared locally
```

Install
-------------------
```
sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update
sudo apt-get install glide

mkdir ~/go/src/github/vsepulve/
cd ~/go/src/github/vsepulve/
git clone https://github.com/vsepulve/gdrift-backend
cd gdrift-backend
go install

$GOPATH/bin/gdrift-backend
```
