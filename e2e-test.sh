#!/usr/bin/env bash

set -euf -o pipefail

cd "$(dirname $0)"

port="$((10000 + (RANDOM % 10000)))"
tmpdir=$(mktemp -d /tmp/aws_ri_exporter_e2e_test.XXXXXX)
fixture='fixtures/e2e-output.txt'

if [ ! -x ./aws_ri_exporter ]; then
    echo './aws_ri_exporter not found. Consider running `go build` first.' >&2
    exit 1
fi

./aws_ri_exporter \
    --web.listen-address "127.0.0.1:${port}" \
    --log.level="debug" >"${tmpdir}/aws_ri_exporter.log" 2>&1 &

echo $! >"${tmpdir}/aws_ri_exporter.pid"

finish() {
    if [ $? -ne 0 -o ${verbose} -ne 0 ]; then
        cat <<EOF >&2
LOG =====================
$(cat "${tmpdir}/aws_ri_exporter.log")
=========================
EOF
    fi

    if [ ${update} -ne 0 ]; then
        cp "${tmpdir}/e2e-output.txt" "${fixture}"
    fi

    if [ ${keep} -eq 0 ]; then
        kill -9 "$(cat ${tmpdir}/aws_ri_exporter.pid)"
        # This silences the "Killed" message
        set +e
        wait "$(cat ${tmpdir}/aws_ri_exporter.pid)" >/dev/null 2>&1
        rm -rf "${tmpdir}"
    fi
}

trap finish EXIT

get() {
    if command -v curl >/dev/null 2>&1; then
        curl -s -f "$@"
    elif command -v wget >/dev/null 2>&1; then
        wget -O - "$@"
    else
        echo "Neither curl nor wget found"
        exit 1
    fi
}

sleep 1

get "127.0.0.1:${port}/metrics" | grep -E -v "${skip_re}" >"${tmpdir}/e2e-output.txt"

diff -u \
    "${fixture}" \
    "${tmpdir}/e2e-output.txt"
