package data

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/lanrat/extsort"
	"github.com/lanrat/extsort/queue"
	log "github.com/sirupsen/logrus"
)

const ValidPromosFile = "couponbase_validpromos"
const DirPrefix = "data/"

const writerBufferSize = 100 * 1024 * 1024
const sortedPromosSuffix = "_sorted"

type promoCode struct {
	code string
}

func (p promoCode) ToBytes() []byte {
	return []byte(p.code)
}

func promoCodeFromBytes(b []byte) extsort.SortType {
	return promoCode{code: string(b)}
}

func comparePromoCodeLess(a, b extsort.SortType) bool {
	return a.(promoCode).code < b.(promoCode).code
}

type heapNode struct {
	promoCode string
	fileIndex int
}

var config = extsort.Config{
	ChunkSize:          1000000,
	NumWorkers:         4,
	ChanBuffSize:       100,
	SortedChanBuffSize: 1000,
	TempFilesDir:       "",
}

func isValidPromoCode(code string) bool {
	length := len(code)
	return length >= 8 && length <= 10
}

func sortFile(inputPath, outputPath string) error {
	inputChan := make(chan extsort.SortType)
	go func() {
		file, err := os.Open(inputPath)
		if err != nil {
			close(inputChan)
			log.Panicf("Error opening file %s: %v", inputPath, err)
			return
		}
		defer file.Close()

		gzReader, err := gzip.NewReader(file)
		if err != nil {
			close(inputChan)
			log.Panicf("Error creating gzip reader: %v", err)
			return
		}
		defer gzReader.Close()

		scanner := bufio.NewScanner(gzReader)
		for scanner.Scan() {
			line := scanner.Text()
			if isValidPromoCode(line) {
				inputChan <- promoCode{code: line}
			}
		}
		close(inputChan)
	}()

	sorter, outputChan, errChan := extsort.New(inputChan, promoCodeFromBytes, comparePromoCodeLess, &config)
	sorter.Sort(context.Background())

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	writer := bufio.NewWriterSize(outputFile, writerBufferSize)
	defer writer.Flush()

	var prior string
	for data := range outputChan {
		code := data.(promoCode).code
		if prior == code {
			continue
		}
		fmt.Fprintln(writer, code)
		prior = code
	}

	log.Infof("Sorted input file: %s and wrote to %s", inputPath, outputPath)
	return <-errChan
}

func mergeSortedFiles(tempFiles []string, outputFile string) error {
	readers := make([]*bufio.Scanner, len(tempFiles))
	files := make([]*os.File, len(tempFiles))

	for i, file := range tempFiles {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		files[i] = f

		readers[i] = bufio.NewScanner(f)

		defer f.Close()
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriterSize(outFile, writerBufferSize)
	defer writer.Flush()

	priorityQ := queue.NewPriorityQueue(func(a, b interface{}) bool {
		return a.(*heapNode).promoCode < b.(*heapNode).promoCode
	})

	// Insert initial lines
	for i, reader := range readers {
		if reader.Scan() {
			priorityQ.Push(&heapNode{promoCode: reader.Text(), fileIndex: i})
		}
	}

	var prior string
	count := 0

	for priorityQ.Len() > 0 {
		node := priorityQ.Pop().(*heapNode)

		// Only write duplicates
		if node.promoCode == prior {
			count++
		} else {
			if count > 1 {
				fmt.Fprintln(writer, prior)
			}

			prior = node.promoCode
			count = 1
		}

		// Pull next code from same file that we popped off the PQ earlier
		fileIndex := node.fileIndex
		if readers[fileIndex].Scan() {
			priorityQ.Push(&heapNode{promoCode: readers[fileIndex].Text(), fileIndex: fileIndex})
		}
	}

	// Write remaining duplicated lastCode after we finish
	if count > 1 {
		fmt.Fprintln(writer, prior)
	}

	log.Infof("Duplicate promo codes written to: %s", outputFile)
	return nil
}

func WriteValidPromos() error {
	start := time.Now()

	inputFiles := []string{"data/couponbase1.gz", "data/couponbase2.gz", "data/couponbase3.gz"}
	sortedFiles := make([]string, 0, len(inputFiles))
	for _, input := range inputFiles {
		file := input + sortedPromosSuffix
		sortedFiles = append(sortedFiles, file)
	}

	var wg sync.WaitGroup
	for i, input := range inputFiles {
		wg.Add(1)
		go func(i int, input string) {
			defer wg.Done()
			if err := sortFile(input, sortedFiles[i]); err != nil {
				log.Panicf("Error sorting file %s: %v", input, err)
			}
		}(i, input)
	}

	wg.Wait()

	err := mergeSortedFiles(sortedFiles, DirPrefix+ValidPromosFile)
	if err != nil {
		return err
	}

	for _, file := range sortedFiles {
		os.Remove(file)
	}

	duration := time.Since(start)
	log.Info("Producing valid codes took: ", duration)

	return nil
}
