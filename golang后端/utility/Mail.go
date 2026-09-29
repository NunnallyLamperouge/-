package utility

import (
	"example/config"

	"gopkg.in/gomail.v2"
)

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Config.Mail.SenderName)
	m.SetHeader("From", config.Config.Mail.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.Config.Mail.SmtpHost, config.Config.Mail.SmtpPort, config.Config.Mail.Username, config.Config.Mail.Password)

	return d.DialAndSend(m)
}

func SendMailFormatted(to, subject, content string) error {
	first := `
<!DOCTYPE html>
<html dir="ltr" lang="en">
<head>
  <meta name="viewport" content="width=device-width"/>
  <meta content="text/html; charset=UTF-8" http-equiv="Content-Type"/>
  <meta name="color-scheme" content="light"/>
  <meta name="supported-color-schemes" content="light"/>
  <style>
    @font-face {
      font-family: 'Inter';
      font-style: normal;
      font-weight: 400;
      mso-font-alt: 'sans-serif';
      src: url(https://rsms.me/inter/font-files/Inter-Regular.woff2?v=3.19) format('woff2');
    }
    * {
      font-family: 'Inter', sans-serif;
    }
  </style>
</head>
<body style="margin:0">
  <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="max-width:600px;min-width:300px;width:100%;margin-left:auto;margin-right:auto;padding:0.5rem">
    <tbody>
      <tr style="width:100%">
        <td>
          <h2 style="margin:0 0 12px 0;text-align:left;color:#111827;font-size:30px;line-height:36px;font-weight:700">
            <strong>研途智伴</strong>
          </h2>
          <p style="font-size:15px;line-height:26.25px;margin:0 0 20px 0;color:#374151">`
	second := `
          </p>
          <table align="center" width="100%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="max-width:100%;text-align:left;margin-bottom:0px">
            <tbody>
              <tr style="width:100%">
                <td>
                  <a href="` + config.Config.Gin.Domain + `" style="line-height:100%;text-decoration:none;display:inline-block;max-width:100%;color:#ffffff;background-color:#000000;border:2px solid #000000;font-size:14px;font-weight:500;border-radius:9999px;padding:12px 32px" target="_blank">
                    <span style="max-width:100%;display:inline-block;line-height:120%;mso-padding-alt:0px;mso-text-raise:9px">
                      前往官网 ->
                    </span>
                  </a>
                </td>
              </tr>
            </tbody>
          </table>
          <div style="height:64px"></div>
          <p style="font-size:15px;line-height:26.25px;margin:0 0 20px 0;color:#374151">
            感谢您使用我们强大的由AI驱动的工程化论文编写工具。
          </p>
          <p style="font-size:15px;line-height:26.25px;margin:0 0 20px 0;color:#374151">
            研途智伴团队
          </p>
        </td>
      </tr>
    </tbody>
  </table>
</body>
</html>
`
	totalConetnt := first + content + second
	return SendEmail(to, subject, totalConetnt)
}
