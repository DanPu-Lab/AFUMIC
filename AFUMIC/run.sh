#!/usr/bin/env bash

#例如
families=/路径/families.tsv
barcodes=/路径/refdir/barcodes.fa
sam=/路径/barcodes.sam

output_file=/路径/families.corrected.tsv

names_to_barcodes=/路径/names_to_barcodes.csv
lost_barcodes=/路径/lost_barcodes.txt
barcodes_txt=/路径/barcodes.txt
ab_txt=/路径/b.txt
ba_txt=/路径/ba.txt
corrections=/路径/corrections.txt

./correct -b $barcodes -f $families -s $sam -o $output_file \
	-n $names_to_barcodes \
	-z $lost_barcodes \
	-t $barcodes_txt \
	-a $ab_txt \
	-q $ba_txt \
	-r $corrections
	#--path $corrected_path \
	#--graphOptimize

wc -l $output_file
