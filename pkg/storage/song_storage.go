package storage

import (
	"SongLibrary/pkg/song"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		DB: db,
	}

}

func (s *Storage) GetAll(limit, offset int, songNameFragment, groupNameFragment, year, textFragment, linkExist string) ([]*song.Song, error) {

	var songs []*song.Song
	query := "SELECT * FROM songs "
	placeholderNum := 1
	args := make([]interface{}, 0)

	if songNameFragment != "" {
		if placeholderNum == 1 {
			query += "WHERE "
		}
		query += fmt.Sprintf("song_name LIKE CONCAT('%%',$%d::text,'%%') AND ", placeholderNum)
		placeholderNum++
		args = append(args, songNameFragment)
	}

	if groupNameFragment != "" {
		if placeholderNum == 1 {
			query += "WHERE "
		}
		query += fmt.Sprintf("group_name LIKE CONCAT('%%',$%d::text,'%%') AND ", placeholderNum)
		placeholderNum++
		args = append(args, groupNameFragment)
	}

	if year != "" {
		if placeholderNum == 1 {
			query += "WHERE "
		}

		startDate, err := time.Parse("2006", year)
		if err != nil {
			return nil, s.GetErrorBadDate()
		}
		endDate, _ := time.Parse("02.01.2006", "31.12."+year)

		query += fmt.Sprintf("release_date >= $%d AND release_date <= $%d AND ", placeholderNum, placeholderNum+1)
		placeholderNum += 2
		args = append(args, startDate, endDate)
	}

	if textFragment != "" {
		if placeholderNum == 1 {
			query += "WHERE "
		}
		query += fmt.Sprintf("text LIKE CONCAT('%%',$%d::text,'%%') AND ", placeholderNum)
		placeholderNum++
		args = append(args, textFragment)
	}

	if strings.ToLower(linkExist) == "true" {
		if placeholderNum == 1 {
			query += "WHERE "
		}
		query += "link is not null AND "

	} else if strings.ToLower(linkExist) == "false" {
		if placeholderNum == 1 {
			query += "WHERE "
		}
		query += "link is null AND "
	}

	query = strings.TrimSuffix(query, "AND ")
	query += "ORDER BY id "

	if limit > 0 {
		query += fmt.Sprintf("LIMIT $%d ", placeholderNum)
		placeholderNum++
		args = append(args, limit)
	}

	if offset > 0 {
		query += fmt.Sprintf("OFFSET $%d ", placeholderNum)
		placeholderNum++
		args = append(args, offset)
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		log.Printf("method get all query error: [%s], query: [%s], args: [%v]\n", err.Error(), query, args)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var songName, groupName, releaseDate string
		var text, link sql.NullString
		var date sql.NullTime
		var id int
		err := rows.Scan(&id, &songName, &groupName, &date, &text, &link)
		if err != nil {
			log.Printf("method get all scan error: [%s], query: [%s]\n", err.Error(), query)
			return nil, err
		}

		if date.Valid {
			releaseDate = date.Time.Format("02.01.2006")
		} else {
			releaseDate = ""
		}

		songs = append(songs, &song.Song{
			ID:          id,
			Name:        songName,
			Group:       groupName,
			ReleaseDate: releaseDate,
			Text:        text.String,
			Link:        link.String,
		})
	}

	return songs, nil
}

func (s *Storage) Get(id int) (string, error) {

	var text sql.NullString
	err := s.DB.QueryRow(
		`SELECT text FROM songs WHERE id = $1`, id,
	).Scan(&text)

	if err != nil {
		if err.Error() != s.GetErrorNoRows().Error() {
			log.Printf("method get query error: [%s], id: [%d]\n", err.Error(), id)
		}
		return "", err
	}

	return text.String, nil
}

func (s *Storage) Add(name, group, releaseDate, text, link string) (int, error) {

	var insertID int
	date, err := time.Parse("02.01.2006", releaseDate)
	if err != nil {
		return 0, err
	}

	err = s.DB.QueryRow(
		`INSERT INTO 
	songs(song_name,group_name,release_date,text,link) 
	VALUES($1,$2,$3,$4,$5) 
	RETURNING id`,
		name, group, date, text, link,
	).Scan(&insertID)

	if err != nil {
		if err.Error() != s.GetErrorAlreadyExist().Error() {
			log.Printf("method add query error: [%s], args: [song: %s, group: %s, releaseDate: %s, link: %s]\n", err.Error(), name, group, releaseDate, link)
		}
		return 0, err
	}

	return insertID, nil
}

func (s *Storage) Delete(id int) error {

	result, err := s.DB.Exec(
		`DELETE FROM songs WHERE id = $1`, id,
	)
	if err != nil {
		log.Printf("method delete query error: [%s], id: [%d]\n", err.Error(), id)
		return err
	}

	num, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if num == 0 {
		return s.GetErrorBadID()
	}

	return nil
}

func (s *Storage) Update(id int, releaseDate, text, link string) error {

	query := "UPDATE songs SET "
	placeholderNum := 1
	args := make([]interface{}, 0)

	if releaseDate != "" {
		date, err := time.Parse("02.01.2006", releaseDate)
		if err != nil {
			return s.GetErrorBadDate()
		}
		query += fmt.Sprintf("release_date = $%d, ", placeholderNum)
		placeholderNum++
		args = append(args, date)
	}

	if text != "" {
		query += fmt.Sprintf("text = $%d, ", placeholderNum)
		placeholderNum++
		args = append(args, text)
	}

	if link != "" {
		query += fmt.Sprintf("link = $%d, ", placeholderNum)
		placeholderNum++
		args = append(args, link)
	}

	if placeholderNum == 1 {
		return s.GetErrorNoUpdate()
	}

	query = strings.TrimSuffix(query, ", ")

	query += fmt.Sprintf(" WHERE id = $%d", placeholderNum)
	args = append(args, id)

	result, err := s.DB.Exec(query, args...)
	if err != nil {
		log.Printf("method update query error: [%s], query: [%s], args: [%v]\n", err.Error(), query, args)
		return err
	}

	num, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if num == 0 {
		return s.GetErrorBadID()
	}

	return nil
}

func (s *Storage) GetErrorAlreadyExist() error {
	return fmt.Errorf(`pq: duplicate key value violates unique constraint "songs_song_name_group_name_key"`)
}

func (s *Storage) GetErrorNoUpdate() error {
	return fmt.Errorf("no data to update")
}

func (s *Storage) GetErrorBadID() error {
	return fmt.Errorf("bad id")
}

func (s *Storage) GetErrorBadDate() error {
	return fmt.Errorf("bad date format")
}

func (s *Storage) GetErrorNoRows() error {
	return fmt.Errorf("sql: no rows in result set")
}
