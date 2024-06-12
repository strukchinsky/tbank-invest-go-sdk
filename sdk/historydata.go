package investgo

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	pb "github.com/floatdrop/tbank-invest-go-sdk"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrNotFound        = fmt.Errorf("zip file not found")
	ErrTooManyRequests = fmt.Errorf("too many requests")
	ErrNotOk           = fmt.Errorf("response status is not ok")
)

// DownloadHistoryData downloads zip file with all 1-minute candles from history-data endpoint and
// returns it as slice of protobufs in investapi format.
func (c *Client) DownloadHistoryData(figi string, year int) ([]*pb.HistoricCandle, error) {
	url := fmt.Sprintf("https://invest-public-api.tinkoff.ru/history-data?figi=%s&year=%d", figi, year)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrTooManyRequests
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrNotOk
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return readHistoryZip(content)
}

// readHistoryZip parses downloaded zip file from history_data to slice of protobufs
func readHistoryZip(content []byte) ([]*pb.HistoricCandle, error) {
	archive, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip.NewReader: %w", err)
	}

	out := make(chan *pb.HistoricCandle, 4096)

	var wg sync.WaitGroup
	for _, file := range archive.File {
		f := file
		wg.Add(1)
		go func() {
			defer wg.Done()

			reader, err := f.Open()
			if err != nil {
				return
			}
			defer func() { _ = reader.Close() }()

			for _, row := range parseCsv(reader) {
				out <- row
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	result := make([]*pb.HistoricCandle, 0, 4096)
	for row := range out {
		result = append(result, row)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.GetSeconds() > result[j].Time.GetSeconds()
	})

	return result, nil
}

func parseCsv(reader io.Reader) []*pb.HistoricCandle {
	rows := make([]*pb.HistoricCandle, 0, 1024)

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			break
		}

		data := &pb.HistoricCandle{}

		t, err := time.Parse(time.RFC3339, line[1])
		if err != nil {
			log.Printf("Failed to parse time from '%s'", line[1])
			continue
		}
		data.Time = timestamppb.New(t)

		openPrice, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			log.Printf("Failed to parse open price from '%s'", line[2])
			continue
		}
		data.Open = FloatToQuotation(openPrice)

		closePrice, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			log.Printf("Failed to parse close price from '%s'", line[3])
			continue
		}
		data.Close = FloatToQuotation(closePrice)

		highPrice, err := strconv.ParseFloat(line[4], 64)
		if err != nil {
			log.Printf("Failed to parse high price from '%s'", line[4])
			continue
		}
		data.High = FloatToQuotation(highPrice)

		lowPrice, err := strconv.ParseFloat(line[5], 64)
		if err != nil {
			log.Printf("Failed to parse low price from '%s'", line[5])
			continue
		}
		data.Low = FloatToQuotation(lowPrice)

		volume, err := strconv.ParseInt(line[6], 10, 64)
		if err != nil {
			log.Printf("Failed to parse volume from '%s'", line[6])
			continue
		}
		data.Volume = volume

		rows = append(rows, data)
	}

	return rows
}
