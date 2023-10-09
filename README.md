# go-mod-papertrail

This is a go module to print logs to papertrail

This is a private repo, so in order to use it in other go project, do this:
```
export GOPRIVATE=github.com/Billes/go-mod-papertrail
go get github.com/Billes/go-mod-papertrail@V0.1.5
```

## TLDR; 
Two environment variables has to be set up:

*PAPER_TRAIL_URL* : Url to papertrail including port
*LOG_LEVEL*: Level of logs we want to see. Its including the one you specify.

1. CRITICAL
2. ERROR
3. WARNING
4. INFO
5. DEBUG



In order to activate it do this:
```
	papertrailUrl := os.Getenv("PAPER_TRAIL_URL")
	if papertrailUrl == "" {
		log.Fatal("Need environment variable PAPER_TRAIL_URL")
	}
	papertrailSystem := os.Getenv("PAPER_TRAIL_SYSTEM")
	if papertrailSystem == "" {
		log.Fatal("Need environment variable PAPER_TRAIL_SYSTEM")
	}

	papertrail.Init(papertrailUrl, papertrailSystem)
```

During local developmen you may set the "PAPER_TRAIL_URL" to "", and it will only logg locally. 


And this is how you use it:

```
papertrail.Error([]string{"Subscribe", "GetMessages"}, "Got an error receiving messages", err.Error())

papertrail.Error|Debug|Info|Critical(<array-of-string>,<msg>,<error string>)

<array-of-string> : keywords to search for i papertrail
<msg>             : friform text to explain the error
<error string>    : error string from err  
```

