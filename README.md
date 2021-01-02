# jlcpcbparts
A simple golang webserver that provides a much easier to search interface into the JLCPCB parts inventory.  Currently provides a simple search with the option of showing only those with stock.


Before running program, go download the latest parts XLS file at https://jlcpcb.com/parts.
Once downloaded, run:
```
$> ssconvert NAME_OF_DOWNLOADED_FILE parts.csv
```
or you can load the XLS and export it as a CSV file named parts.csv and your favorite spreadsheet program.
then run the program:
```
$> go run main.go
```

Once running, open your browser to http://localhost:8888

![alt text](/screenshot.png?raw=true)
