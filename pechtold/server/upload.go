package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"image/png"
	"log"
	"net"
	"net/http"

	"github.com/GeenPeil/teken/data"
	"github.com/GeenPeil/teken/storage"
	"github.com/lib/pq"
)

var (
	fieldErr   = "form values missing or invalid"
	captchaErr = "captcha invalid"
	imgErr     = "image invalid"
)

func (s *Server) newUploadHandlerFunc() http.HandlerFunc {

	stmtInsertNawHash, err := s.db.Prepare(`INSERT INTO nawhashes (hash) VALUES ($1)`)
	if err != nil {
		log.Fatalf("error preparing stmtInsertNAWHash: %v", err)
	}

	stmtInsertHandtekening, err := s.db.Prepare(`INSERT INTO handtekeningen (insert_time, iphash) VALUES (NOW(), $1) RETURNING ID`)
	if err != nil {
		log.Fatalf("error preparing stmtInsertHandtekening: %v", err)
	}

	saver, err := storage.NewSaver(s.options.StoragePubkeyFile, s.options.StorageLocation)
	if err != nil {
		log.Fatalf("error creating storage.Saver: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h := &data.Handtekening{}
		err := json.NewDecoder(r.Body).Decode(h)
		if err != nil {
			http.Error(w, "input json error", http.StatusInternalServerError)
			log.Printf("error decoding json in request from %s: %v", r.RemoteAddr, err)
		}

		out := &uploadOutput{}

		{
			remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)

			// check captcha
			if !s.options.CaptchaDisable {
				valid, err := s.captcha.Verify(h.CaptchaResponse, remoteIP)
				if err != nil {
					http.Error(w, "server error", http.StatusInternalServerError)
					log.Printf("error verifying captcha: %v", err)
					return
				}
				if !valid {
					out.Error = captchaErr
					log.Printf("invalid captcha in request from %s", r.RemoteAddr)
					goto Response
				}
			}

			// checken of alle data is ingevuld
			if len(h.Voornaam) == 0 {
				out.Error = fieldErr
				log.Printf("missing field voornaam in request from %s", r.RemoteAddr)
				goto Response
			}

			if len(h.Achternaam) == 0 {
				out.Error = fieldErr
				log.Printf("missing field achternaam in request from %s", r.RemoteAddr)
				goto Response
			}

			if len(h.Geboortedatum) == 0 {
				out.Error = fieldErr
				log.Printf("missing field geboortedatum in request from %s", r.RemoteAddr)
				goto Response
			}

			if len(h.Geboorteplaats) == 0 {
				out.Error = fieldErr
				log.Printf("missing field geboorteplaats in request from %s", r.RemoteAddr)
				goto Response
			}

			if len(h.Straat) == 0 {
				out.Error = fieldErr
				log.Printf("missing field straat in request from %s", r.RemoteAddr)
				goto Response
			}

			if len(h.Huisnummer) == 0 {
				out.Error = fieldErr
				log.Printf("missing field huisnummer in request from %s", r.RemoteAddr)
				goto Response
			}

			if len(h.Postcode) == 0 {
				out.Error = fieldErr
				log.Printf("missing field postcode in request from %s", r.RemoteAddr)
				goto Response
			}

			if len(h.Woonplaats) == 0 {
				out.Error = fieldErr
				log.Printf("missing field woonplaats in request from %s", r.RemoteAddr)
				goto Response
			}

			if len(h.Handtekening) == 0 {
				out.Error = fieldErr
				log.Printf("missing field handtekening in request from %s", r.RemoteAddr)
				goto Response
			}

			// check (decode, etc.) handtekening
			hImgPNG := make([]byte, base64.StdEncoding.DecodedLen(len(h.Handtekening)))
			_, err = base64.StdEncoding.Decode(hImgPNG, h.Handtekening)
			if err != nil {
				out.Error = imgErr
				log.Printf("invalid base64 image received from %s: %v", r.RemoteAddr, err)
				goto Response
			}

			_, err = png.Decode(bytes.NewBuffer(hImgPNG))
			if err != nil {
				out.Error = imgErr
				log.Printf("invalid image from %s: %v", r.RemoteAddr, err)
				goto Response
			}

			// all ok
			out.Success = true

			// naw hash check (false positive)
			nawHash := sha256.New()
			nawHash.Write([]byte(h.Voornaam))
			nawHash.Write([]byte(h.Achternaam))
			nawHash.Write([]byte(h.Geboortedatum))
			nawHash.Write([]byte(h.Geboorteplaats))
			nawHash.Write([]byte(h.Straat))
			nawHash.Write([]byte(h.Huisnummer))
			nawHash.Write([]byte(h.Postcode))
			nawHash.Write([]byte(h.Woonplaats))
			nawHashBytes := nawHash.Sum(nil)
			_, err = stmtInsertNawHash.Exec(nawHashBytes)
			if err != nil {
				if perr, ok := err.(*pq.Error); ok {
					if perr.Code == "23505" {
						log.Printf("duplicate n.a.w. hash from %s: %x", r.RemoteAddr, nawHashBytes)
						goto Response // return direclty with a 'false positive'
					} else {
						log.Printf("error inserting naw hash from %s: %v", r.RemoteAddr, err)
						http.Error(w, "server error", http.StatusInternalServerError)
						return
					}
				}
			}

			// insert handtekening entry into db, get inserted ID
			ipHash := sha256.New()
			ipHashBytes := ipHash.Sum([]byte(remoteIP))
			insertHandtekeningRows, err := stmtInsertHandtekening.Query(ipHashBytes)
			if err != nil {
				log.Printf("error inserting handtekening entry in db for %s: %v", r.RemoteAddr, err)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			insertHandtekeningRows.Next()
			var ID int
			err = insertHandtekeningRows.Scan(&ID)
			if err != nil {
				log.Printf("error getting ID for new handtekening entry for %s: %v", r.RemoteAddr, err)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			// save to disk
			err = saver.Save(ID, h)
			if err != nil {
				log.Printf("error saving handtekening for %s: %v", r.RemoteAddr, err)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
		}

	Response:
		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			log.Printf("error encoding response json: %v", err)
		}
	}
}
