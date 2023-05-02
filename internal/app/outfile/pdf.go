package outfile

import (
	"context"
	"github.com/balibuild/winio/pkg/guid"
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
	font             model.StdFontName
	boldFont         model.StdFontName
	tableBorderWidth float64
	db               database.IDatabase
}

func New(cfg *config.Config, db database.IDatabase) (*PDFFile, error) {
	pdf := PDFFile{}

	pdf.outFilesDir = cfg.Parser.OutFilesDirectory
	pdf.apiKey = cfg.Parser.PdfApiKey
	pdf.font = model.StdFontName(cfg.Parser.Font)
	pdf.boldFont = model.StdFontName(cfg.Parser.BoldFont)
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

	if len(records) == 0 {
		return nil
	}

	// Write to output file.
	if err = c.WriteToFile(f.outFilesDir + "\\" + records[0].UnitGuid.String() + ".pdf"); err != nil {
		return err
	}

	return nil
}

func getUniqueGUid(records []database.Record) []guid.GUID {

	var uniqGuids []guid.GUID
	var prevGuid guid.GUID
	for _, record := range records {
		if record.UnitGuid == prevGuid {
			continue
		}

		uniqGuids = append(uniqGuids, record.UnitGuid)
		prevGuid = record.UnitGuid
	}

	return uniqGuids
}

func (f *PDFFile) createPdf(records []database.Record) (*creator.Creator, error) {
	c := creator.New()
	pageSize := creator.PageSize{creator.PageSizeA4[1], creator.PageSizeA4[0]}
	c.SetPageSize(pageSize)

	table := c.NewTable(15)
	table.SetMargins(0, 0, 10, 0)

	// Draw table header.
	err := f.addHeader(c, table)
	if err != nil {
		return nil, err
	}

	font, err := model.NewStandard14Font(f.font)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		f.addRow(c, table, font, record)
	}

	err = c.Draw(table)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (f *PDFFile) addHeader(c *creator.Creator, table *creator.Table) error {
	font, err := model.NewStandard14Font(f.boldFont)
	if err != nil {
		return err
	}

	addCell(c, table, "n", font, f.tableBorderWidth)
	addCell(c, table, "mqtt", font, f.tableBorderWidth)
	addCell(c, table, "invid", font, f.tableBorderWidth)
	addCell(c, table, "unit_guid", font, f.tableBorderWidth)
	addCell(c, table, "msg_id", font, f.tableBorderWidth)
	addCell(c, table, "text", font, f.tableBorderWidth)
	addCell(c, table, "context", font, f.tableBorderWidth)
	addCell(c, table, "class", font, f.tableBorderWidth)
	addCell(c, table, "level", font, f.tableBorderWidth)
	addCell(c, table, "area", font, f.tableBorderWidth)
	addCell(c, table, "addr", font, f.tableBorderWidth)
	addCell(c, table, "block", font, f.tableBorderWidth)
	addCell(c, table, "type", font, f.tableBorderWidth)
	addCell(c, table, "bit", font, f.tableBorderWidth)
	addCell(c, table, "invert_bit", font, f.tableBorderWidth)

	return nil
}

func (f *PDFFile) addRow(c *creator.Creator, table *creator.Table, font *model.PdfFont, record database.Record) {
	addCell(c, table, strconv.Itoa(record.N), font, f.tableBorderWidth)
	addCell(c, table, string(record.MQTT), font, f.tableBorderWidth)
	addCell(c, table, record.InvId, font, f.tableBorderWidth)
	addCell(c, table, record.UnitGuid.String(), font, f.tableBorderWidth)
	addCell(c, table, record.MsgId, font, f.tableBorderWidth)
	addCell(c, table, record.Text, font, f.tableBorderWidth)
	addCell(c, table, string(record.Context), font, f.tableBorderWidth)
	addCell(c, table, record.Class, font, f.tableBorderWidth)
	addCell(c, table, strconv.Itoa(record.Level), font, f.tableBorderWidth)
	addCell(c, table, record.Area, font, f.tableBorderWidth)
	addCell(c, table, record.Addr, font, f.tableBorderWidth)
	addCell(c, table, record.Block, font, f.tableBorderWidth)
	addCell(c, table, record.Type, font, f.tableBorderWidth)
	addCell(c, table, strconv.Itoa(record.Bit), font, f.tableBorderWidth)
	addCell(c, table, strconv.Itoa(record.InvertBit), font, f.tableBorderWidth)
}

func addCell(c *creator.Creator, table *creator.Table, text string, font *model.PdfFont, borderWidth float64) *creator.TableCell {
	cell := table.NewCell()

	p := c.NewStyledParagraph()
	p.Append(text).Style.Font = font

	cell.SetContent(p)
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, borderWidth)

	return cell
}
