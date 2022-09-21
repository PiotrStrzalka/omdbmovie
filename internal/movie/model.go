package movie

import (
	"strconv"
	"strings"
)

type MovieRecord []string

func (m *MovieRecord) Id() *string {
	return &(*m)[0]
}

func (m *MovieRecord) TitleType() *string {
	return &(*m)[1]
}

func (m *MovieRecord) PrimaryTitle() *string {
	return &(*m)[2]
}

func (m *MovieRecord) OriginalTitle() *string {
	return &(*m)[3]
}

func (m *MovieRecord) StartYear() (int, error) {
	return strconv.Atoi((*m)[5])
}

func (m *MovieRecord) EndYear() (int, error) {
	return strconv.Atoi((*m)[6])
}

func (m *MovieRecord) Genres() []string {
	return strings.Split((*m)[8], ",")
}

func (m *MovieRecord) RuntimeMinutes() (int, error) {
	return strconv.Atoi((*m)[7])
}

type MovieResult struct {
	Id    string
	Title string
	Plot  string
}
