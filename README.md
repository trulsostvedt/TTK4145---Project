# TTK4145---Project

Kjør programmet:
```
go run main.go -id=your_id
```

Dette prosjektet tar for seg sanntidsprogrammering av N heiser som oppererer over M etasjer. 


Peer to Peer

- Cyclic counter (bekrefta/ikke-bekrefta/ingen ordre)
    - Lys på når bekrefta tilstand
- UDP Broadcast heartbeat (felles verdensbilde)
- Når noen forsvinntrulsostvedt/TTK4145-projecter trenger ingen å ta over
- TCP forbindelser mellom seg
- Broadcast, jeg vet dette, okei nå vet jeg det, okei nå vet alle det, da blir et kall bekrefta. Siden alle har samme verdensbilde så vil alle komme til at den samme personen skal ta oppdraget. Den utfører oppdraget, når utført broadcaster den ingen ordre, eller fjerner den bekrefta ordren, og alle skjønner at den er løst, siden det er ingen andre årsaker til at en bekrefta ordre skal forsvinne. 
- Når nettverksfeil/crash så må den forholde seg til de andre heisenes broadcast ved oppstart. 
- Jeg fikk inn en ny ordre i 3. etasje, *plutselig ny trykk*, nå har man to meldinger man må resende, da må man bruke TCP(?) (kanskje snakk om dataspill skjønte han ikke helt)
- 
-  Den structen I broadcast, hva skal egentlig være i det verdensbilde? Hva trenger å være der, ikke putt inn noe som ikke trengs å være der, json eller noe sånt. 
    - Etasje jeg er 
    - Om jeg beveger meg
    - Cyclic counter, finnes tilfeller der du kanskje trenger 4 argumenter.. 
    - 3x4 cyclic counter
- Ikke se bort i fra at om generalisering av etasjer og heiser gjør koden bedre, men det er ikke noe spesielt viktig for Sverre eller Anders


Overordna struktur/plan: 
- Start broadcaster
- Så er det å starte et program 
- Så er det å teste, «Jeg har mottatt en sånn melding» 
- Kjempebra nå har vi kommunikasjon
- Så kan vi begynne med ordentlige data (etasje div osv)
- Så merge funksjon som sier at jeg fikk inn hans verdensbilde, noe jeg må merke meg? 
- Koden som går over og teller har alle en ubekrefta ordre? Den må ligge noe sted
- Så er det heis, nå kan vi begynne å regne ut, jeg som heis, hva skal jeg gjøre nå?
