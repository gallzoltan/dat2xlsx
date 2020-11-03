package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type telrow struct {
	date     string //time.Time
	time     string
	ext      string
	co       string
	number   string
	ring     string
	duration string
	cost     string
	acc      string
	cd       string
	dt       string
}

func excelTime(t string) float32 {
	var result float32
	var validStr = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`)
	if validStr.MatchString(t) {
		split := strings.Split(t, ":")
		hh, _ := strconv.Atoi(split[0])
		mm, _ := strconv.Atoi(split[1])
		ss, _ := strconv.Atoi(split[2])
		result = (0.041666 * float32(hh)) + (0.000694 * float32(mm)) + (0.00001157 * float32(ss))
	} else {
		result = 0.0
	}
	return result
}

func readTelDat(path string, fromdt string) ([]telrow, error) {
	var lines []telrow
	var validLine = regexp.MustCompile(`^\d{2}\/\d{2}\/\d{2}`)
	from, err := strconv.ParseUint(fromdt, 10, 32)

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if validLine.MatchString(line) {
			u, err := strconv.ParseUint(line[120:126], 10, 32)
			if err != nil {
				log.Fatalf("readLines: %s", err)
			}

			if u >= from {
				//ext, err := strconv.Parseint(line[15:23], 6, 64)
				rowdate := "20" + strings.Replace(line[:8], "/", ".", 2)
				rowring := strings.Replace(line[79:83], "'", ":", 1)
				rowdura := strings.Replace(line[84:95], "'", ":", 1)
				row := telrow{
					date:     strings.TrimSpace(rowdate),
					time:     strings.TrimSpace(line[9:14]),
					ext:      strings.TrimSpace(line[15:22]),
					co:       strings.TrimSpace(line[23:27]),
					number:   strings.TrimSpace(line[28:77]),
					ring:     rowring,
					duration: strings.TrimSpace(rowdura),
					cost:     strings.TrimSpace(line[96:104]),
					acc:      strings.TrimSpace(line[105:115]),
					cd:       strings.TrimSpace(line[116:119]),
					dt:       strings.TrimSpace(line[120:131]),
				}

				lines = append(lines, row)
			}
		}
	}

	return lines, scanner.Err()
}

func writeExcel(lines []telrow, xlsxFile string) {
	xlsx := excelize.NewFile()
	xlsx.SetSheetName(xlsx.GetSheetName(0), "teledat")
	hmm, err := xlsx.NewStyle(`{"number_format": 20}`)
	if err != nil {
		fmt.Println(err)
	}
	mmss, err := xlsx.NewStyle(`{"number_format": 45}`)
	if err != nil {
		fmt.Println(err)
	}
	hmmss, err := xlsx.NewStyle(`{"number_format": 46}`)
	if err != nil {
		fmt.Println(err)
	}

	xlsx.SetCellStr("teledat", "A1", "Date")
	xlsx.SetCellStr("teledat", "B1", "Time")
	xlsx.SetCellStr("teledat", "C1", "Ext")
	xlsx.SetCellStr("teledat", "D1", "CO")
	xlsx.SetCellStr("teledat", "E1", "Dial Number")
	xlsx.SetCellStr("teledat", "F1", "Ring")
	xlsx.SetCellStr("teledat", "G1", "Duration")
	xlsx.SetCellStr("teledat", "H1", "Cost")
	xlsx.SetCellStr("teledat", "I1", "ACC")
	xlsx.SetCellStr("teledat", "J1", "CD")
	xlsx.SetCellStr("teledat", "K1", "Date Time")
	for i, line := range lines {
		n := i + 2
		axis, _ := excelize.CoordinatesToCellName(1, n)
		xlsx.SetCellValue("teledat", axis, line.date)
		axis, _ = excelize.CoordinatesToCellName(2, n)
		xlsx.SetCellValue("teledat", axis, line.time)
		xlsx.SetCellStyle("teledat", axis, axis, hmm)

		axis, _ = excelize.CoordinatesToCellName(3, n)
		xlsx.SetCellStr("teledat", axis, line.ext)
		axis, _ = excelize.CoordinatesToCellName(4, n)
		xlsx.SetCellStr("teledat", axis, line.co)
		axis, _ = excelize.CoordinatesToCellName(5, n)
		xlsx.SetCellStr("teledat", axis, line.number)
		axis, _ = excelize.CoordinatesToCellName(6, n)
		xlsx.SetCellValue("teledat", axis, line.ring)
		xlsx.SetCellStyle("teledat", axis, axis, mmss)

		axis, _ = excelize.CoordinatesToCellName(7, n)
		//xlsx.SetCellValue("teledat", axis, line.duration)
		xlsx.SetCellValue("teledat", axis, excelTime(line.duration))
		xlsx.SetCellStyle("teledat", axis, axis, hmmss)

		axis, _ = excelize.CoordinatesToCellName(8, n)
		xlsx.SetCellStr("teledat", axis, line.cost)
		axis, _ = excelize.CoordinatesToCellName(9, n)
		xlsx.SetCellStr("teledat", axis, line.acc)
		axis, _ = excelize.CoordinatesToCellName(10, n)
		xlsx.SetCellStr("teledat", axis, line.cd)
		axis, _ = excelize.CoordinatesToCellName(11, n)
		xlsx.SetCellStr("teledat", axis, line.dt)
	}
	//if err := xlsx.SaveAs(".\\bin\\testfiles\\exportdat.xlsx"); err != nil {
	if err := xlsx.SaveAs(xlsxFile); err != nil {
		log.Fatalf("readLines: %s", err)
	}
}

func main() {
	log.Println("The program has started...")
	now := time.Now()
	beginDate := now.AddDate(0, -1, 0)
	fromPtr := flag.String("begin", beginDate.Format("060102"), "a megadott yymmdd dátumtól kezdje el az exportálást")
	datFile := flag.String("in", "./tele.dat", "a tele.dat fájl helye")
	xlsxFile := flag.String("out", "./datexport.xlsx", "az xlsx export fájl helye")
	flag.Parse()

	//lines, err := readTelDat(".\\bin\\testfiles\\tele.dat", *fromPtr)
	lines, err := readTelDat(*datFile, *fromPtr)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	writeExcel(lines, *xlsxFile)
	log.Println("The program has finished running.")
}
