## ScalingoTest: Get the 100 last gtihub repositories created

### Installation
* Install go.
```
sudo apt-get install go
```
* Create a workspace directory and set _GOPATH_ accordingly.
```
mkdir $HOME/path/to/your/repository
export GOPATH=$HOME/path/to/your/repository
```
* Add the workspace's _bin_ subdirectory to your _PATH_:
```
export PATH=$PATH:$GOPATH/bin
```
* Get the __ScalingoTest project__.
```
go get github.com/utay/ScalingoTest
```
* Go automatically create a ScalingoTest binary in _GOPATH/bin_.

### How to use it?
* Launch the server.
```
cd $GOPATH
./bin/ScalingoTest
```
* Connect your browser to __localhost:4242__
![Alt text](https://github.com/utay/ScalingoTest/blob/master/images/search.png?raw=true "Search")
* Search GitHub repositories by name in the 100 last created
![Alt text](https://github.com/utay/ScalingoTest/blob/master/images/results.png?raw=true "Results")
