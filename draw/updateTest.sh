#!/bin/bash

script="testdata/draw.txt"

cat > "$script" <<-"HEADER"
# draw the BigTestFlowData and compare the resulting SVG with the expected one:
drawBigTestFlowData false false
cmp markdown-false-false.actual markdown-false-false.expected 
cmp flowdev/flow-bigTestFlow.svg flowdev/flow-bigTestFlow.expected
# draw the BigTestFlowData split up in many SVGs and in dark mode:
drawBigTestFlowData true true
cmp markdown-true-true.actual markdown-true-true.expected 
HEADER

for fnam in $(basename -a -s .svg flowdev/flow-bigTestFlow-*.svg | sort) ; do
	echo "cmp flowdev/$fnam.svg flowdev/$fnam.expected" >> "$script"
done

echo "" >> "$script"
echo "" >> "$script"

echo "-- markdown-false-false.expected --" >> "$script"
cat "./markdown-false-false.actual.md" >> "$script"

echo "-- markdown-true-true.expected --" >> "$script"
cat "./markdown-true-true.actual.md" >> "$script"

for fnam in $(basename -a -s .svg flowdev/*.svg | sort) ; do
	echo "-- flowdev/$fnam.expected --" >> "$script"
	cat "flowdev/$fnam.svg" >> "$script"
done
