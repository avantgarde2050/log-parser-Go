# LogParser

LogParser is a command-line tool for parsing, filtering, and saving logs. It allows you to extract log data based on log level, time range, and format it in either plain text or JSON.

## Features
- Parse logs from a file
- Filter logs by level (INFO, ERROR, WARN, DEBUG)
- Filter logs by time range
- Save logs in text or JSON format

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/logparser.git
   cd logparser
   ```
2. Build the project:
   ```sh
   go build -o logparser
   ```

## Usage
### Parse and Save Logs
Save all logs in a file:
```sh
logparser save --file output.txt --format text
logparser save --file output.json --format json
```

### Filter Logs
Filter logs by level:
```sh
logparser filter --level ERROR
```
Filter logs by time range:
```sh
logparser filter --since "2025-02-26 12:00:00" --until "2025-02-27 12:00:00"
```

## License
MIT License

## Contributing
Pull requests are welcome!

