#!/usr/local/bin/bash

if [ -n "$DRYCC_USER" ] && [ -n "$DRYCC_TOKEN" ] && [ -n "$DRYCC_ENDPOINT" ]; then
    mkdir -p ~/.drycc
    if [ ! -f ~/.drycc/client.json ]; then
        cat > ~/.drycc/client.json <<EOF
{
    "username": "${DRYCC_USER}",
    "ssl_verify": ${DRYCC_SSL_VERIFY:-true},
    "controller": "${DRYCC_ENDPOINT}",
    "token": "${DRYCC_TOKEN}",
    "response_limit": ${DRYCC_RESPONSE_LIMIT:-0}
}
EOF
    fi
fi

_main() {
    if [ "$1" == 'bash' ]; then
        shift
        exec bash "$@"
    elif [ "$1" == 'wait' ]; then
        shift
        /etc/bash_completion.d/timeout
        while [ "$(date +%s)" -le "$(flock -x -w 10 /tmp/timeout.lock cat /etc/wait/timeout)" ]; do
            sleep ${TIMEOUT:-30}
        done
    else
        exec drycc "$@"
    fi
}

_main "$@"
