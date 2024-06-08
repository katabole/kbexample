package actions

/*
import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	"github.com/golang/gddo/httputil"
	mailgun "github.com/mailgun/mailgun-go/v4"
	"github.com/unrolled/secure"
)

type EmailFn = func(context.Context, string, string, string) error

var SendEmail = func(ctx context.Context, subject string, body string, receiver string) error {
	mgDomain := envy.Get("MAILGUN_DOMAIN", "")
	mgKey := envy.Get("MAILGUN_API_KEY", "")
	if mgDomain == "" || mgKey == "" {
		return fmt.Errorf("missing either MAILGUN_DOMAIN or MAILGUN_API_KEY env-var")
	}
	mg := mailgun.NewMailgun(mgDomain, mgKey)
	message := mg.NewMessage("noreply@mg.gracepointonline.org", subject, body, receiver)
	_, _, err := mg.Send(ctx, message)
	return err
}

*/
