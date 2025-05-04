#!/bin/bash

install () {
    echo "Installing schema cli"
    git clone https://github.com/Ayobami6/schema_dump.git
    chmod u+x ./schema_dump/bin/linux/schema
    # get the os name
    os_name=$(uname)
    if [[ $os_name == "Darwin" ]]; then
        chmod u+x ./schema_dump/bin/mac/schema
        sudo cp ./schema_dump/bin/mac/schema /usr/local/bin
    elif [[ $os_name == "Linux" ]]; then
        sudo cp ./schema_dump/bin/linux/schema /usr/local/bin
    fi

    # clean up
    rm -rf ./schema_dump
}

install 