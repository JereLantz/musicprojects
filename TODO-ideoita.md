- API:t sessio tarkistuksen taakse.
- Lisää conffitiedosto, jossa voidaa muokata mm. sessioiden pituus, sessioiden cleanup funtion
kutsujen välistä aikaa, portti jota tämä palvelin kuuntelee.
- Sessioiden tallennus tietokantaan.
- Sessiot api:n taakse.
- Ota sessioiden luonnissa IP osoite talteen ja tarkista että se pysyy samana seuraavissa requesteissa. Resetoi kirjautuneiden käyttäjien sessio jos ip osoite vaihtuu.
- Tapa tallentaa käyttäjien state.
   - Tallenna rekisteröityjen käyttäjien state tietokantaan.
