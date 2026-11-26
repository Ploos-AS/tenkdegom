# Tenk deg om

**Personvern angår oss alle.**

[Tenkdegom.no](https://tenkdegom.no) er et norsk folkeopplysningsprosjekt om digitalt personvern, eid og utgitt av **Ploos AS**.

## Formål

Digitalt personvern skal være forståelig, relevant og praktisk for alle, uavhengig av alder og teknisk bakgrunn. Grunnleggende kunnskap og verktøy skal være gratis.

## Status – 26. november 2026

- Norsk forside, læringsportal med 19 artikler, kildelenker og relaterte artikler.
- Personvernsjekken: 12 spørsmål, lokal behandling av svar, handlingsplan og utskrift/PDF via nettleseren.
- Ressursside, kontaktskjema, personvernerklæring og redaksjonelle prinsipper.
- Sitemap, robots.txt, canonical, Open Graph og strukturerte artikkeldata.
- Grunnleggende tilgjengelighetstiltak; full WCAG 2.2 AA-revisjon gjenstår.
- CI med Astro-kontroll, regresjonstester og statisk bygg; egen produksjonsworkflow.
- Manuell ekstern lenkekontroll.

Status beskriver implementasjon i repositoriet, ikke en uavhengig revisjon eller bekreftelse på at alle produksjonstjenester fungerer.

## Redaksjon og finansiering

Se [EDITORIAL.md](EDITORIAL.md) for arbeidsrutiner og [Om oss](https://tenkdegom.no/om/) for offentlige prinsipper. Ingen leverandør kan kjøpe en positiv vurdering. Relevante kommersielle forbindelser skal opplyses, og reklame skal merkes. Frivillige bidrag, digitale publikasjoner, kurs og tydelig merkede affiliate-lenker er mulige inntektskilder.

## Teknologi og drift

Astro og TypeScript med statisk frontend. Personvernsjekken behandler svar lokalt; kontaktskjemaet er et separat serveravhengig endepunkt (/api/contact). Ingen tredjepartssporing i nettstedets kildekode. Produksjonsworkflow bygger og tester før publisering via SSH med verifisert vertsnøkkel. Caddy og HTTPS inngår i driftsoppsettet.

## Lokal utvikling

Bruk Node.js 22 (som i CI) og npm.

```bash
npm install
npm run dev
```

Kontroller før PR:

```bash
npm run check
node --test tests/*.test.mjs
npm run build
```

Manuell kildekontroll:

```bash
npm run check:links
```

Undersøk DEAD og REVIEW manuelt. Ekstern lenkekontroll er ikke en obligatorisk CI-port og bekrefter ikke at kilden støtter en påstand. Lokal frontend alene gir ikke et fungerende kontaktendepunkt.

## Videre arbeid

Se [ROADMAP.md](ROADMAP.md). Oppdater redaksjonelle data og anbefalingskart sammen ved nye artikler. Ikke legg hemmeligheter eller personopplysninger i repositoriet.
