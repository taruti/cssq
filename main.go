package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var html = flag.Bool("h", false, "HTML")

func main() {
	flag.Parse()

	d, e := goquery.NewDocumentFromReader(os.Stdin)
	if e != nil {
		log.Fatal("parse error:", e)
	}
	var res [][]string
	for _, f := range flag.Args() {
		var imgSrc = false
		if strings.HasPrefix(f, ":src:") {
			imgSrc = true
			f = f[5:]
		}
		var r []string
		d.Find(f).Each(func(i int, s *goquery.Selection) {
			var dst string
			if *html {
				dst, e = s.Html()
				if e != nil {
					log.Println(e.Error())
				}
			} else if imgSrc {
				n := s.Get(s.Index())
				if n != nil {
					for _, a := range n.Attr {
						if a.Key == "src" {
							dst = a.Val
							break
						}
					}
				}
			} else {
				dst = s.Text()
			}
			dst = strings.TrimSpace(dst)
			r = append(r, dst)
		})
		res = append(res, r)
	}
	cw := csv.NewWriter(os.Stdout)
	any := true
	//	log.Println(res)
	for i := 0; any; i++ {
		var col []string
		any = false
		for _, ss := range res {
			var tmp string
			if len(ss) > i {
				tmp = ss[i]
				any = true
			}
			col = append(col, tmp)
			//			log.Println(i, len(ss), ss, "=>", col)
		}
		cw.Write(col)
	}
	cw.Flush()
}
