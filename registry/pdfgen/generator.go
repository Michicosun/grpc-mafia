package pdfgen

import (
	"bytes"
	"text/template"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func NewPDF() []byte {
	var templ *template.Template
	var err error

	// use Go's default HTML template generation tools to generate your HTML
	if templ, err = template.ParseFiles("/home/mdartemenko/hse/soa/grpc-mafia/registry/pdfgen/template.html"); err != nil {
		panic(err)
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err = templ.Execute(&body, nil); err != nil {
		panic(err)
	}

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		panic(err)
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)

	// add the page to your generator
	pdfg.AddPage(page)

	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	// pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)

	// magic
	err = pdfg.Create()
	if err != nil {
		panic(err)
	}

	return pdfg.Bytes()
}
