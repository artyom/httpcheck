`httpcheck` periodically issues HTTP `HEAD` request on URLs from file and logs
URL and response code to stderr.

	Usage of httpcheck:
	  -delay=30s: delay between each check cycle
	  -n=10: number of concurrent checks to make
	  -nokeepalive=false: disable keep alive
	  -urls="": file with urls, one url per line
