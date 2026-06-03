package scraper

func sources() []Source {
	return []Source{
		{
			Brand: "Pertamina",
			URL: getenv(
				"PERTAMINA_PRICE_URL",
				"https://pertaminapatraniaga.com/api/api/v1/post/get-by-slug/page/harga-terbaru-bbm?language=id",
			),
		},
		{
			Brand: "Shell",
			URL: getenv(
				"SHELL_PRICE_URL",
				"https://www.shell.co.id/in_id/pengendara-bermotor/bahan-bakar-shell/harga-bahan-bakar-shell.html",
			),
		},
		{
			Brand: "BP",
			URL: getenv(
				"BP_PRICE_URL",
				"https://www.bp.com/id_id/indonesia/home/produk-dan-layanan/spbu/harga.html",
			),
		},
		{
			Brand: "Vivo",
			URL: getenv(
				"VIVO_STOCK_URL",
				"https://spbuvivo-ketersediaanstok.azurewebsites.net",
			),
		},
	}
}
