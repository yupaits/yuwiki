module github.com/yupaits/yuwiki

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dchest/captcha v0.0.0-20170622155422-6a29415a8364
	github.com/denisenkom/go-mssqldb v0.0.0-20190315220205-a8ed825ac853 // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/gin-contrib/sessions v0.0.0-20190226023029-1532893d996f
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/jinzhu/now v1.0.0 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/lestrrat-go/envload v0.0.0-20180220234015-a3eb8ddeffcc // indirect
	github.com/lestrrat-go/file-rotatelogs v2.2.0+incompatible
	github.com/lestrrat-go/strftime v0.0.0-20180821113735-8b31f9c59b0f // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/matoous/go-nanoid v0.0.0-20181114085210-eab626deece6
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robfig/cron v1.1.0
	github.com/sirupsen/logrus v1.4.1
	github.com/stretchr/testify v1.4.0
	github.com/tebeka/strftime v0.0.0-20140926081919-3f9c7761e312 // indirect
	github.com/ugorji/go/codec v0.0.0-20190320090025-2dc34c0b8780 // indirect
	golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529
	golang.org/x/tools v0.0.0-20191206204035-259af5ff87bd // indirect
)

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190219172222-a4c6cb3142f2
	golang.org/x/net => github.com/golang/net v0.0.0-20190213061140-3a22650c66bd
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190222072716-a9d3bda3a223
)
