# 1. Check if go is installed, install if not
if ! command -v go &> /dev/null
then
    echo "go command could not be found"
    echo "refer to https://go.dev/doc/install on installing go"
    exit 1
else
    go version
fi

# 2. Check if mysql is installed, install if not
if ! command -v mysql &> /dev/null
then
    echo "mysql command could not be found"
    echo "if you're on mac, use brew install mysql "
    echo "if you're on ubuntu, use sudo apt install mysql-server"
    exit 1
else
    mysql --version
fi