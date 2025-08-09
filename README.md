# AFUMIC #

Accurate detection of low-frequency DNA variants below 1% is critical. Here, we present AFUMIC, an alignment-free UMI clustering framework that systematically addresses these limitations through collision-resilient UMI grouping and a consensus quality score (CQS)–guided strategy for high-fidelity consensus sequence generation.

## Install ##

git clone --recursive https://github.com/DanPu-Lab/AFUMIC.git

##Dependencies

This pipeline requires the following dependencies:

| Program | Version | Purpose                                    |
| ------- | ------- | ------------------------------------------ |
| [python](https://www.python.org/)| 3.9 or higher | Perform multiple sequence alignment |  
| [Bowtie](http://bowtie-bio.sourceforge.net/) | 1.3.1   | Align sequence                |
| [BWA](http://bio-bwa.sourceforge.net/) | 0.7.17   | Align sequence to the reference genome                |
| [Umi-tools](https://github.com/CGATOxford/UMI-tool))| 1.1.4   | Extract UMI              |
| [Networkx](https://networkx.org)) | 2.8   | luster UMI                |
| [MAFFT](http://samtools.sourceforge.net/)| 7.505   | Multiple sequence alignment              |

## Running AFUMIC ##

### 1.	UMI clustering ###

Sequences with the same UMI are grouped into the same family. Run make-families.sh with required input parameters:

'''

read_1.fastq  FASTQ1,    FASTQ containing Read 1 of paired-end reads.

read_2.fastq  FASTQ2,    FASTQ containing Read 2 of paired-end reads.


'''

bash make-families.sh read_1.fastq read_2.fastq > families.tsv

DESCRIPTION: The sequences clustered by UMI are written into the families.tsv file.

### 2.	Create an index for UMI ###

bash baralign.sh families.tsv refdir barcodes.sam

DESCRIPTION: Among them, refdir is the generated barcode index folder, and barcodes.sam is the generated file. 

### 3.	Correct errors in UMI sequences ###

UMI sequences that are identical or highly similar are clustered into a single cluster, with errors in the UMIs corrected. To run run.sh, the following parameters need to be modified:

'''

output_file,     Sequence cluster files after alignment with the reference genome.

names_to_barcodes,     The ID of UMI. 

lost_barcodes,     UMI sequences with a similarity less than the Hamming distance.

barcodes_txt,     UMI sequences with a similarity more than the Hamming distance.

ab_txt,     The order of UMI.

ba_txt,     The order of UMI.

corrections,     The sequence file after clustering via run.sh

'''

bash run.sh > correct.txt

DESCRIPTION: Among them, you need to modify the path of the run.sh file. The families, barcodes, and sam files are generated in step 2. The subsequent files are output files, and their paths need to be modified as well.

### 4.	Multiple sequence alignment ###

Perform multiple sequence alignment on UMI-corrected sequence clusters.

python align-families.py families.corrected.tsv > families.msa.tsv

### 5.	Generate consensus sequences ###

Run CQS/run.sh required these parameters:

'''

-f file,   	  Input the file after multiple sequence alignment.

-anchor,     The number of bases upstream and downstream of the midpoint of the core region. The default parameter is set to 15.


'''

bash run.sh > CQS.txt

DESCRIPTION: Modify the file path in run.sh. The families.msa.tsv file is derived from the output file generated in Step 4. This script merges duplicate reads from the aligned files into single-strand consensus sequences (SSCS), which are then combined into double-strand consensus sequences (DCS).

## OVERVIEW  ##

<img src="https://github.com/DanPu-Lab/AFUMIC/blob/master/AFUMIC/Overview.jpg" width="50%" height="50%">


