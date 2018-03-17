#!/bin/bash

for BOLD in `seq 0 1`
do
	echo "==== BOLD=${BOLD} ===="
	echo "FG↓   BG→  0   1   2   3   4   5   6   7"
	for FG in `seq 0 7`
	do
		echo -ne " ${FG}   \e[m \e[${BOLD};3${FG}m X \e[m"
		for BG in `seq 0 7`
		do
			echo -ne "\e[m \e[${BOLD};3${FG};4${BG}m X \e[m"
		done
		echo ""
	done
done

echo "Foreground: \\033[3?m"
echo "Background: \\033[4?m"
