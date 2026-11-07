const form=document.querySelector<HTMLFormElement>('#contact')!;
const feedback=document.querySelector<HTMLElement>('#feedback')!;
const send=document.querySelector<HTMLButtonElement>('#send')!;
const botstatus=document.querySelector<HTMLElement>('#botstatus')!;
const email=document.querySelector<HTMLInputElement>('#email')!;
const message=document.querySelector<HTMLTextAreaElement>('#message')!;
const validationSummary=document.querySelector<HTMLElement>('#validation-summary')!;
function setFieldError(field: HTMLInputElement | HTMLTextAreaElement, errorId: string, error: string) {
 const target=document.getElementById(errorId)!;
 target.textContent=error;
 target.hidden=!error;
 if(error) field.setAttribute('aria-invalid','true');
 else field.removeAttribute('aria-invalid');
}
function validateContact(): boolean {
 const emailValue=email.value.trim();
 const messageValue=message.value.trim();
 let emailError='';
 if(!emailValue) emailError='Oppgi e-postadressen din.';
 else if(!email.validity.valid) emailError='Oppgi en gyldig e-postadresse.';
 const messageError=!messageValue?'Skriv en melding.':messageValue.length<10?'Skriv en melding på minst 10 tegn.':'';
 setFieldError(email,'email-error',emailError);
 setFieldError(message,'message-error',messageError);
 if(emailError || messageError){
  validationSummary.textContent='Kontroller feltene som er markert nedenfor.';
  validationSummary.hidden=false;
  (emailError?email:message).focus();
  return false;
 }
 validationSummary.textContent='';
 validationSummary.hidden=true;
 return true;
}
for(const field of [email,message]) field.addEventListener('input',()=>{
 if(field.getAttribute('aria-invalid')==='true') validateContact();
});
async function proofOfWork() {
 const response=await fetch('/api/contact?challenge=1',{cache:'no-store'});
 if(!response.ok)throw new Error('challenge unavailable');
 const c=await response.json();
 if(c.algorithm!=='SHA-256'||!Number.isInteger(c.maxnumber)||c.maxnumber>100000)throw new Error('invalid challenge');
 const encoder=new TextEncoder();
 for(let n=0;n<c.maxnumber;n++){
  const hash=await crypto.subtle.digest('SHA-256',encoder.encode(c.salt+String(n)));
  const hex=Array.from(new Uint8Array(hash),b=>b.toString(16).padStart(2,'0')).join('');
  if(hex===c.challenge){
   const payload={algorithm:c.algorithm,challenge:c.challenge,salt:c.salt,signature:c.signature,number:n};
   return btoa(JSON.stringify(payload));
  }
  if(n%250===0)await new Promise(resolve=>setTimeout(resolve,0));
 }
 throw new Error('challenge unsolved');
}
form.addEventListener('submit',async e=>{
 e.preventDefault();
 if(!validateContact()) return;
 feedback.textContent='Sender …';send.disabled=true;
 try{
  botstatus.textContent='Kontrollerer at innsendingen er ekte …';
  const token=await proofOfWork();
  const fields=new URLSearchParams(new FormData(form) as any);
  fields.set('altcha',token);
  botstatus.textContent='';
  const response=await fetch('/api/contact',{method:'POST',body:fields,headers:{'Accept':'application/json','Content-Type':'application/x-www-form-urlencoded'},credentials:'same-origin'});
  if(!response.ok)throw new Error('Kunne ikke sende meldingen. Prøv igjen senere.');
  form.reset();feedback.textContent='Takk! Meldingen er sendt.';
 }catch{feedback.textContent='Meldingen ble ikke sendt. Prøv igjen senere.';}
 finally{send.disabled=false;botstatus.textContent='';}
});
