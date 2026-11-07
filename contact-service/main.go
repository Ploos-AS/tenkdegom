// Contact service for tenkdegom.no. Bind only to loopback; Caddy handles public TLS.
package main

import (
 "crypto/tls"
 "crypto/hmac"
 "crypto/rand"
 "crypto/sha256"
 "encoding/base64"
 "encoding/hex"
 "encoding/json"
 "strconv"
 "errors"
 "fmt"
 "log"
 "net"
 "net/http"
 "net/mail"
 "net/smtp"
 "net/textproto"
 "os"
 "strings"
 "sync"
 "time"
)

var (
 mu sync.Mutex
 attempts=map[string][]time.Time{}
 used=map[string]time.Time{}
 challengeKey=[32]byte{}
)
type challenge struct {
 Algorithm string `json:"algorithm"`
 Challenge string `json:"challenge"`
 Salt string `json:"salt"`
 Signature string `json:"signature"`
 MaxNumber int `json:"maxnumber"`
}
type solution struct {
 Algorithm string `json:"algorithm"`
 Challenge string `json:"challenge"`
 Salt string `json:"salt"`
 Signature string `json:"signature"`
 Number int `json:"number"`
}
const maxNumber=30000
func env(k string) string { v:=os.Getenv(k);if v=="" {log.Fatalf("missing %s",k)};return v }
func hexSHA(s string) string { sum:=sha256.Sum256([]byte(s));return hex.EncodeToString(sum[:]) }
func signature(s string) string { h:=hmac.New(sha256.New,challengeKey[:]);h.Write([]byte(s));return hex.EncodeToString(h.Sum(nil)) }
func newChallenge() (challenge,error) {
 var saltBytes [16]byte
 if _,err:=rand.Read(saltBytes[:]);err!=nil{return challenge{},err}
 var numBytes [4]byte
 if _,err:=rand.Read(numBytes[:]);err!=nil{return challenge{},err}
 number:=int(uint32(numBytes[0])<<24|uint32(numBytes[1])<<16|uint32(numBytes[2])<<8|uint32(numBytes[3]))%maxNumber
 salt:=hex.EncodeToString(saltBytes[:])+"?expires="+strconv.FormatInt(time.Now().Add(5*time.Minute).Unix(),10)
 hash:=hexSHA(salt+strconv.Itoa(number))
 return challenge{"SHA-256",hash,salt,signature(hash),maxNumber},nil
}
func verifySolution(encoded string) bool {
 raw,err:=base64.StdEncoding.DecodeString(encoded);if err!=nil||len(raw)>2048{return false}
 var s solution
 if json.Unmarshal(raw,&s)!=nil||s.Algorithm!="SHA-256"||s.Number<0||s.Number>=maxNumber{return false}
 if len(s.Challenge)!=64||len(s.Signature)!=64||!hmac.Equal([]byte(signature(s.Challenge)),[]byte(s.Signature)){return false}
 parts:=strings.Split(s.Salt,"?expires=")
 if len(parts)!=2||len(parts[0])!=32{return false}
 expires,err:=strconv.ParseInt(parts[1],10,64);if err!=nil||time.Now().Unix()>expires||expires>time.Now().Add(6*time.Minute).Unix(){return false}
 if hexSHA(s.Salt+strconv.Itoa(s.Number))!=s.Challenge{return false}
 mu.Lock();defer mu.Unlock()
 now:=time.Now()
 for key,until:=range used {if now.After(until){delete(used,key)}}
 if _,exists:=used[s.Salt];exists{return false}
 used[s.Salt]=time.Unix(expires,0)
 return true
}
func limited(ip string) bool {
 mu.Lock();defer mu.Unlock()
 now:=time.Now()
 total:=0
 for key,values:=range attempts {
  recent:=values[:0]
  for _,t:=range values {if now.Sub(t)<time.Minute{recent=append(recent,t)}}
  if len(recent)==0{delete(attempts,key)}else{attempts[key]=recent;total+=len(recent)}
 }
 if len(attempts[ip])>=5||total>=30{return true}
 attempts[ip]=append(attempts[ip],now)
 return false
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
func smtpFailureDetail(err error) string {
 var smtpErr *textproto.Error
 if errors.As(err,&smtpErr) {
  return fmt.Sprintf("SMTP response code %d",smtpErr.Code)
 }
 return fmt.Sprintf("error type %T",err)
}
func main(){
 host:=env("SMTP_HOST");user:=env("SMTP_USER");password:=env("SMTP_TOKEN");to:=env("CONTACT_TO")
 if _,err:=rand.Read(challengeKey[:]);err!=nil{log.Fatal("cannot initialize challenge key")}
 if _,err:=mail.ParseAddress(user);err!=nil{log.Fatal("invalid SMTP_USER")}
 if _,err:=mail.ParseAddress(to);err!=nil{log.Fatal("invalid CONTACT_TO")}
 mux:=http.NewServeMux()
 mux.HandleFunc("/api/contact",func(w http.ResponseWriter,r *http.Request){
  if r.Method=="GET"&&r.URL.Query().Get("challenge")=="1" {
   w.Header().Set("Cache-Control","no-store")
   w.Header().Set("Content-Type","application/json")
   c,err:=newChallenge();if err!=nil{http.Error(w,"Unavailable",503);return}
   _=json.NewEncoder(w).Encode(c);return
  }
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
  ip,_,_:=net.SplitHostPort(r.RemoteAddr)
  // Trust Caddy's forwarded client IP only for connections from loopback.
  if parsed:=net.ParseIP(ip);parsed!=nil&&parsed.IsLoopback(){
   forwarded:=strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"),",")[0])
   if candidate:=net.ParseIP(forwarded);candidate!=nil{ip=candidate.String()}
  }
  if limited(ip){http.Error(w,"Too many requests",429);return}
  if !verifySolution(r.PostForm.Get("altcha")){http.Error(w,"Bot verification failed",403);return}
  if err:=sendMail(host,user,password,to,name,reply,message);err!=nil{log.Printf("SMTP delivery failed: %s",smtpFailureDetail(err));http.Error(w,"Delivery failed",502);return}
  w.Header().Set("Content-Type","application/json; charset=utf-8");fmt.Fprintln(w,`{"ok":true}`)
 })
 server:=&http.Server{Addr:"127.0.0.1:8787",Handler:mux,ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:35*time.Second,IdleTimeout:30*time.Second}
 log.Println("contact service listening on loopback:8787");log.Fatal(server.ListenAndServe())
}
