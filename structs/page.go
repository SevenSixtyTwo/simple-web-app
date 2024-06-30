package structs

import "os"

var pagesPath = "./data/"

type Page struct {
	Title string
	Body  []byte
}

// 0600 means that the file should be created with read-write permissions for the current user only
func (p *Page) Save() error {
	filename := p.Title + ".txt"
	filePath := pagesPath + filename

	return os.WriteFile(filePath, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	filePath := pagesPath + filename

	body, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}
