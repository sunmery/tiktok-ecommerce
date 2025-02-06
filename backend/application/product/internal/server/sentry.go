package server

import "github.com/getsentry/sentry-go"

func InitSentry(){
	sentry.Init(sentry.ClientOptions{
		Dsn: "https://6904a52570f3193e3ca484c1b01048ed@o4508714183819264.ingest.us.sentry.io/4508714186833920",
		AttachStacktrace: true, // recommended
	})
}
