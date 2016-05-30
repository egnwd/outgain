package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type attribute struct {
	class, fill string
}

type style struct {
	attributes []attribute
}

func SVGSpriteHandler(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		const size = 64
		vars := mux.Vars(r)
		colour := vars["colour"]

		s := generateStyle(colour)

		svg, err := os.Open(staticDir + "/images/creature-base.svg")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		body, err := ioutil.ReadAll(svg)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		defer svg.Close()

		// Find & Replace style
		const styleComment = "<!-- style -->"
		replace := strings.NewReplacer(styleComment, s.String())
		newSvg := replace.Replace(string(body))
		// ioutil.WriteFile("/tmp/creature", []byte(newSvg), 0644)
		fmt.Fprint(w, newSvg)
	})
}

func generateStyle(colour string) style {
	outline := attribute{class: "outline", fill: "FFFFFF"}

	outerEnv := attribute{class: "outer-envelope", fill: colour}
	innerEnv := attribute{class: "inner-envelope", fill: colour}

	body1 := attribute{class: "body-1", fill: colour}
	body2 := attribute{class: "body-2", fill: colour}
	body3 := attribute{class: "body-3", fill: colour}
	body4 := attribute{class: "body-4", fill: colour}
	body5 := attribute{class: "body-5", fill: colour}
	body6 := attribute{class: "body-6", fill: colour}

	return style{
		attributes: []attribute{
			outline,
			outerEnv, innerEnv,
			body1, body2, body3, body4, body5, body6,
		},
	}
}

func (a attribute) String() string {
	return fmt.Sprintf(".%s { fill:#%s; }", a.class, a.fill)
}

func (s style) String() string {
	var attributes string

	for _, a := range s.attributes {
		attributes += a.String()
	}

	return fmt.Sprintf(`<style type="text/css">%s</style>`, attributes)
}
