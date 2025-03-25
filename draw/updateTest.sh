#!/bin/bash

script="testdata/draw.txt"

cat > "$script" <<-"HEADER"
# width values:
# 1550 (no break at all)
# 50 (break everywhere)
# 750 (break short after bigMerge and long after Split1 and short after Split2)
# 850 (break long after bigMerge)
# 1150 (break before Split1 and Split2)
# 1250 (break after Split1 and before Split2)
# 1350 (break long and short before lastMerge)
#
# draw the BigTestFlowData and compare the resulting SVG with the expected one:
drawBigTestFlowData false false 1550
cmp markdown-false-false-1550.actual markdown-false-false-1550.expected
cmp flowdev/flow-bigTestFlow1550.svg flowdev/flow-bigTestFlow1550.expected
# draw the BigTestFlowData split up in many SVGs and in dark mode:
drawBigTestFlowData true true 750
cmp markdown-true-true-750.actual markdown-true-true-750.expected
HEADER

for fnam in $(basename -a -s .svg flowdev/*.svg | sort) ; do
	echo "cmp flowdev/$fnam.svg flowdev/$fnam.expected" >> "$script"
done

echo "" >> "$script"
echo "" >> "$script"

echo "-- markdown-false-false-1550.expected --" >> "$script"
cat "./markdown-false-false-1550.actual.md" >> "$script"

echo "-- markdown-true-true-750.expected --" >> "$script"
cat "./markdown-true-true-750.actual.md" >> "$script"

for fnam in $(basename -a -s .svg flowdev/*.svg | sort) ; do
	echo "-- flowdev/$fnam.expected --" >> "$script"
	cat "flowdev/$fnam.svg" >> "$script"
done
