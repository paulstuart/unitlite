rm -rf /tmp/dqlite-demo
dqlite-demo start 1 &
dqlite-demo start 2 &
dqlite-demo start 3 &
dqlite-demo add 2
dqlite-demo add 3
dqlite-demo update mykey myvalues
dqlite-demo query mykey
#dqlite-demo start 4 &
#dqlite-demo add 4 
#dqlite-demo update key1 oneval
