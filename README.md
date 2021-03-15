A basic package used for scraping emails from a website where parameters are defined in the config file


## Configuration

The configuration is set using the config.toml file, e.g.

```
-URL            "https://bsuir.by/en"       :website to visit       
-MaxDepth           1                       :maximum depth of tree to reach         
-OutputFile       bsuir.xml                 :xml file for writing emails to                 
```

## How to run

Working from root folder "github.com", ensure dependencies are installed with -
 ```
go mod vendor
 ```

Build binaries with 
 ```
go build -o app cmd/main.go 
 ```

Run .exe file
 ```
./app
 ```
