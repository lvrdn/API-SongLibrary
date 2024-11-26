package song

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// @Summary Get songs list with pagination and filtering by all fields
// @Description Get songs list with pagination and filtering by all fields
// @Tags songs
// @ID get-all
// @Param limit query int false "Limit" Default(0)
// @Param offset query int false "Offset" Default(0)
// @Param song query string false "song name"
// @Param group query string false "group name"
// @Param releaseDate query string false "year"
// @Param text query string false "text, word, letters"
// @Param link query string false "if need song with video use: true, else use:false"
// @Produce json
// @Success 200 {object} song.Response{response=[]song.Song}
// @Failure 400 {object} song.Response{error=song.Response{timestamp=time,message=string,path=string}}
// @Failure 500 "something bad with db or marshal/unmarshal data"
// @Router /api/songs [get]
func (sh *SongHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	var offset, limit int
	var err error
	if r.FormValue("limit") == "" {
		limit = 0
	} else {
		limit, err = strconv.Atoi(r.FormValue("limit"))
		if err != nil {
			dataToSend, err := GetAnswerWithError("limit must be number", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer limit read error [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}
	}
	if r.FormValue("offset") == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(r.FormValue("offset"))
		if err != nil {
			dataToSend, err := GetAnswerWithError("offset must be number", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer offset read error [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}
	}

	if offset < 0 || limit < 0 {
		dataToSend, err := GetAnswerWithError("offset or limit cannot be negative", r.URL.Path)
		if err != nil {
			log.Printf("marshal answer offset/limit negative error [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dataToSend)
		return
	}

	link := r.FormValue("link")
	if link != "false" && link != "true" && link != "" {
		dataToSend, err := GetAnswerWithError("link must be true or false", r.URL.Path)
		if err != nil {
			log.Printf("marshal answer bad link error [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dataToSend)
		return
	}

	//get songs
	songs, err := sh.Storage.GetAll(
		limit,
		offset,
		r.FormValue("song"),
		r.FormValue("group"),
		r.FormValue("releaseDate"),
		r.FormValue("text"),
		r.FormValue("link"),
	)

	if err != nil {
		if err.Error() == sh.Storage.GetErrorBadDate().Error() {
			dataToSend, err := GetAnswerWithError("bad date, date must have format year", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer with get all songs error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}

		log.Printf("get all song text error: [%s], data: [limit: %d, offset: %d, song: %s, group: %s, releaseDate: %s, text: %s, link: %s]\n",
			err.Error(), limit, offset, r.FormValue("song"), r.FormValue("group"), r.FormValue("releaseDate"), r.FormValue("text"), r.FormValue("link"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := Response{
		"response": songs,
	}

	dataResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("marshal response with songs error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(dataResponse)
	if err != nil {
		log.Printf("sending response error: [%s], user agent: [%s], path: [%s], method: [%s]\n", err.Error(), r.Header.Get("User-Agent"), r.URL.Path, r.Method)
	} else {
		log.Printf("response sended: [%s], user agent: [%s], path: [%s], method: [%s]\n", string(dataResponse), r.Header.Get("User-Agent"), r.URL.Path, r.Method)
	}
}

// @Summary Get song text with verse pagination
// @Description Get song text with verse pagination
// @Tags songs
// @ID get
// @Param limit query int false "Limit" Default(0)
// @Param offset query int false "Offset" Default(0)
// @Param id path int true "song id"
// @Produce json
// @Success 200 {object} song.Response{response=song.Response{id=int,verses=string,resesInSong=int}}
// @Failure 400 {object} song.Response{error=song.Response{timestamp=time,message=string,path=string}}
// @Failure 500 "something bad with db or marshal/unmarshal data"
// @Router /api/songs/{id} [get]
func (sh *SongHandler) Get(w http.ResponseWriter, r *http.Request) {
	var err error

	id := r.PathValue("id")
	song := &Song{}
	song.ID, err = strconv.Atoi(id)
	if err != nil {
		dataToSend, err := GetAnswerWithError("id must be number", r.URL.Path)
		if err != nil {
			log.Printf("marshal answer id read error [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dataToSend)
		return
	}

	var offset, limit int
	if r.FormValue("limit") == "" {
		limit = 0
	} else {
		limit, err = strconv.Atoi(r.FormValue("limit"))
		if err != nil {
			dataToSend, err := GetAnswerWithError("limit must be number", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer limit read error [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}
	}
	if r.FormValue("offset") == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(r.FormValue("offset"))
		if err != nil {
			dataToSend, err := GetAnswerWithError("offset must be number", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer offset read error [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}
	}

	if offset < 0 || limit < 0 {
		dataToSend, err := GetAnswerWithError("offset or limit cannot be negative", r.URL.Path)
		if err != nil {
			log.Printf("marshal answer offset/limit negative error [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dataToSend)
		return
	}

	//get text of song
	song.Text, err = sh.Storage.Get(song.ID)
	if err != nil {
		if err.Error() == sh.Storage.GetErrorNoRows().Error() {
			dataToSend, err := GetAnswerWithError("no song with this id", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer get song text error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}
		log.Printf("get song text error: [%s], data: [id: %d]\n", err.Error(), song.ID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	verses := strings.Split(song.Text, "\n\n")

	//get desired verses of the song text
	var desiredVerses []string
	switch {
	case offset == 0 && limit == 0:
		desiredVerses = verses
	case offset+limit <= len(verses):
		desiredVerses = verses[offset : offset+limit]
	case offset < len(verses):
		desiredVerses = verses[offset:]
	}

	response := Response{
		"response": Response{
			"id":           song.ID,
			"verses":       strings.Join(desiredVerses, "\n\n"),
			"versesInSong": len(verses),
		},
	}

	dataResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("marshal response with verse error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(dataResponse)
	if err != nil {
		log.Printf("sending response error: [%s], user agent: [%s], path: [%s], method: [%s]\n", err.Error(), r.Header.Get("User-Agent"), r.URL.Path, r.Method)
	} else {
		log.Printf("response sended: [%s], user agent: [%s], path: [%s], method: [%s]\n", string(dataResponse), r.Header.Get("User-Agent"), r.URL.Path, r.Method)
	}
}

// @Summary Add new song to library
// @Description Add new song to library
// @Tags songs
// @ID new
// @Param bodyJSON body string true "song and group names" SchemaExample({"song":"sone song name","group":"some group name"})
// @Accept json
// @Produce json
// @Success 201 {object} song.Response{response=song.Response{id=int}}
// @Failure 400 {object} song.Response{error=song.Response{timestamp=time,message=string,path=string}}
// @Failure 500 "something bad with db or marshal/unmarshal data"
// @Failure 502 "something bad with another server for getting song info"
// @Router /api/songs [put]
func (sh *SongHandler) New(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("request body read error: [%s], user agent: [%s], path: [%s], method: [%s]\n", err.Error(), r.Header.Get("User-Agent"), r.URL.Path, r.Method)

		dataToSend, err := GetAnswerWithError("request body read error", r.URL.Path)
		if err != nil {
			log.Printf("marshal answer body read error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dataToSend)
		return
	}

	song := &Song{}

	err = json.Unmarshal(body, song)
	if err != nil {
		log.Printf("unmarshal body error: [%s], user agent: [%s], path: [%s], method: [%s]\n", err.Error(), r.Header.Get("User-Agent"), r.URL.Path, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if song.Group == "" || song.Name == "" {
		dataToSend, err := GetAnswerWithError("song and group values must be not empty", r.URL.Path)
		if err != nil {
			log.Printf("marshal answer error with empty song, group names: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dataToSend)
		return
	}

	//prepare to send request to another server
	searcherParams := url.Values{}
	searcherParams.Add("group", song.Group)
	searcherParams.Add("song", song.Name)

	newRequest, err := http.NewRequest("GET", sh.ExternalAPI+"/info"+"?"+searcherParams.Encode(), nil)
	if err != nil {
		log.Printf("creating new request error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("getting response error from server [%s]: error [%s], path: [%s]\n", sh.ExternalAPI, err.Error(), newRequest.URL.Path)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("bad status code from server [%s]: code [%d], path: [%s]\n", sh.ExternalAPI, resp.StatusCode, newRequest.URL.Path)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("request body read error from server [%s]: error [%s], path: [%s], method: [%s]\n", sh.ExternalAPI, err.Error(), newRequest.URL.Path, newRequest.Method)

		dataToSend, err := GetAnswerWithError("request body read error from server", sh.ExternalAPI)
		if err != nil {
			log.Printf("marshal answer body read error from server [%s]: error: [%s], path: [%s], method: [%s]\n", sh.ExternalAPI, err.Error(), r.URL.Path, r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		w.Write(dataToSend)
		return
	}

	m := make(map[string](map[string]string))

	err = json.Unmarshal(respBody, &m)
	if err != nil {
		log.Printf("unmarshal body error from server [%s]: error [%s], path: [%s], method: [%s]\n", sh.ExternalAPI, err.Error(), newRequest.URL.Path, newRequest.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	song.ReleaseDate = m["SongDetail"]["releaseDate"]
	song.Text = m["SongDetail"]["text"]
	song.Link = m["SongDetail"]["link"]

	//add new song to storage
	song.ID, err = sh.Storage.Add(song.Name, song.Group, song.ReleaseDate, song.Text, song.Link)

	if err != nil {
		if err.Error() == sh.Storage.GetErrorAlreadyExist().Error() {
			dataToSend, err := GetAnswerWithError("song of this group is already added", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer add song error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}

		log.Printf("add new song error: [%s], data: [song: %s, group: %s, date: %s, link: %s]\n", err.Error(), song.Name, song.Group, song.ReleaseDate, song.Link)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := Response{
		"response": Response{
			"id": song.ID,
		},
	}

	dataResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("marshal answer with new song id error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(dataResponse)
	if err != nil {
		log.Printf("sending response error: [%s], user agent: [%s], path: [%s], method: [%s]\n", err.Error(), r.Header.Get("User-Agent"), r.URL.Path, r.Method)
	} else {
		log.Printf("response sended: [%s], user agent: [%s], path: [%s], method: [%s]\n", string(dataResponse), r.Header.Get("User-Agent"), r.URL.Path, r.Method)
	}
}

// @Summary Update song data
// @Description Update song data
// @Tags songs
// @ID update
// @Param bodyJSON body string true "song id and at least one of the listed parameters required" SchemaExample({"id":2,"releaseDate":"25.02.2012","text":"some text","link":"some link"})
// @Success 200
// @Failure 400 {object} song.Response{error=song.Response{timestamp=time,message=string,path=string}}
// @Failure 500 "something bad with db or marshal/unmarshal data"
// @Router /api/songs [post]
func (sh *SongHandler) Update(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("request body read error: [%s], user agent: [%s], path: [%s], method: [%s]\n", err.Error(), r.Header.Get("User-Agent"), r.URL.Path, r.Method)

		dataToSend, err := GetAnswerWithError("request body read error", r.URL.Path)
		if err != nil {
			log.Printf("marshal answer body read error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dataToSend)
		return
	}

	song := &Song{}

	err = json.Unmarshal(body, song)
	if err != nil {
		log.Printf("unmarshal body error: [%s], user agent: [%s], path: [%s], method: [%s]\n", err.Error(), r.Header.Get("User-Agent"), r.URL.Path, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//update song data (link and/or release date and/or text)
	err = sh.Storage.Update(song.ID, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		if err.Error() == sh.Storage.GetErrorNoUpdate().Error() {
			dataToSend, err := GetAnswerWithError("no data to update, release date or link or text must be not empty", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer with update song error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}

		if err.Error() == sh.Storage.GetErrorBadID().Error() {
			dataToSend, err := GetAnswerWithError("bad id, nothing updated", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer with update song error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}

		if err.Error() == sh.Storage.GetErrorBadDate().Error() {
			dataToSend, err := GetAnswerWithError("bad date, date must have format day.month.year", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer with update song error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}

		log.Printf("update song error: [%s], data: [id: %d, text: %s, date: %s, link: %s]\n", err.Error(), song.ID, song.Text, song.ReleaseDate, song.Link)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("song updated, id: [%d], user agent: [%s], path: [%s], method: [%s]\n", song.ID, r.Header.Get("User-Agent"), r.URL.Path, r.Method)
}

// @Summary Delete song
// @Description Delete song
// @Tags songs
// @ID delete
// @Param bodyJSON body string true "song id" SchemaExample({"id":2})
// @Success 200
// @Failure 400 {object} song.Response{error=song.Response{timestamp=time,message=string,path=string}}
// @Failure 500 "something bad with db or marshal/unmarshal data"
// @Router /api/songs [delete]
func (sh *SongHandler) Delete(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("request body read error: [%s], user agent: [%s], path: [%s], method: [%s]\n", err.Error(), r.Header.Get("User-Agent"), r.URL.Path, r.Method)

		dataToSend, err := GetAnswerWithError("request body read error", r.URL.Path)
		if err != nil {
			log.Printf("marshal answer body read error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dataToSend)
		return
	}

	song := &Song{}

	err = json.Unmarshal(body, song)
	if err != nil {
		log.Printf("unmarshal body error: [%s], user agent: [%s], path: [%s], method: [%s]\n", err.Error(), r.Header.Get("User-Agent"), r.URL.Path, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//delete song from storage
	err = sh.Storage.Delete(song.ID)
	if err != nil {
		if err.Error() == sh.Storage.GetErrorBadID().Error() {
			dataToSend, err := GetAnswerWithError("bad id, nothing deleted", r.URL.Path)
			if err != nil {
				log.Printf("marshal answer with delete song error: [%s], path: [%s], method: [%s]\n", err.Error(), r.URL.Path, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(dataToSend)
			return
		}

		log.Printf("delete song error: [%s] [id: %d]\n", err.Error(), song.ID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("song deleted, id: [%d], user agent: [%s], path: [%s], method: [%s]\n", song.ID, r.Header.Get("User-Agent"), r.URL.Path, r.Method)
}
