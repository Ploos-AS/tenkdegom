export type Question = { category: string; prompt: string; advice: string };
export const questions: Question[] = [
  {
    "category": "Nettleser",
    "prompt": "Oppdaterer du nettleseren og enheten regelmessig?",
    "advice": "Oppdater nettleseren og operativsystemet. Sikkerhetsoppdateringer tetter kjente svakheter."
  },
  {
    "category": "Nettleser",
    "prompt": "Har du gjennomgått personvern- og sporingsinnstillingene i nettleseren?",
    "advice": "Se gjennom sporingsbeskyttelse, tillatelser og informasjonskapsler i nettleseren."
  },
  {
    "category": "Passord",
    "prompt": "Bruker du unike passord til viktige kontoer?",
    "advice": "Bruk unike passord. En passordbehandler kan gjøre dette enklere."
  },
  {
    "category": "Passord",
    "prompt": "Har du aktivert tofaktorautentisering på e-post og andre viktige kontoer?",
    "advice": "Aktiver tofaktorautentisering, gjerne med autentiseringsapp eller sikkerhetsnøkkel."
  },
  {
    "category": "E-post",
    "prompt": "Kontrollerer du avsender og lenker før du åpner uventede meldinger?",
    "advice": "Vær skeptisk til uventede meldinger. Gå direkte til tjenestens nettsted ved tvil."
  },
  {
    "category": "E-post",
    "prompt": "Vet du hvordan du kan oppdage og håndtere mistenkelige e-poster?",
    "advice": "Lær å kjenne igjen svindelforsøk og bruk rapporteringsfunksjonen i e-posttjenesten."
  },
  {
    "category": "Apper",
    "prompt": "Har du gått gjennom hvilke apper som har tilgang til kamera, mikrofon og posisjon?",
    "advice": "Fjern tillatelser apper ikke trenger, særlig posisjon, mikrofon og kontakter."
  },
  {
    "category": "Apper",
    "prompt": "Sletter du apper du ikke lenger bruker?",
    "advice": "Avinstaller apper du ikke bruker, og kontroller tilknyttede kontoer."
  },
  {
    "category": "Skylagring",
    "prompt": "Vet du hvilke tjenester som lagrer bildene og dokumentene dine?",
    "advice": "Kartlegg hvor filer lagres og hvilke delingsinnstillinger som gjelder."
  },
  {
    "category": "Skylagring",
    "prompt": "Kontrollerer du hvem som har tilgang til delte mapper og lenker?",
    "advice": "Gjennomgå delingslenker og fjern tilgang som ikke lenger trengs."
  },
  {
    "category": "KI",
    "prompt": "Unngår du å dele sensitive personopplysninger med KI-tjenester uten å kjenne vilkårene?",
    "advice": "Ikke legg inn sensitive opplysninger før du forstår hvordan tjenesten behandler dem."
  },
  {
    "category": "KI",
    "prompt": "Sjekker du personverninnstillingene i KI-tjenestene du bruker?",
    "advice": "Undersøk lagring, treningsbruk og muligheten til å slette samtaler."
  }
];
