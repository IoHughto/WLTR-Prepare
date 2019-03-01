# WLTR-Prepare
Script to manage WLTR event preparation

## Usage
For now, all interaction with this software is via the command line and with various command-line flags.

The following flags are used:
*  -add string
   * New file to add
*  -byes string
   * CSV file with byes
*  -data string
   * File containing stored data (default "./local.json")
*  -drop string
   * Comma-separated list (no spaces!) of DCI numbers to drop
*  -new
   * Only consider new players
*  -obfuscate
   * Print obfuscated info instead of full info
*  -out string
   * Output file (default "./local.json")
*  -print
   * Print players in CSV format
*  -save
   * Save file

The typical workflow is: 
1) Add the my.cfbe file from an url like: /tools/sk/wltrme/\<event id\>, and save it. 
   * `go run prepare.go --add=cfbe.file --save`
2) Add the LPDB from WER of the players entered on-site, and save it. 
   * `go run prepare.go --add=wer.lpdb --save`
3) Import byes and save. 
   * `go run --byes=byes.file.xlsx --save`
4) If you need an obfuscated list for uploading to the website, do the following. 
   * `go run --print --obfuscate`
5) If you need to drop players, put their DCIs in a comma-separated list. 
   * `go run --drop=123456,87654321 --save`
6) Print out a csv to put into a file for importing. 
   * `go run --print`