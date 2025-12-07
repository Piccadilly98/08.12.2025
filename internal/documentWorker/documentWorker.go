package document_worker

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/signintech/gopdf"
)

const (
	startX  = 100.0
	startY  = 70.0
	cellW   = 200.0
	cellH   = 20.0
	finishY = 782
)

func CreateDocument(data map[int64]map[string]string) ([]byte, error) {
	font := ""
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err := pdf.AddTTFFont("SF", "../../materials/fonts/SFNS.ttf")
	if err != nil {
		log.Print(err.Error())
		err := pdf.AddTTFFont("NY", "../../materials/fonts/NewYork.ttf")
		if err != nil {
			return nil, err
		}
		font = "NY"
	} else {
		font = "SF"
	}
	pdf.AddPage()

	err = pdf.SetFont(font, "", 20)
	if err != nil {
		log.Print(err.Error())
		return nil, err

	}
	pdf.SetXY(230, 50)
	pdf.Text("LINKS REPORT")
	pdf.Br(30)
	pageCount := 1
	drawTableHeader(&pdf, pageCount, font)
	currentY := startY + cellH
	pdf.SetFont(font, "", 10)
	sort := getSortMapKey(data)
	for _, k := range sort {
		v := data[k]
		if currentY >= finishY {
			pageCount++
			pdf.AddPage()
			y := drawTableHeader(&pdf, pageCount, font)
			currentY = y + 20
		}
		drawNumberBucket(&pdf, &currentY, k, font)
		for k1, v1 := range v {
			pdf.SetXY(startX, currentY)
			if len(k1) >= 45 {
				k1 = k1[:20] + "..." + k1[len(k1)-10:]
			}
			pdf.CellWithOption(&gopdf.Rect{W: cellW, H: cellH}, k1, gopdf.CellOption{
				Align:  gopdf.Center,
				Border: gopdf.AllBorders,
			})
			pdf.SetXY(startX+cellW, currentY)
			pdf.CellWithOption(&gopdf.Rect{W: cellW, H: cellH}, v1, gopdf.CellOption{
				Align:  gopdf.Center,
				Border: gopdf.AllBorders,
			})

			currentY += cellH
		}
	}
	currentY += cellH * 2
	pdf.SetFont(font, "", 10)
	if currentY >= finishY {
		pdf.AddPage()
		pdf.SetXY(startX, startY-10)
	} else {
		pdf.SetXY(350, currentY)
	}
	finishStr := fmt.Sprintf("Document created in %s", time.Now().Format("2006-01-02 15:04:05"))
	pdf.Text(finishStr)
	res := new(bytes.Buffer)
	_, err = pdf.WriteTo(res)
	log.Printf("create document!Page count %d\n", pageCount)
	return res.Bytes(), err
}

func drawTableHeader(pdf *gopdf.GoPdf, pageCount int, font string) float64 {
	pdf.SetFont(font, "", 14)
	y := startY
	if pageCount > 1 {
		y -= 20
	}
	pdf.SetXY(startX, y)
	pdf.CellWithOption(&gopdf.Rect{W: cellW, H: cellH}, "Link", gopdf.CellOption{
		Align:  gopdf.Center,
		Border: gopdf.AllBorders,
	})

	pdf.SetXY(startX+cellW, y)
	pdf.CellWithOption(&gopdf.Rect{W: cellW, H: cellH}, "Status", gopdf.CellOption{
		Align:  gopdf.Center,
		Border: gopdf.AllBorders,
	})
	pdf.SetFont(font, "", 10)
	return y
}

func drawNumberBucket(pdf *gopdf.GoPdf, currentY *float64, bucketNum int64, font string) {
	pdf.SetFont(font, "", 12)
	pdf.SetXY(startX, *currentY)
	str := fmt.Sprintf("Bucket number: %d", bucketNum)
	pdf.CellWithOption(&gopdf.Rect{W: cellW * 2, H: cellH}, str, gopdf.CellOption{
		Align:  gopdf.Center,
		Border: gopdf.AllBorders,
	})
	*currentY += cellH
	pdf.SetFont(font, "", 10)
}

func getSortMapKey(data map[int64]map[string]string) []int64 {
	res := []int64{}
	for k := range data {
		res = append(res, k)
	}
	slices.Sort(res)
	return res
}
