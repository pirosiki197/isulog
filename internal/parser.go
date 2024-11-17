package internal

import (
	"bufio"
	"cmp"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"github.com/samber/lo"
)

type Parser struct {
	r        io.Reader
	profiles []Profile
}

type Profile struct {
	Count  int
	Method string
	Path   string
	Min    time.Duration
	Max    time.Duration
	Sum    time.Duration
	Avg    time.Duration
	P90    time.Duration
	P95    time.Duration
	P99    time.Duration
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		r: r,
	}
}

func (p *Parser) Parse() error {
	data := make(map[string][]Record)
	dec := gob.NewDecoder(p.r)
	var record Record
	for {
		err := dec.Decode(&record)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			break
		}

		key := genKey(record)
		data[key] = append(data[key], record)
	}

	p.profiles = make([]Profile, 0, len(data))
	for _, records := range data {
		durations := lo.Map(records, func(r Record, _ int) time.Duration { return r.ResponseTime.Truncate(100 * time.Microsecond) })
		slices.Sort(durations)

		var (
			minTime time.Duration = time.Hour
			maxTime time.Duration
			sumTime time.Duration
		)
		for _, d := range durations {
			minTime = min(minTime, d)
			maxTime = max(maxTime, d)
			sumTime += d
		}

		profile := Profile{
			Count:  len(records),
			Method: records[0].Method,
			Path:   records[0].Path,
			Min:    minTime,
			Max:    maxTime,
			Sum:    sumTime,
			Avg:    (sumTime / time.Duration(len(durations))),
			P90:    durations[len(durations)*9/10],
			P95:    durations[len(durations)*95/100],
			P99:    durations[len(durations)*99/100],
		}
		p.profiles = append(p.profiles, profile)
	}

	return nil
}

var keys = []string{"count", "method", "path", "min", "max", "sum", "avg", "p90", "p95", "p99"}

func (p *Parser) Print(sort string) {
	// TODO: sort
	slices.SortFunc(p.profiles, func(a, b Profile) int { return -1 * cmp.Compare(a.Sum, b.Sum) })

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	print := func(len int, arg any) {
		format := fmt.Sprintf(" %%-%dv |", len) // ` %-${len}v |`
		fmt.Fprintf(w, format, arg)
	}
	strLength := calculateMaxStrLength(p.profiles)

	w.WriteByte('|')
	for _, k := range keys {
		print(strLength[k], k)
	}
	w.WriteByte('\n')

	// fmt.Fprintln(w, "| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |")

	for _, profile := range p.profiles {
		w.WriteByte('|')
		print(strLength["count"], profile.Count)
		print(strLength["method"], profile.Method)
		print(strLength["path"], profile.Path)
		print(strLength["min"], profile.Min)
		print(strLength["max"], profile.Max)
		print(strLength["sum"], profile.Sum)
		print(strLength["avg"], profile.Avg)
		print(strLength["p90"], profile.P90)
		print(strLength["p95"], profile.P95)
		print(strLength["p99"], profile.P99)
		w.WriteByte('\n')
	}
}

func calculateMaxStrLength(profiles []Profile) map[string]int {
	strLength := make(map[string]int)
	chMax := func(key string, value int) {
		strLength[key] = max(strLength[key], value)
	}
	for _, p := range profiles {
		chMax("count", len(fmt.Sprintf("%d", p.Count)))
		chMax("method", len(p.Method))
		chMax("path", len(p.Path))
		chMax("min", len(p.Min.String()))
		chMax("max", len(p.Max.String()))
		chMax("sum", len(p.Sum.String()))
		chMax("avg", len(p.Avg.String()))
		chMax("p90", len(p.P90.String()))
		chMax("p95", len(p.P95.String()))
		chMax("p99", len(p.P99.String()))
	}
	for _, k := range keys {
		chMax(k, len(k))
	}
	return strLength
}

func genKey(record Record) string {
	return fmt.Sprintf("%s|%s", record.Path, record.Method)
}
