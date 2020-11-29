package kunst

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/dbscan"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/ihleven/errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRepo(database string) (*Repo, error) {

	ctx := context.Background()
	dbpool, err := pgxpool.Connect(ctx, fmt.Sprintf("postgresql://%s@localhost:5432/%s", database, database))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to database: %q", database)
	}

	repo := Repo{
		ctx:    ctx,
		dbpool: dbpool,
	}

	return &repo, nil
}

type Repo struct {
	ctx    context.Context
	dbpool *pgxpool.Pool
}

func (r *Repo) Select(dst interface{}, query string, args ...interface{}) error {

	err := pgxscan.Select(r.ctx, r.dbpool, dst, query, args...)
	if err != nil {
		// if errors.As(err, &pgx.ErrNoRows) {
		// 	return errors.NewWithCode(errors.NotFound, "Not found: %s", query)
		// }
		return errors.Wrap(err, "Error")
	}
	return nil
}

func (r *Repo) LoadAusstellungen() ([]Ausstellung, error) {

	stmt := "SELECT * FROM ausstellung ORDER BY id ASC"

	var ausstellungen []Ausstellung
	err := r.Select(&ausstellungen, stmt)
	fmt.Println(ausstellungen, err)
	if err != nil {
		return nil, err
	}

	return ausstellungen, nil
}

func (r *Repo) InsertAusstellung(a *Ausstellung) (int, error) {

	stmt := "INSERT INTO ausstellung (titel, untertitel, typ, jahr, von, bis, ort, venue, kommentar) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id"

	row := r.dbpool.QueryRow(r.ctx, stmt, a.Titel, a.Untertitel, a.Typ, a.Jahr, a.Von, a.Bis, a.Ort, a.Venue, a.Kommentar)
	var returnid int
	err := row.Scan(&returnid)
	return returnid, errors.Wrap(err, "Could not insert ausstellung %v", a)
}

func (r *Repo) SaveAusstellung(a *Ausstellung) error {
	fmt.Println("ausstellung:", a)
	stmt := `
		UPDATE ausstellung SET titel=$1, untertitel=$2, typ=$3, jahr=$4, von=$5, bis=$6, ort=$7, venue=$8, kommentar=$9 WHERE id=$10
	`
	_, err := r.dbpool.Exec(r.ctx, stmt, a.Titel, a.Untertitel, a.Typ, a.Jahr, a.Von, a.Bis, a.Ort, a.Venue, a.Kommentar, a.ID)
	return errors.Wrap(err, "Could not save ausstellung %v", a)
}

// LoadAusstellung ...
func (r *Repo) LoadAusstellung(id int) (*Ausstellung, error) {

	stmt := "SELECT * FROM ausstellung WHERE id=$1"
	fmt.Println("LoadAusstellung")

	ausstellung := Ausstellung{Bilder: make([]Bild, 0), Fotos: make([]Foto, 0)}
	err := pgxscan.Get(r.ctx, r.dbpool, &ausstellung, stmt, id)

	if dbscan.NotFound(err) {
		fmt.Println(dbscan.NotFound(err))
		return nil, errors.NewWithCode(404, "not row found")
	} else if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	fmt.Println("LoadAusstellung", id, ausstellung)
	err = r.Select(&ausstellung.Bilder, "SELECT b.* FROM bild b, enthalten e WHERE b.id=e.bild_id AND ausstellung_id=$1", id)
	if err != nil {
		return nil, errors.Wrap(err, "db error bilder")
	}

	// err = r.Select(&ausstellung.Fotos, "SELECT * FROM foto WHERE aust_id=$1", id)
	// if err != nil {
	// 	return nil, err
	// }
	return &ausstellung, nil
}

func (r *Repo) LoadSerien() ([]Serie, error) {

	stmt := "SELECT * FROM serie ORDER BY id ASC"

	var serien []Serie
	err := r.Select(&serien, stmt)
	fmt.Println(serien, err)
	if err != nil {
		return nil, err
	}

	return serien, nil
}

func (r *Repo) InsertSerie(s *Serie) (int, error) {

	stmt := "INSERT INTO serie (jahr, titel, anzahl, kommentar) VALUES ($1,$2,$3,$4) RETURNING id"

	row := r.dbpool.QueryRow(r.ctx, stmt, s.Jahr, s.Titel, s.Anzahl, s.Kommentar)
	var returnid int
	err := row.Scan(&returnid)
	return returnid, errors.Wrap(err, "Could not insert serie %v", s)
}

// LoadSerie ...
func (r *Repo) LoadSerie(id int) (*Serie, error) {

	stmt := "SELECT * FROM serie WHERE id=$1"

	serie := Serie{Bilder: make([]Bild, 0), Fotos: make([]Foto, 0)}
	err := pgxscan.Get(r.ctx, r.dbpool, &serie, stmt, id)
	if err != nil {
		return nil, err
	}

	err = r.Select(&serie.Bilder, "SELECT * FROM bild WHERE serie=$1", serie.Titel)
	if err != nil {
		return nil, err
	}

	// err = r.Select(&serie.Fotos, "SELECT * FROM foto WHERE aust_id=$1", id)
	// if err != nil {
	// 	return nil, err
	// }
	return &serie, nil
}

func (r *Repo) LoadBilder() ([]Bild, error) {

	stmt := "SELECT * FROM bild ORDER BY id DESC"

	var bilder []Bild
	err := r.Select(&bilder, stmt)
	fmt.Println(bilder, err)
	if err != nil {
		return nil, err
	}

	return bilder, nil
}

func (r *Repo) LoadBild(id int) (*Bild, error) {

	stmt := "SELECT * FROM bild WHERE id=$1"

	bild := Bild{Fotos: make([]Foto, 0)}
	err := pgxscan.Get(r.ctx, r.dbpool, &bild, stmt, id)
	if err != nil {
		return nil, err
	}

	err = r.Select(&bild.Fotos, "SELECT * FROM foto WHERE bild_id=$1", id)
	if err != nil {
		return nil, err
	}
	return &bild, nil
}

func (r *Repo) InsertBild(bild *Bild) (int, error) {

	stmt := "INSERT INTO bild (titel, jahr, technik, traeger, hoehe, breite, tiefe, flaeche, foto_id, anmerkungen, kommentar, ordnung, phase) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id"

	row := r.dbpool.QueryRow(r.ctx, stmt, bild.Titel, bild.Jahr, bild.Technik, bild.Bildträger, bild.Höhe, bild.Breite, bild.Tiefe, bild.Fläche, bild.IndexFotoID, bild.Anmerkungen, bild.Kommentar, bild.Überordnung, bild.Schaffensphase)
	var returnid int
	err := row.Scan(&returnid)
	return returnid, errors.Wrap(err, "Could not insert bild %v", bild)
}

func (r *Repo) SaveBild(id int, bild *Bild) error {

	stmt := "UPDATE bild set  titel=$2, jahr=$3, technik=$4, traeger=$5, hoehe=$6, breite=$7, tiefe=$8, flaeche=$9, foto_id=$10, anmerkungen=$11, kommentar=$12, ordnung=$13, phase=$14 where id=$1"

	i, err := r.dbpool.Exec(r.ctx, stmt, id, bild.Titel, bild.Jahr, bild.Technik, bild.Bildträger, bild.Höhe, bild.Breite, bild.Tiefe, bild.Fläche, bild.IndexFotoID, bild.Anmerkungen, bild.Kommentar, bild.Überordnung, bild.Schaffensphase)
	if err != nil {
		return err
	}
	fmt.Println(bild.Titel)

	fmt.Println("i:", i, err)
	return err
}

func (r *Repo) InsertFoto(id, index int, name string, size int, path, format string, width, height int, taken time.Time, caption, kommentar string) (int, error) {

	stmt := "INSERT INTO foto (bild_id,index,name,size,uploaded,path,format,width,height,taken,caption,kommentar) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id"

	row := r.dbpool.QueryRow(r.ctx, stmt, id, index, name, size, time.Now(), path, format, width, height, taken, caption, kommentar)
	var returnid int
	err := row.Scan(&returnid)
	return returnid, errors.Wrap(err, "Could not insert foto %v", id)

}

func (r *Repo) LoadFoto(id int) (*Foto, error) {

	stmt := "SELECT * FROM foto WHERE id=$1"

	var foto Foto
	err := pgxscan.Get(r.ctx, r.dbpool, &foto, stmt, id)
	if err != nil {
		return nil, err
	}

	return &foto, nil
}

// func (r *Repo) LoadFotos() ([]Foto, error) {

// 	stmt := "SELECT * FROM foto ORDER BY id"

// 	var fotos []Foto
// 	err := pgxscan.Select(r.ctx, r.dbpool, &fotos, stmt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return fotos, nil
// }