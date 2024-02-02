# checkbiz-exporter
A simple web server written in [Go](https://go.dev/) that exposes [Checkbiz](https://checkbiz.se/) API call statistics as [Prometheus](https://prometheus.io/) metrics.

## Running
Build the docker image:
```sh
docker build -t checkbiz-exporter .
```
Run the image
```sh
docker run -it --rm -e CHECKBIZ_TOKEN=AAAABBBBBCCCDDD checkbiz-exporter:latest
```
alternatively, put your token in a file
```sh
docker run -it --rm -v /your/token/file:/token  -e CHECKBIZ_TOKEN_FILE=/token checkbiz-exporter:latest
```

## Example output

```
# HELP checkbiz_call_count The number of API calls made
# TYPE checkbiz_call_count counter
checkbiz_call_count{package="personsok",period="LastYear",product="PersonSearch"} 538
checkbiz_call_count{product="PersonSearch",package="personsok",period="ThisYear"} 0
checkbiz_call_count{package="personsok",period="LastMonth",product="PersonSearch"} 0
checkbiz_call_count{product="PersonSearch",package="personsok",period="ThisMonth"} 0
checkbiz_call_count{period="LastYear",product="CompanyInformation",package="CompanyInformation"} 30
checkbiz_call_count{period="ThisYear",product="CompanyInformation",package="CompanyInformation"} 10
checkbiz_call_count{period="LastMonth",product="CompanyInformation",package="CompanyInformation"} 9
checkbiz_call_count{product="CompanyInformation",package="CompanyInformation",period="ThisMonth"} 1
checkbiz_call_count{package="CompanyAutocomplete",period="ThisMonth",product="CompanyAutocomplete"} 7
checkbiz_call_count{product="CompanyAutocomplete",package="CompanyAutocomplete",period="LastYear"} 90
checkbiz_call_count{product="CompanyAutocomplete",package="CompanyAutocomplete",period="ThisYear"} 39
checkbiz_call_count{product="CompanyAutocomplete",package="CompanyAutocomplete",period="LastMonth"} 32

# HELP checkbiz_call_count The total number of API calls made for each period
# TYPE checkbiz_call_count counter

checkbiz_call_count_total{period="LastYear"} 658
checkbiz_call_count_total{period="ThisYear"} 49
checkbiz_call_count_total{period="LastMonth"} 41
checkbiz_call_count_total{period="ThisMonth"} 8
```
## Caveats
This is actually my first Go program which I converted from a Python reference implementation. PRs welcome!
