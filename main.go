package main

import (
	"github.com/rawnly/pdf-template-poc/pkg/qrcode"
	"github.com/signintech/gopdf"
)

func main() {

	pdf := gopdf.GoPdf{}
	pageSize := *gopdf.PageSizeA5
	pdf.Start(gopdf.Config{
		PageSize: pageSize,
	})
	pdf.AddPage()
	tpl1 := pdf.ImportPage("pdf_template_flat.pdf", 1, "/MediaBox")
	pdf.UseImportedTemplate(tpl1, 0, 0, pageSize.W, pageSize.H)

	_ = qrcode.DrawStyledQRCode(&pdf, qrcode.Options{
		X:        84,
		Y:        113,
		Size:     74,
		Content:  "https://staging.satispay.com/promo/10PERTE/qrcode/demo",
		IconSize: qrcode.IconSizeMedium,
		ShowLogo: true,
	})

	_ = qrcode.DrawStyledQRCode(&pdf, qrcode.Options{
		X:        279,
		Y:        113,
		Size:     74,
		Content:  "https://staging.satispay.com/promo/10PERTE/qrcode/demo",
		IconSize: qrcode.IconSizeMedium,
		ShowLogo: true,
	})

	_ = qrcode.DrawStyledQRCode(&pdf, qrcode.Options{
		X:        264,
		Y:        403,
		Size:     82,
		Content:  "https://staging.satispay.com/promo/10PERTE/qrcode/demo",
		IconSize: qrcode.IconSizeMedium,
		ShowLogo: true,
	})

	pdf.WritePdf("out.pdf")
}
