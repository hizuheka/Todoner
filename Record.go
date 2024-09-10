package main

import "time"

// データを格納する構造体
type Record struct {
	IP              string
	TerminalName    string
	AgreementNumber string
	LifeEvent       string
	ApplicationNum  string
	ApplicationCnt  string
	MovePersonCnt   string
	AppFormName     string
	StartDateTime   string
	ProcessDateTime string
	UsageTime       time.Duration
}

func (r Record) GetKu() string {
	return r.AgreementNumber[0:1]
}

// CSV行を構造体に変換
func parseRecord(fields []string) (Record, error) {
	if ut, err := parseTime(fields[10]); err != nil {
		return Record{}, err
	} else {
		return Record{
			IP:              fields[0],
			TerminalName:    fields[1],
			AgreementNumber: fields[2],
			LifeEvent:       fields[3],
			ApplicationNum:  fields[4],
			ApplicationCnt:  fields[5],
			MovePersonCnt:   fields[6],
			AppFormName:     fields[7],
			StartDateTime:   fields[8],
			ProcessDateTime: fields[9],
			UsageTime:       ut,
		}, nil
	}
}

func parseTime(timeStr string) (time.Duration, error) {
	// 時間形式をパース
	parsedTime, err := time.Parse("3:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	// パースした時間をDurationに変換
	return time.Duration(parsedTime.Hour())*time.Hour + time.Duration(parsedTime.Minute())*time.Minute + time.Duration(parsedTime.Second())*time.Second, nil
}
