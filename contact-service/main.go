// Contact service for tenkdegom.no. Bind only to loopback; Caddy handles public TLS.
package main

import (
 "crypto/tls"
 "errors"
 "fmt"
 "log"
 "net"
 "net/http"
 "net/mail"
 "net/smtp"
 "os"
 "strings"
 "sync"
 "time"
)

var (
 mu sync.Mutex
 attempts []time.Time
)

func env(k string) string { v:=os.Getenv(k);if v=="" {log.Fatalf("missing %s",k)};return v }
func limited() bool {
 mu.Lock();defer mu.Unlock()
 now:=time.Now();recent:=attempts[:0]
 for _,t:=range attempts {if now.Sub(t)<time.Minute {recent=append(recent,t)}}
 attempts=recent
 if len(attempts)>=10{return true}
 attempts=append(attempts,now);return false
}
func clean(s string) string {return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s,"\r"," "), "\n"," "))}
func sendMail(host,user,password,to,name,reply,message string) error {
 addr:=net.JoinHostPort(host,"587")
 conn,err:=net.DialTimeout("tcp",addr,10*time.Second);if err!=nil{return err}
 defer conn.Close()
 _=conn.SetDeadline(time.Now().Add(25*time.Second))
 client,err:=smtp.NewClient(conn,host);if err!=nil{return err}
 defer client.Close()
 if err=client.StartTLS(&tls.Config{ServerName:host,MinVersion:tls.VersionTLS12});err!=nil{return err}
 if ok,_:=client.Extension("AUTH");!ok{return errors.New("SMTP AUTH unavailable")}
 if err=client.Auth(smtp.PlainAuth("",user,password,host));err!=nil{return err}
 if err=client.Mail(user);err!=nil{return err}
 if err=client.Rcpt(to);err!=nil{return err}
 writer,err:=client.Data();if err!=nil{return err}
 subject:="Ny henvendelse fra Tenk deg om"
 body:=fmt.Sprintf("From: %s\r\nTo: %s\r\nReply-To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nNavn: %s\r\nSvaradresse: %s\r\n\r\n%s\r\n",user,to,reply,subject,name,reply,strings.ReplaceAll(strings.ReplaceAll(message,"\r\n","\n"),"\n","\r\n"))
 _,err=writer.Write([]byte(body));if err!=nil{return err}
 if err=writer.Close();err!=nil{return err}
 return client.Quit()
}
func main(){
 host:=env("SMTP_HOST");user:=env("SMTP_USER");password:=env("SMTP_TOKEN");to:=env("CONTACT_TO")
 if _,err:=mail.ParseAddress(user);err!=nil{log.Fatal("invalid SMTP_USER")}
 if _,err:=mail.ParseAddress(to);err!=nil{log.Fatal("invalid CONTACT_TO")}
 mux:=http.NewServeMux()
 mux.HandleFunc("/api/contact",func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("X-Content-Type-Options","nosniff")
  if r.Method!="POST"{w.Header().Set("Allow","POST");http.Error(w,"Method not allowed",405);return}
  if origin:=r.Header.Get("Origin");origin!=""&&origin!="https://tenkdegom.no"{http.Error(w,"Forbidden",403);return}
  if !strings.HasPrefix(r.Header.Get("Content-Type"),"application/x-www-form-urlencoded"){http.Error(w,"Unsupported media type",415);return}
  r.Body=http.MaxBytesReader(w,r.Body,12<<10)
  if err:=r.ParseForm();err!=nil{http.Error(w,"Invalid request",400);return}
  if r.PostForm.Get("website")!=""{w.WriteHeader(200);fmt.Fprintln(w,"Takk!");return}
  name:=clean(r.PostForm.Get("name"));reply:=clean(r.PostForm.Get("email"));message:=strings.TrimSpace(r.PostForm.Get("message"))
  if len(name)>120||len(reply)>254||len(message)<10||len(message)>5000||strings.ContainsAny(reply,"\r\n") {http.Error(w,"Invalid fields",400);return}
  parsed,err:=mail.ParseAddress(reply);if err!=nil||parsed.Address!=reply{http.Error(w,"Invalid email",400);return}
  if limited(){http.Error(w,"Too many requests",429);return}
  if err:=sendMail(host,user,password,to,name,reply,message);err!=nil{log.Printf("SMTP delivery failed: %T",err);http.Error(w,"Delivery failed",502);return}
  w.Header().Set("Content-Type","application/json; charset=utf-8");fmt.Fprintln(w,`{"ok":true}`)
 })
 server:=&http.Server{Addr:"127.0.0.1:8787",Handler:mux,ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:35*time.Second,IdleTimeout:30*time.Second}
 log.Println("contact service listening on loopback:8787");log.Fatal(server.ListenAndServe())
}
