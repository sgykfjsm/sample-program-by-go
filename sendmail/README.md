# Mail Sample Program

## Testing

```
# Start mail server
$ ./MailHog_darwin_amd64 # https://github.com/mailhog/MailHog
```

```
# In another terminal window
$ go run ./main.go
```

Open your browser and then go to http://127.0.0.1:8025

## Reference

- https://github.com/go-gomail/gomail
- https://godoc.org/gopkg.in/gomail.v2
- https://github.com/golang/go/wiki/SendingMail
- http://tmichel.github.io/2014/10/12/golang-send-test-email/
