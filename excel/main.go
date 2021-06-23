package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

type mento struct {
	info     string
	priority int
}

func main() {
	Mentos := make(map[int]*mento)
	i := 0
	xlFile, err := xlsx.OpenFile("./mentoList.xlsx")
	if err != nil {
		panic(err)
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			line := ""
			for _, cell := range row.Cells {
				text := cell.String()
				text += "//"
				if text != "" {
					line += text
				}
			}

			// IT+투자+IR+비즈니스모델
			if !strings.Contains(line, "의료") && !strings.Contains(line, "자재관리") &&
				!strings.Contains(line, "자율 주행차") && !strings.Contains(line, "공정") &&
				!strings.Contains(line, "스마트 팜") && !strings.Contains(line, "3D") &&
				!strings.Contains(line, "드론") && !strings.Contains(line, "섬유") && !strings.Contains(line, "로봇") &&
				!strings.Contains(line, "전시회") && !strings.Contains(line, "AR/VR") && !strings.Contains(line, "전기/전자") &&
				!strings.Contains(line, "기술전문가") && !strings.Contains(line, "전기/전자") &&
				!strings.Contains(line, "기계") && !strings.Contains(line, "공예") && !strings.Contains(line, "시제품") &&
				!strings.Contains(line, "유통채널") && !strings.Contains(line, "제품") &&
				!strings.Contains(line, "창업이음") && !strings.Contains(line, "해외 마케팅") &&
				!strings.Contains(line, "사업자등록") && !strings.Contains(line, "에 너 지") {
				mento := new(mento)
				mento.info = line
				slice := strings.Split(line, "//")

				key, err := strconv.Atoi(slice[0])
				if err == nil {
					Mentos[key] = mento
				}
				// fmt.Println(line)
				// fmt.Println()
				i++
			}
		}
	}
	// 우선순위
	keywords := []string{"벤처투자", "엔젤투자", "엑셀러레이터"}
	sortedMentos := []*mento{}

	for _, val := range Mentos {
		pri := 0
		for _, keyword := range keywords {
			if strings.Contains(val.info, keyword) {
				pri++
			}
		}
		val.priority = pri
		// fmt.Println(val.priority)
		sortedMentos = append(sortedMentos, val)
	}

	sort.Slice(sortedMentos, func(i, j int) bool {
		return sortedMentos[i].priority > sortedMentos[j].priority
	})

	// for _, val := range sortedMentos {
	// 	fmt.Println("priority: ", val.priority, "info: ", val.info)
	// }
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("멘토명단")
	if err != nil {
		panic(err)
	}

	category := []string{"우선순위", "연번", "구분", "전담멘토명", "출생년도", "소속", "직위", "분야", "직군"}
	row = sheet.AddRow()
	for _, v := range category {
		cell = row.AddCell()
		cell.Value = v
	}

	for _, val := range sortedMentos {
		row = sheet.AddRow()
		slice := strings.Split(val.info, "//")
		val.priority = 4 - val.priority
		pri := strconv.Itoa(val.priority)
		slice = append([]string{pri}, slice...)
		// fmt.Println(pri)
		for _, section := range slice {
			cell = row.AddCell()
			cell.Value = section
		}
	}
	err = file.Save("멘토명단.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Println(i)
}
