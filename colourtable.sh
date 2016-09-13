#!/bin/sh

CSI="["
for i in `seq 7`
do
	echo "3${i} ${CSI}3${i}mFG${CSI}m — 4${i} ${CSI}4${i}mBG${CSI}m"
done
