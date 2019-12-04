set -x
rm -rf /tmp/dqlite-demo
rm -f /tmp/dqdemo*

# fail at first error
set -e

export PATH=$PWD:$PATH
dqlite-demo start 1 > /tmp/dqdemo1.txt 2>&1 &
sleep 1
dqlite-demo start 2 > /tmp/dqdemo2.txt 2>&1 &
sleep 1
dqlite-demo start 3 > /tmp/dqdemo3.txt 2>&1 &
sleep 1
dqlite-demo add 2
sleep 1
dqlite-demo add 3
sleep 1
dqlite-demo update mykey myvalues
dqlite-demo query mykey
