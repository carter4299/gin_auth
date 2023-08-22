#!/bin/bash

LOG_FILE="./my_server_config/startup/startup.log"
ENCRYPTED_FILE="./my_server_config/startup/encrypted.enc"
DECRYPTED_FILE="./my_server_config/startup/_/decrypted.txt"
ENV_FILE="./my_server_config/startup/_/.env"
GO_M_FILE="./my_server_config/startup/_/go_mod.txt"
PY_INIT="./my_server_config/startup/_/init.py"
INIT_DB="./my_server_config/startup/_/r0.txt"
SET_MAIN="./my_server_config/startup/_/r1.txt"
TEST_API_S0="./my_server_config/startup/_/s0.txt"
TEST_API_S1="./my_server_config/startup/_/s1.txt"
TEST_API_C0="./my_server_config/startup/_/c0.txt"
TEST_API_C1="./my_server_config/startup/_/c1.txt"
TEST_API_M0="./my_server_config/startup/_/m0.txt"
TEST_API_M1="./my_server_config/startup/_/m1.txt"
SERVER_FILE="./my_server/server.go"
CORS_FILE="./my_server_config/cors.go"
MIDDLEWARE_FILE="./middleware/middleware.go"
MAIN_RUNNER="main.go"

log() {
    echo "$(date +"%Y-%m-%d %H:%M:%S"): $1" >> $LOG_FILE
}

check_dependencies() {
    dependencies=("openssl" "yq" "go" "python3")
    
    log "Checking dependencies (openssl, yq, go, python3)..."

    for dep in "${dependencies[@]}"; do
        if [ "$dep" == "go" ] && [ -x ~/go/bin/go ]; then
            export PATH=$PATH:~/go/bin
            continue
        fi

        command -v $dep >/dev/null 2>&1 || {
            if [ "$dep" == "go" ]; then
                read -p "If $dep is installed. Would you like to update your GO_PATH (y/n): " choice
                if [ "$choice" == "y" ]; then
                    read -p "Enter your GOPATH: " gopath
                    echo "export GOPATH=$gopath" >> ~/.bashrc
                    echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bashrc
                    source ~/.bashrc
                    check_dependencies
                else
                    echo "The script requires $dep to run. Please install it and rerun the script."
                    exit 1
                fi
            fi
            read -p "$dep is not installed. Would you like to install it now? (y/n): " choice
            if [ "$choice" == "y" ]; then
                sudo apt-get install $dep
            else
                echo "The script requires $dep to run. Please install it and rerun the script."
                exit 1
            fi
        }
    done
}

encrypt_data() {
    log "Encrypting data..."
    local data="$1"
    local password="$2"
    echo -n "$data" | openssl enc -aes-256-cbc -salt -out "$ENCRYPTED_FILE" -k "$password" -pbkdf2
}

decrypt_file() {
    log "Decrypting file..."
    openssl enc -d -aes-256-cbc -in "$ENCRYPTED_FILE" -out "$DECRYPTED_FILE" -k "$1" -pbkdf2
}

read_and_set_env_vars() {
    log "Reading and setting environment variables from $1..."
    while IFS= read -r line; do
        export "$line"
    done < "$1"
}

prompt_for_encryption_password() {
    read -s -p "Enter a password to encrypt env: " encryption_password
    echo ""
    read -s -p "Confirm password: " confirm_password
    echo ""
    if [ "$encryption_password" != "$confirm_password" ]; then
        echo "Passwords do not match. Please try again."
        prompt_for_encryption_password
    fi
    encrypt_data "$(cat $ENV_FILE)" "$encryption_password"
    rm "$ENV_FILE"
}

prompt_for_decryption_password() {
    read -s -p "Enter your password: " encryption_password
    echo ""
    decrypt_file "$encryption_password"
    read_and_set_env_vars "$DECRYPTED_FILE"
    rm "$DECRYPTED_FILE"
}

init_go() {
    log "Initializing go..."
    local replacements=(
        "s/carter4299/${github_username}/g"
        "s/gin_auth/${github_repo}/g"
    )
    for r in "${replacements[@]}"; do
        find . -type f -not -path "./test_api/*" -not -path "./load.sh" -not -path "./readme.md" -exec sed -i "$r" {} +
    done
    
    sleep 5

    go mod init github.com/$github_username/$github_repo
    local go_packages=(
        "github.com/gin-gonic/gin"
        "github.com/mattn/go-sqlite3"
        "golang.org/x/crypto/bcrypt"
        "github.com/sirupsen/logrus"
        "github.com/gin-contrib/cors"
        "github.com/githubnemo/CompileDaemon"
    )
    for p in "${go_packages[@]}"; do
        go get $p
    done

    sleep 2
    go mod tidy
}
prompt_server_choice() {
    echo "You can remove these prompts of code from load_my_server.sh ( line 162 )."
    read -p "Would you like to test the API? (y/n): " u_choice
    if [ "$u_choice" == "y" ]; then
        cp $TEST_API_S1 $SERVER_FILE
        cp $TEST_API_C1 $CORS_FILE
        cp $TEST_API_M1 $MIDDLEWARE_FILE
    else
        cp $TEST_API_S0 $SERVER_FILE
        cp $TEST_API_C0 $CORS_FILE
        cp $TEST_API_M0 $MIDDLEWARE_FILE
    fi
    sleep 2
}
cleanup() {
    pkill -f CompileDaemon
    echo "Cleaned up processes."
}
cleanup_py() {
    pkill -f python3
    read -p "Would you like to (0) - exit or (1) - restart the python script ? (0/1): " choice
    case $choice in
        0)
            exit 1
            ;;
        1)
            python3 $PY_INIT
            ;;
        *)
            echo "Please enter 0 or 1"
            cleanup_py
            ;;
    esac
}

trap cleanup SIGINT
prompt_choice() {
    read -p "Would you like to use CompileDaemon for live updates? (y/n): " choice
    echo $choice

    case $choice in
        y)
            go get github.com/githubnemo/CompileDaemon
            go install github.com/githubnemo/CompileDaemon
            sleep 2
            CompileDaemon -command="./$github_repo"
            ;;
        n)
            go run .
            ;;
        *)
            echo "Please enter y or n"
            prompt_choice
            ;;
    esac
}
trap - INT
# Main script
if [ -e "$LOG_FILE" ]; then
    touch $LOG_FILE
fi

check_dependencies

if [ -e "$ENCRYPTED_FILE" ]; then
    log "Returning user detected..."
    prompt_server_choice
    source $GO_M_FILE
    prompt_for_decryption_password
    env | grep -E '(^KEY=|^SECRET=)'
    sleep 2
    prompt_choice
else
    log "First-time setup..."
    prompt_server_choice
    cp $INIT_DB $MAIN_RUNNER
    trap cleanup_py SIGINT
    python3 $PY_INIT
    trap - INT
    sleep 2

    if [ -e "$ENV_FILE" ]; then
        prompt_for_encryption_password
    fi

    source $GO_M_FILE
    init_go
    prompt_for_decryption_password
    env | grep -E '(^KEY=|^SECRET=)'
    sleep 2

    go run .

    cp $SET_MAIN $MAIN_RUNNER
    sleep 2
    prompt_choice
fi
