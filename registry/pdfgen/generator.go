package pdfgen

import (
	"bytes"
	"grpc-mafia/util"
	"html/template"
	"path/filepath"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	zlog "github.com/rs/zerolog/log"
)

type PdfGen struct {
	pdf_template     *template.Template
	default_ava_path string
}

func (g *PdfGen) Render(request RenderRequest) ([]byte, error) {
	var body bytes.Buffer

	if len(request.User.AvatarFilename) == 0 {
		request.User.AvatarFilename = g.default_ava_path
	}

	zlog.Info().Msg("executing pdf template")

	// apply the parsed HTML template data and keep the result in a Buffer
	if err := g.pdf_template.Execute(&body, request); err != nil {
		return nil, err
	}

	zlog.Info().Msg("New pdf generator")

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	zlog.Info().Msg("New page reader")

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	zlog.Info().Msg("configuring")

	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)

	// add the page to your generator
	pdfg.AddPage(page)

	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	zlog.Info().Msg("creating pdf")

	// magic
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	zlog.Info().Msg("pdf created")

	return pdfg.Bytes(), nil
}

func NewPDFGen() *PdfGen {
	root, err := util.GetProjectRoot()
	if err != nil {
		panic(err) // can't work without template
	}

	ava_path, err := filepath.Abs(filepath.Join(root, "registry", "pdfgen", "default-avatar.png"))
	if err != nil {
		panic(err) // can't work without default ava
	}

	templ_path, err := filepath.Abs(filepath.Join(root, "registry", "pdfgen", "template.html"))
	if err != nil {
		panic(err) // can't work without template
	}

	pdf_template, err := template.ParseFiles(templ_path)
	if err != nil {
		panic(err) // can't work without template
	}

	return &PdfGen{
		pdf_template:     pdf_template,
		default_ava_path: ava_path,
	}
}
