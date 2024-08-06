package spreadsheets

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	scopes = []string{
		drive.DriveFileScope,
		sheets.SpreadsheetsScope,
	}

	headers = []any{
		"Name", "Surname", "IDNum", "Birthday", "Citizen",
		"Visited", "Purpose", "Country", "City", "Address",
		"Postal", "Phone", "Email", "Time", "Image",
	}

	sheetRange        = "A%d:O%d"
	sheetHeadersRange = fmt.Sprintf(sheetRange, 1, 1)
)

const mimeTypeSpreadsheet = "application/vnd.google-apps.spreadsheet"

type Spreadsheets struct {
	client *http.Client

	period            time.Duration
	minAbnormalPeriod time.Duration
	maxAbnormalPeriod time.Duration

	folder       string
	sheetID      string
	lastSubmited time.Time

	sheetsSrv *sheets.Service
	driveSrv  *drive.Service
}

type Spreadsheeter interface {
	Spreadsheets() *Spreadsheets
}

func NewSpreadsheeter(getter kv.Getter) Spreadsheeter {
	return &spreadsheeter{
		getter: getter,
	}
}

type spreadsheeter struct {
	once   comfig.Once
	getter kv.Getter
}

func (c *spreadsheeter) Spreadsheets() *Spreadsheets {
	return c.once.Do(func() interface{} {
		var cfg struct {
			Credentials       string        `fig:"credentials,required"`
			Folder            string        `fig:"folder,required"`
			Period            time.Duration `fig:"period,required"`
			MinAbnormalPeriod time.Duration `fig:"min_abnormal_period,required"`
			MaxAbnormalPeriod time.Duration `fig:"max_abnormal_period,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "spreadsheets")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out spreadsheets config: %w", err))
		}

		creds, err := os.ReadFile(cfg.Credentials)
		if err != nil {
			panic(fmt.Errorf("unable to read client secret file: %w", err))
		}

		config, err := google.JWTConfigFromJSON(creds, scopes...)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}

		client := config.Client(context.Background())

		sheetsSrv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
		if err != nil {
			panic(fmt.Errorf("failed to create sheets service: %w", err))
		}

		driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
		if err != nil {
			panic(fmt.Errorf("failed to create drive service: %w", err))
		}

		return &Spreadsheets{
			client:            client,
			period:            cfg.Period,
			minAbnormalPeriod: cfg.MinAbnormalPeriod,
			maxAbnormalPeriod: cfg.MaxAbnormalPeriod,

			folder: cfg.Folder,

			sheetsSrv: sheetsSrv,
			driveSrv:  driveSrv,
		}
	}).(*Spreadsheets)
}

func (s *Spreadsheets) CreateTable() error {
	sheet, err := s.driveSrv.Files.Create(&drive.File{
		Name:     time.Now().UTC().Format("01/02/2006 15:04"),
		MimeType: mimeTypeSpreadsheet,
		Parents:  []string{s.folder},
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to create spreadsheet: %w", err)
	}

	_, err = s.sheetsSrv.Spreadsheets.Values.Update(sheet.Id, sheetHeadersRange, &sheets.ValueRange{
		Values: [][]any{headers},
	}).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("failed to set table headers: %w", err)
	}

	s.sheetID = sheet.Id

	return nil
}

func (s *Spreadsheets) FillTable(data [][]any) error {
	dataRange := len(data) + 1
	_, err := s.sheetsSrv.Spreadsheets.Values.Update(s.sheetID, fmt.Sprintf(sheetRange, 2, dataRange), &sheets.ValueRange{
		Values: data,
	}).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("failed to insert user data in table: %w", err)
	}

	return nil
}
