package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tobiaskohlbau/dustrobo/pkg/music"
)

type musicHandler struct {
	player *music.Player
}

func (m *musicHandler) Play(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to read request body").Error(), http.StatusInternalServerError)
		return
	}
	var rd struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &rd); err != nil {
		http.Error(w, errors.Wrap(err, "failed to unmarshal request").Error(), http.StatusBadRequest)
		return
	}
	log.Printf("got request for %s", rd.URL)

	res, err := http.Get(rd.URL)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to retrieve url").Error(), http.StatusBadRequest)
		return
	}
	if res.StatusCode != http.StatusOK {
		http.Error(w, errors.Wrap(err, "got bad response while retrieving url").Error(), http.StatusBadRequest)
		return
	}

	m.player, err = music.NewMP3Player(res.Body)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed creating mp3 player").Error(), http.StatusInternalServerError)
		return
	}

	m.player.Play()
}

func (m *musicHandler) Pause(w http.ResponseWriter, r *http.Request) {
	if err := m.player.Pause(); err != nil {
		http.Error(w, errors.Wrap(err, "failed to pause player").Error(), http.StatusInternalServerError)
		return
	}
}

func (m *musicHandler) Stop(w http.ResponseWriter, r *http.Request) {
	if err := m.player.Stop(); err != nil {
		http.Error(w, errors.Wrap(err, "failed to stop player").Error(), http.StatusInternalServerError)
		return
	}
}

func (m *musicHandler) Volume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to read request body").Error(), http.StatusInternalServerError)
		return
	}

	var rd struct {
		Volume float64
	}
	if err := json.Unmarshal(body, &rd); err != nil {
		http.Error(w, errors.Wrap(err, "failed to unmarshal request").Error(), http.StatusBadRequest)
		return
	}

	if m.player == nil {
		return
	}

	if err := m.player.Volume(rd.Volume); err != nil {
		http.Error(w, errors.Wrap(err, "failed to set volume").Error(), http.StatusBadRequest)
		return
	}
}
