#!/bin/bash

# file=/mnt/yufei/method/dunovo/SC8/our/families.msa.tsv
#file=/mnt/yufei/method/dunovo/ABL1/our/families.msa.tsv
# file=/mnt/yufei/method/Calib/dunovo-8/families.msa.tsv
# file=/mnt/yufei/method/graduate/results/N0015/families.msa.tsv
#file=/mnt/yufei/keti/test/SRR23581200/families.msa.tsv
#file=/mnt/yufei/keti/ABL1/families.msa.tsv
file=/mnt/yufei/keti/SC8/families.msa.tsv
#file=/mnt/yufei/method/graduate/Germline/SRR14055021/families.msa.tsv
#file=/mnt/yufei/method/graduate/CQS/consensus-go-lib/HD701/anchor_15/HD701_without_clustering/families.msa.tsv
# file=/mnt/yufei/method/dunovo/GERMLINE/new_folder/families.msa.tsv
# file=/mnt/nfs/yufei/sinoduplex_data/code/dunovo/work/families.msa.tsv
# file=/mnt/nfs/yufei/sinoduplex_data/code/graduate/sim_UMIGen_results/families.msa.tsv 
# file=/mnt/yufei/method/graduate/dunovo/sim_UMIGen_results/families.msa.tsv
# file=/mnt/yufei/method/simulate/UMIGen/simulate/150_data/1000X/processed/families.msa.tsv
# file=/mnt/yufei/method/graduate/results/sim_UMIGen_1500x/without_UMI_correction/families.msa.tsv
#file=/mnt/yufei/method/graduate/results/HD734/SRR2556944/families.msa.tsv
# file=/mnt/yufei/method/compare/umi_clustering_compare/families.msa.tsv

./consensus -f $file -sscs1 ./sscs_1.fq -sscs2 ./sscs_2.fq -dcs1 ./dcs_1.fq -dcs2./dcs_2.fq -anchor 14


#./consensus -f $file -sscs1 /mnt/yufei/keti/sampleB/SRR13224665/sscs_1.fq -sscs2 /mnt/yufei/keti/sampleB/SRR13224665/sscs_2.fq -dcs1 /mnt/yufei/keti/sampleB/SRR13224665/dcs_1.fq -dcs2 /mnt/yufei/keti/sampleB/SRR13224665/dcs_2.fq -anchor 14

wc -l dcs_1.fq
wc -l dcs_2.fq
wc -l sscs_1.fq
wc -l sscs_2.fq

