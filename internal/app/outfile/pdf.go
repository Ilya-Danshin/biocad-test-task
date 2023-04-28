package outfile

import (
	"context"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"

	"test_task/internal/app/database"
)

type PDFFile struct {
	outFilesDir string
	db          database.IDatabase
}

func New(outFilesDir string, db database.IDatabase) (*PDFFile, error) {
	pdf := PDFFile{}

	pdf.outFilesDir = outFilesDir
	pdf.db = db

	err := license.SetMeteredKey(os.Getenv(`PDF_API_KEY`))
	if err != nil {
		return nil, err
	}

	return &pdf, nil
}

func (f *PDFFile) WriteData(ctx context.Context, records []database.Record) error {
	// sort slice
	sort.Slice(records, func(i, j int) bool {
		return records[i].UnitGuid.String() > records[j].UnitGuid.String()
	})

	// get all unique guid
	uniqGuids := getUniqueGUid(records)

	for _, guid := range uniqGuids {
		// get all records for guid
		allRec, err := f.db.GetRecordsByGuid(ctx, guid)
		if err != nil {
			return nil
		}

		// write all records in pdf file
		err = f.WriteToPdf(allRec)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *PDFFile) WriteToPdf(records []database.Record) error {
	c, err := createPdf(records)
	if err != nil {
		return err
	}

	// Write to output file.
	if err := c.WriteToFile(f.outFilesDir + "\\" + records[0].UnitGuid.String() + ".pdf"); err != nil {
		log.Fatal(err)
	}

	return nil
}

func getUniqueGUid(records []database.Record) []uuid.UUID {

	var uniqGuids []uuid.UUID
	var prevGuid uuid.UUID
	for _, record := range records {
		if record.UnitGuid != prevGuid {
			uniqGuids = append(uniqGuids, record.UnitGuid)
			prevGuid = record.UnitGuid
		} else {
			continue
		}
	}

	return uniqGuids
}

func createPdf(records []database.Record) (*creator.Creator, error) {
	// Create report fonts.
	font, err := model.NewStandard14Font("Helvetica")
	if err != nil {
		return nil, err
	}

	fontBold, err := model.NewStandard14Font("Helvetica-Bold")
	if err != nil {
		return nil, err
	}

	c := creator.New()
	pageSize := creator.PageSize{creator.PageSizeA4[1], creator.PageSizeA4[0]}
	c.SetPageSize(pageSize)

	table := c.NewTable(15)
	table.SetMargins(0, 0, 10, 0)

	// Draw table header.
	addCell(c, table, "n", fontBold)
	addCell(c, table, "mqtt", fontBold)
	addCell(c, table, "invid", fontBold)
	addCell(c, table, "unit_guid", fontBold)
	addCell(c, table, "msg_id", fontBold)
	addCell(c, table, "text", fontBold)
	addCell(c, table, "context", fontBold)
	addCell(c, table, "class", fontBold)
	addCell(c, table, "level", fontBold)
	addCell(c, table, "area", fontBold)
	addCell(c, table, "addr", fontBold)
	addCell(c, table, "block", fontBold)
	addCell(c, table, "type", fontBold)
	addCell(c, table, "bit", fontBold)
	addCell(c, table, "invert_bit", fontBold)

	for _, record := range records {
		addCell(c, table, strconv.Itoa(record.N), font)
		addCell(c, table, string(record.MQTT), font)
		addCell(c, table, record.InvId, font)
		addCell(c, table, record.UnitGuid.String(), font)
		addCell(c, table, record.MsgId, font)
		addCell(c, table, record.Text, font)
		addCell(c, table, string(record.Context), font)
		addCell(c, table, record.Class, font)
		addCell(c, table, strconv.Itoa(record.Level), font)
		addCell(c, table, record.Area, font)
		addCell(c, table, record.Addr, font)
		addCell(c, table, record.Block, font)
		addCell(c, table, record.Type, font)
		addCell(c, table, strconv.Itoa(record.Bit), font)
		addCell(c, table, strconv.Itoa(record.InvertBit), font)
	}

	err = c.Draw(table)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func addCell(c *creator.Creator, table *creator.Table, text string, font *model.PdfFont) *creator.TableCell {
	cell := table.NewCell()

	p := c.NewStyledParagraph()
	p.Append(text).Style.Font = font

	cell.SetContent(p)
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)

	return cell
}
