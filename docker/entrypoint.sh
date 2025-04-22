#!/bin/bash

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
    if [ "$1" != 'bash' ]; then
        drycc "$@"
    else
        exec "$@"
    fi
}

_main "$@"