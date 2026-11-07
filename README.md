# Tenk deg om

**Personvern angår oss alle.**

[Tenkdegom.no](https://tenkdegom.no) er et norsk folkeopplysningsprosjekt om digitalt personvern, eid og utgitt av **Ploos AS**.

## Formål

Vi vil gjøre digitalt personvern forståelig, relevant og praktisk for alle – uavhengig av alder og teknisk bakgrunn. Grunnleggende kunnskap og verktøy skal være gratis.

## Redaksjonelle prinsipper

- **Folkeopplysning først:** Saklige, forståelige og handlingsrettede råd.
- **Uavhengighet:** Ingen leverandør kan kjøpe seg til en bedre vurdering.
- **Åpenhet:** Vi opplyser om eierskap, finansiering og kommersielle forbindelser.
- **Personvern i praksis:** Dataminimering, ingen unødvendige brukerkontoer og ingen tredjepartssporing.
- **Tydelig reklame:** Eventuelle affiliate-lenker skal merkes tydelig.

Mulige inntektskilder er frivillige bidrag, affiliate-lenker, digitale kompendier og relevante kurs. Inntekter skal støtte driften uten å styre redaksjonelle vurderinger.

## Første versjon

Prosjektet er under utvikling. Første leveranse omfatter:

- Norsk forside og informasjon om prosjektet
- Fem kvalitetssikrede artikler (planlagt)
- Gratis personvernsjekk med svar behandlet lokalt i nettleseren (planlagt)
- Redaksjonelle prinsipper og åpen finansieringsinformasjon

**Merk:** Personvernsjekken er foreløpig en plassholderside, ikke et fungerende verktøy.

## Teknologi

- [Astro](https://astro.build/) og TypeScript
- Statiske sider som utgangspunkt
- Ingen database eller brukerkontoer i første versjon
- Mulig Go-backend senere ved behov
- Planlagt drift på europeisk VPS med Caddy og HTTPS

## Lokal utvikling

Krever Node.js og npm.

```bash
npm install
npm run dev
```

Bygg produksjonsfiler:

```bash
npm run build
```

Bygget er ikke verifisert i CI ennå.

## Veikart

Se [ROADMAP.md](ROADMAP.md) for milepæler og videre utvikling.

## Eier

**Ploos AS** – redaksjonelt ansvar og utgiver av Tenk deg om.
