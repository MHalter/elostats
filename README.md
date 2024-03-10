# elostats

## Projekt erstellen

### abhängige Bibliotheken laden 
- go-mssqldb  
  ``go get github.com/denisenkom/go-mssqldb``


## Einbinden in telegraf
Hierzu die telegraf.conf anpassen

```` toml
[[inputs.exec]]
  commands = ["./elostats"]
  timeout = "5s"
  data_format = "influx"

  [inputs.exec.config]
    DBHost = "localhost"
    DBName = "database"
    DBUser = "user"
    DBPassword = "password"
````
Außerdem muss das kompilierte Plugin in das Telegraf-Verzeichnis abgelegt werden