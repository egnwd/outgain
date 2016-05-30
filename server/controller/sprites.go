package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
		fmt.Fprint(w, newSvg)
	})
}

func generateStyle(colour string) style {
	outline := attribute{class: "outline", fill: "FFFFFF"}

	outerEnv := attribute{class: "outer-envelope", fill: shiftColour(colour, 0.4)}
	innerEnv := attribute{class: "inner-envelope", fill: shiftColour(colour, -0.2)}

	body1 := attribute{class: "body-1", fill: shiftColour(colour, 0)}
	body2 := attribute{class: "body-2", fill: shiftColour(colour, 0.1)}
	body3 := attribute{class: "body-3", fill: shiftColour(colour, 0.2)}
	body4 := attribute{class: "body-4", fill: shiftColour(colour, 0.3)}
	body5 := attribute{class: "body-5", fill: shiftColour(colour, 0.4)}

	return style{
		attributes: []attribute{
			outline,
			outerEnv, innerEnv,
			body1, body2, body3, body4, body5,
		},
	}
}

func shiftColour(colour string, correctionFactor float32) string {
	c, _ := strconv.ParseInt(colour, 16, 64)

	const rShift, gShift = 16, 8
	const max = 0xFF

	red := float32(c >> 16)
	green := float32((c >> 8) & max)
	blue := float32(c & max)

	if correctionFactor < 0 {
		correctionFactor++
		red *= correctionFactor
		green *= correctionFactor
		blue *= correctionFactor
	} else {
		red = (max-red)*correctionFactor + red
		green = (max-green)*correctionFactor + green
		blue = (max-blue)*correctionFactor + blue
	}

	value := ((int(red) & max) << 16) + ((int(green) & max) << 8) + (int(blue) & max)

	return fmt.Sprintf("%X", value)
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
