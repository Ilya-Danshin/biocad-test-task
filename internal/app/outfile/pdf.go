package outfile

import (
	"context"
	"github.com/balibuild/winio/pkg/guid"
	"log"
	"sort"
	"strconv"
	"test_task/internal/app/config"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"

	"test_task/internal/app/database"
)

type PDFFile struct {
	outFilesDir      string
	apiKey           string
	font             string
	boldFont         string
	tableBorderWidth float64
	db               database.IDatabase
}

func New(cfg *config.Config, db database.IDatabase) (*PDFFile, error) {
	pdf := PDFFile{}

	pdf.outFilesDir = cfg.Parser.OutFilesDirectory
	pdf.apiKey = cfg.Parser.PdfApiKey
	pdf.font = cfg.Parser.Font
	pdf.boldFont = cfg.Parser.BoldFont
	pdf.tableBorderWidth = cfg.Parser.TableBorderWidth
	pdf.db = db

	err := license.SetMeteredKey(pdf.apiKey)
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
	c, err := f.createPdf(records)
	if err != nil {
		return err
	}

	// Write to output file.
	if err := c.WriteToFile(f.outFilesDir + "\\" + records[0].UnitGuid.String() + ".pdf"); err != nil {
		log.Fatal(err)
	}

	return nil
}

func getUniqueGUid(records []database.Record) []guid.GUID {

	var uniqGuids []guid.GUID
	var prevGuid guid.GUID
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

func (f *PDFFile) createPdf(records []database.Record) (*creator.Creator, error) {
	// Create report fonts.
	font, err := model.NewStandard14Font(model.StdFontName(f.font))
	if err != nil {
		return nil, err
	}

	fontBold, err := model.NewStandard14Font(model.StdFontName(f.boldFont))
	if err != nil {
		return nil, err
	}

	c := creator.New()
	pageSize := creator.PageSize{creator.PageSizeA4[1], creator.PageSizeA4[0]}
	c.SetPageSize(pageSize)

	table := c.NewTable(15)
	table.SetMargins(0, 0, 10, 0)

	// Draw table header.
	f.addCell(c, table, "n", fontBold)
	f.addCell(c, table, "mqtt", fontBold)
	f.addCell(c, table, "invid", fontBold)
	f.addCell(c, table, "unit_guid", fontBold)
	f.addCell(c, table, "msg_id", fontBold)
	f.addCell(c, table, "text", fontBold)
	f.addCell(c, table, "context", fontBold)
	f.addCell(c, table, "class", fontBold)
	f.addCell(c, table, "level", fontBold)
	f.addCell(c, table, "area", fontBold)
	f.addCell(c, table, "addr", fontBold)
	f.addCell(c, table, "block", fontBold)
	f.addCell(c, table, "type", fontBold)
	f.addCell(c, table, "bit", fontBold)
	f.addCell(c, table, "invert_bit", fontBold)

	for _, record := range records {
		f.addCell(c, table, strconv.Itoa(record.N), font)
		f.addCell(c, table, string(record.MQTT), font)
		f.addCell(c, table, record.InvId, font)
		f.addCell(c, table, record.UnitGuid.String(), font)
		f.addCell(c, table, record.MsgId, font)
		f.addCell(c, table, record.Text, font)
		f.addCell(c, table, string(record.Context), font)
		f.addCell(c, table, record.Class, font)
		f.addCell(c, table, strconv.Itoa(record.Level), font)
		f.addCell(c, table, record.Area, font)
		f.addCell(c, table, record.Addr, font)
		f.addCell(c, table, record.Block, font)
		f.addCell(c, table, record.Type, font)
		f.addCell(c, table, strconv.Itoa(record.Bit), font)
		f.addCell(c, table, strconv.Itoa(record.InvertBit), font)
	}

	err = c.Draw(table)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (f *PDFFile) addCell(c *creator.Creator, table *creator.Table, text string, font *model.PdfFont) *creator.TableCell {
	cell := table.NewCell()

	p := c.NewStyledParagraph()
	p.Append(text).Style.Font = font

	cell.SetContent(p)
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, f.tableBorderWidth)

	return cell
}
