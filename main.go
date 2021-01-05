package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var totalParts, availParts int

type Price struct {
	Min  int
	Cost float32
}

func parsePrice(s string) []Price {
	z := make([]Price, 0)
	x := strings.Split(s, ",")
	for _, a := range x {
		colon := strings.Index(a, ":")
		if colon == -1 {
			return z
		}
		qtyr := a[:colon]
		cost := a[colon+1:]
		dash := strings.Index(qtyr, "-")
		i, _ := strconv.ParseInt(qtyr[:dash], 10, 64)
		cst, _ := strconv.ParseFloat(cost, 64)
		c := Price{Min: int(i), Cost: float32(cst)}
		z = append(z, c)
	}
	return z
}

type Part struct {
	LCSC      string
	Cat1      string
	Cat2      string
	ManuPart  string
	Pkg       string
	Solder    int
	Manu      string
	BasicType bool
	Descr     string
	Datasheet string
	Price     []Price
	Stock     int
}

// ssconvert parts.xls parts.csv
var parts map[string]Part

func main() {
	parts = make(map[string]Part)
	log.Println("Opening parts library file.....")
	csvfile, err := os.Open("parts.csv")
	if err != nil {
		log.Panic("Error tying to open parts.csv :", err)
	}
	stat, _ := csvfile.Stat()
	log.Printf("Loading parts library %v.....", stat.ModTime().Format(time.RFC3339))
	r := csv.NewReader(csvfile)
	r.Read() // Header
	for {
		r, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		p := Part{}
		p.LCSC = strings.ToUpper(r[0])
		p.Cat1 = strings.ToUpper(r[1])
		p.Cat2 = strings.ToUpper(r[2])
		p.ManuPart = strings.ToUpper(r[3])
		p.Pkg = strings.ToUpper(r[4])
		i, _ := strconv.ParseInt(r[5], 10, 64)
		p.Solder = int(i)
		p.Manu = strings.ToUpper(r[6])
		p.BasicType = r[7] == "Basic"
		p.Descr = strings.ToUpper(r[8])
		p.Datasheet = strings.ToUpper(r[9])
		p.Price = parsePrice(r[10])
		i, _ = strconv.ParseInt(r[11], 10, 64)
		p.Stock = int(i)
		parts[p.LCSC] = p
		totalParts++
		if p.Stock > 0 {
			availParts++
		}
	}
	log.Println("Starting webserver")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d := struct {
			Totalparts     int
			Availparts     int
			PartsTimestamp time.Time
		}{
			Totalparts:     totalParts,
			Availparts:     availParts,
			PartsTimestamp: stat.ModTime(),
		}
		tmpl := template.Must(template.ParseFiles("www/index.html"))
		tmpl.Execute(w, d)
	})
	http.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	txt := r.FormValue("txt")
	found := find(strings.ToUpper(txt))
	instock := r.FormValue("instock") == "true"
	p := make([]Part, 0)
	for _, a := range found {
		if !instock || parts[a].Stock > 0 {
			p = append(p, parts[a])
		}
	}

	j, _ := json.Marshal(p)
	w.Write([]byte(j))
}

func find(s string) []string {
	r := make([]string, 0)
	for a, b := range parts {
		if strings.Index(a, s) >= 0 {
			r = append(r, a)
			continue
		}
		if strings.Index(b.Descr, s) >= 0 {
			r = append(r, a)
			continue
		}
		if strings.Index(b.Cat1, s) >= 0 {
			r = append(r, a)
			continue
		}
		if strings.Index(b.Cat2, s) >= 0 {
			r = append(r, a)
			continue
		}
		if strings.Index(b.ManuPart, s) >= 0 {
			r = append(r, a)
			continue
		}
		if strings.Index(b.Manu, s) >= 0 {
			r = append(r, a)
			continue
		}
	}
	return r
}
