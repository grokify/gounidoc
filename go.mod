module github.com/grokify/gounidoc

go 1.26.0

require (
	github.com/grokify/mogo v0.74.6
	github.com/jessevdk/go-flags v1.6.1
	github.com/modelcontextprotocol/go-sdk v1.6.1
	github.com/plexusone/omniskill v0.8.0
	github.com/spf13/cobra v1.10.2
	github.com/unidoc/unioffice v1.39.0
	github.com/unidoc/unipdf/v3 v3.69.0
)

require (
	github.com/adrg/strutil v0.3.1 // indirect
	github.com/adrg/sysfont v0.1.2 // indirect
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/caarlos0/env/v11 v11.4.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/jsonschema-go v0.4.3 // indirect
	github.com/gorilla/i18n v0.0.0-20150820051429-8b358169da46 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/inconshreveable/log15 v3.0.0-testing.5+incompatible // indirect
	github.com/inconshreveable/log15/v3 v3.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/llgcode/draw2d v0.0.0-20231212091825-f55e0c776b44 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/segmentio/encoding v0.5.4 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/unidoc/emf v0.1.0 // indirect
	github.com/unidoc/freetype v0.2.3 // indirect
	github.com/unidoc/garabic v0.0.0-20220702200334-8c7cb25baa11 // indirect
	github.com/unidoc/pkcs7 v0.2.0 // indirect
	github.com/unidoc/timestamp v0.0.0-20200412005513-91597fd3793a // indirect
	github.com/unidoc/unichart v0.4.0 // indirect
	github.com/unidoc/unitype v0.5.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.ngrok.com/muxado/v2 v2.0.1 // indirect
	golang.ngrok.com/ngrok v1.13.0 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/exp v0.0.0-20260603202125-055de637280b // indirect
	golang.org/x/image v0.41.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/term v0.43.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Force log15/v3 to version compatible with ngrok v1.13.0 (has ext.RandId)
replace github.com/inconshreveable/log15/v3 => github.com/inconshreveable/log15/v3 v3.0.0-testing.5
