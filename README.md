# TTK4145 Heisprosjekt

Dette prosjektet implementerer et sanntidssystem for styring av flere (N) heiser over flere (M) etasjer. Systemet fokuserer på distribusjon, feiltoleranse og robust håndtering av nettverksproblemer – uten bruk av en sentral server.

---

## Kom i gang

### Kjøring av programmet

Start en heisprosess med følgende kommando:

```bash
./start.bash -id=<your_id> -port=<your_port>
```

- `id`: En unik identifikator for heisen (f.eks. 1, 2, 3).
- `port`: Porten for lokal kommunikasjon mellom heisen og simulatoren (standard: 15657).

---

## Funksjonalitet

### Distribuert ordrehåndtering

- Heisene fordeler hall-ordre (opp/ned) dynamisk ved hjelp av en ekstern "Hall Request Assigner".
- Cab-ordre håndteres alltid av den heisen der knappen ble trykket.
- Ordrene blir bekreftet og synkronisert over nettverket.
- Når en heis utfører en ordre, fjernes den fra alle heisers datastruktur.

### Peer-to-peer kommunikasjon (UDP broadcast)

- Heisstatus sendes og mottas kontinuerlig via UDP broadcast.
- Verdensbildet oppdateres hos alle heiser basert på denne informasjonen.
- Lys og tilstand holdes synkronisert mellom heisene under normale forhold.

### Feilhåndtering

- Nettverksbrudd: Heisen går i "offline mode" og utfører kun hall-ordre den var kjent med før nettverksbruddet oppstod. Heisen kan fremdeles motta og utføre cab-ordre. Når alle ordre er fullført, forsøker den å koble seg på nytt.
- Bevegelsesfeil: Dersom en heis ikke registrerer ny etasje innen 10 sekunder mens den er i bevegelse, tolkes det som motorfeil og restart initieres.
- Autorestart: Ved programkrasj, strømbrudd eller motorfeil forsøker heisen å restarte seg selv automatisk og synkronisere med nettverket.

---

## Teknisk oversikt

### Heisens tilstand (`ElevatorInstance`)

- `Floor`: Nåværende etasje
- `Direction`: Retning (opp, ned, stopp)
- `State`: Idle, Moving, eller DoorOpen
- `Queue`: 3 x 4 matrise over ordrestatus (Uninitialized, No Order, Unconfirmed, Confirmed)

### Kommunikasjon

- UDP broadcast brukes for all meldingsutveksling.
- Peer discovery skjer automatisk via `peers`-modulen.
- Hver heis sender sin tilstand til de andre og mottar deres tilstand.

---

## Oppfylte krav

- Alle ordre blir utført, selv ved krasj, nettverksbrudd eller strømbrudd.
- Cab-ordre gjenopptas etter restart.
- Offline heis fungerer lokalt og håndterer egne ordre uavhengig av nettverket.
- Lyssystemet oppfører seg korrekt, og lys slukkes kun når ordren er utført.
- Døroppførsel følger krav: 3 sekunder åpen dør, og døren forblir åpen ved obstruksjon.

---

## Testing

Systemet er testet i miljø med:

- 3 heiser, 4 etasjer
- Nettverksbrudd, motorfeil og programkrasj
- Offline-mode med rekobling
- Synkronisering av ordre og lys
- Automatisk reinitialisering av heiser

---

## Bidragsytere

Prosjektet er levert anonymt i henhold til retningslinjene for TTK4145.

---