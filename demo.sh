set -x
rm -rf /tmp/dqlite-demo

export PATH=$PWD:$PATH
dqlite-demo start 1 &
sleep 1
dqlite-demo start 2 &
sleep 1
dqlite-demo start 3 &
sleep 1
dqlite-demo add 2
sleep 1
dqlite-demo add 3
sleep 1
dqlite-demo update mykey myvalues
dqlite-demo query mykey
