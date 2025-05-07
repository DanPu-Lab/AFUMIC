package cmd

import (
	"correct-go/src/core"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	familiesPath string
	barcodesPath string
	samPath      string
	outputPath   string
	pos          int
	mapq         int
	dist         int
	limit        int
	checkIds     bool
	// 输出文件路径
	namesToBarcodesPath string
	lostBarcodesPath    string
	barcodesTxtPath     string
	abTxtPath           string
	baTxtPath           string
	correctionsPath     string
)

var rootCmd = &cobra.Command{
	Use: "correct",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		output := make(chan string, 100)
		var wg sync.WaitGroup
		go func() {
			defer close(output)
			defer wg.Done()
			wg.Add(1)
			corrections, reversedBarcodes := mainProgress(start)
			core.GenerateCorrectedOutput(familiesPath, corrections, reversedBarcodes, output)
			fmt.Println("GenerateCorrectedOutput completed time: ", time.Since(start))
		}()

		outputFile, err := os.Create(outputPath)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		go func() {
			for line := range output {
				_, err = outputFile.WriteString(line + "\n")
				if err != nil {
					panic(err)
				}
			}
		}()
		wg.Wait()
		fmt.Printf("correct-go time: %v\n", time.Since(start))
	},
}

func init() {
	rootCmd.Flags().StringVarP(&familiesPath, "families", "f", "", "families file")
	rootCmd.MarkFlagRequired("families")
	rootCmd.Flags().StringVarP(&barcodesPath, "barcodes", "b", "", "barcodes file")
	rootCmd.MarkFlagRequired("barcodes")
	rootCmd.Flags().StringVarP(&samPath, "sam", "s", "", "sam file")
	rootCmd.MarkFlagRequired("sam")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "families.corrected.tsv", "output file")
	rootCmd.Flags().IntVarP(&pos, "pos", "p", 2, "position")
	rootCmd.Flags().IntVarP(&mapq, "mapq", "m", 20, "map quality")
	rootCmd.Flags().IntVarP(&dist, "dist", "d", 3, "distance")
	rootCmd.Flags().IntVarP(&limit, "limit", "l", 0, "limit")
	rootCmd.Flags().BoolVarP(&checkIds, "checkIds", "c", true, "check ids")

	rootCmd.Flags().StringVarP(&namesToBarcodesPath, "namesToBarcodes", "n", "", "names to barcodes file")
	rootCmd.Flags().StringVarP(&lostBarcodesPath, "lostBarcodes", "z", "", "lost barcodes file")
	rootCmd.Flags().StringVarP(&barcodesTxtPath, "barcodesTxt", "t", "", "barcodes txt file")
	rootCmd.Flags().StringVarP(&abTxtPath, "abTxt", "a", "", "ab txt file")
	rootCmd.Flags().StringVarP(&baTxtPath, "baTxt", "q", "", "ba txt file")
	rootCmd.Flags().StringVarP(&correctionsPath, "corrections", "r", "", "corrections file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainProgress(start time.Time) (map[string]string, map[string]struct{}) {
	var wg sync.WaitGroup
	namesToBarcodes := core.MapNamesToBarcodes(barcodesPath)
	fmt.Println("MapNamesToBarcodes completed time: ", time.Since(start))

	wg.Add(1)
	go func(barcodesMap map[int]string) {
		defer wg.Done()
		if len(namesToBarcodesPath) == 0 {
			return
		}

		file, err := os.Create(namesToBarcodesPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		for k, v := range barcodesMap {
			line := fmt.Sprintf("%d->%s\n", k, v)
			_, err = file.WriteString(line)
			if err != nil {
				panic(err)
			}
		}
	}(namesToBarcodes)

	passingAlignments := core.FilterAlignment(namesToBarcodes, samPath, lostBarcodesPath, pos, mapq, dist)
	fmt.Println("FilterAlignment completed time: ", time.Since(start))

	graph, reversedBarcodes, _ := core.ReadAlignments(passingAlignments, namesToBarcodes)
	fmt.Println("ReadAlignments completed time: ", time.Since(start))

	wg.Add(1)
	go func(alignments []core.AlignmentData, namesToBarcodes map[int]string) {
		defer wg.Done()
		if len(barcodesTxtPath) == 0 {
			return
		}

		file, err := os.Create(barcodesTxtPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		for _, data := range alignments {
			rseq := namesToBarcodes[data.RName]
			qseq := namesToBarcodes[data.QName]
			line := fmt.Sprintf("%s->%s\n", rseq, qseq)
			_, err = file.WriteString(line)
			if err != nil {
				panic(err)
			}
		}

	}(passingAlignments, namesToBarcodes)

	familyCounts, _ := core.GetFamilyCounts(familiesPath, limit, checkIds)
	fmt.Println("GetFamilyCounts completed time: ", time.Since(start))

	wg.Add(1)

	go func(familyCounts map[string]map[string]int) {
		defer wg.Done()
		if len(abTxtPath) == 0 || len(baTxtPath) == 0 {
			return
		}

		abFile, err := os.Create(abTxtPath)
		if err != nil {
			panic(err)
		}
		defer abFile.Close()

		baFile, err := os.Create(baTxtPath)
		if err != nil {
			panic(err)
		}
		defer baFile.Close()

		for k, v := range familyCounts {
			num := v["ab"]
			if num != 0 {
				_, err = abFile.WriteString(k + " ab " + fmt.Sprint(num) + "\n")
				if err != nil {
					panic(err)
				}
			} else {
				num = v["ba"]
				if num != 0 {
					_, err = baFile.WriteString(k + " ba " + fmt.Sprint(num) + "\n")
					if err != nil {
						panic(err)
					}
				}
			}

		}
	}(familyCounts)

	corrections := core.MakeCorrectionTable(graph, familyCounts)
	fmt.Println("MakeCorrectionTable completed time: ", time.Since(start))

	wg.Add(1)
	go func(corrections map[string]string) {
		defer wg.Done()
		if len(correctionsPath) == 0 {
			return
		}

		file, err := os.Create(correctionsPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		for k, v := range corrections {
			line := fmt.Sprintf("%s->%s\n", k, v)
			_, err = file.WriteString(line)
			if err != nil {
				panic(err)
			}
		}
	}(corrections)

	wg.Wait()
	return corrections, reversedBarcodes
}
