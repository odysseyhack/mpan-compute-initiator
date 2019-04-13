#!/bin/sh

# start the compute initiators (blockchain (-> CI -> mpc node) x5)

/home/mpc/compute-initiator/main -filename /tmp/nodepipe1&
/home/mpc/compute-initiator/main -filename /tmp/nodepipe2&
/home/mpc/compute-initiator/main -filename /tmp/nodepipe3&
/home/mpc/compute-initiator/main -filename /tmp/nodepipe4&
/home/mpc/compute-initiator/main -filename /tmp/nodepipe5&

/home/mpc/gabriele/node.py /tmp/nodepipe1&
/home/mpc/gabriele/node.py /tmp/nodepipe2&
/home/mpc/gabriele/node.py /tmp/nodepipe3&
/home/mpc/gabriele/node.py /tmp/nodepipe4&
/home/mpc/gabriele/node.py /tmp/nodepipe5&
