# jlcpcbparts
A simple golang webserver that provides a much easier to search interface into the JLCPCB parts inventory.


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

