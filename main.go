package main

import (
	"fmt"
	"net/http"
	"regexp"

	"./models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// func parseAllFiles([]string, error) {
// 	serchDir := "templates/*"

// 	fileList := []string{}
// 	filepath.Walk(serchDir, func(path string, f os.FileInfo, err error) error {
// 		fileList = append(fileList, path)
// 		return nil
// 	})

// 	for _, file := range fileList {
// 		fmt.Println(file)
// 	}
// }

// Тут временно хранятся все посты
var posts map[string]*models.Post

// indexHandler показать главную
func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println(w, err.Error())
	}

	// fmt.Println(posts)

	t.ExecuteTemplate(w, "index", posts)
}

// writeHandler показать страницу редактирования\создания
func writeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println(w, err.Error())
	}
	t.ExecuteTemplate(w, "write", nil)
}

// editHandler изменить существующий пост
func editHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println(w, err.Error())
	}

	id := r.FormValue("id")
	post, found := posts[id]
	if !found {
		http.NotFound(w, r)
	}

	t.ExecuteTemplate(w, "write", post)
}

// savePostHandler сохранить новый пост
func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	re := regexp.MustCompile(`[[:punct:]]|[[:space:]]`) // Обрезает странные символы и пробелы
	strTitle := re.ReplaceAllString(title, "")
	strContent := re.ReplaceAllString(content, "")

	var post *models.Post // Найденный или новый пост

	/*
	* Если строка не пустая, то редактировать
	* Если пустая, то создать новый пост
	 */
	if len(id) != 0 {
		post = posts[id] // Выбрать id

		if strTitle != "" && strContent != "" {
			post.Title = title     // Заменить старый title на новый
			post.Content = content // Заменить старый content на новый
		} else {
			http.Error(w, "Error", 302)
		}
	} else {
		id = GenereteID()                          // Сгенирировать id
		post := models.NewPost(id, title, content) // Создать новый пост
		posts[post.ID] = post                      // Добавить пост в map
	}

	http.Redirect(w, r, "/", 302)
}

// deleteHandler удалить пост
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.NotFound(w, r)
	}

	delete(posts, id)

	http.Redirect(w, r, "/", 302)
}

func main() {

	posts = make(map[string]*models.Post, 0)

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                // Specify what path to load the templates from.
		Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		// Funcs:           []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		// Delims:          render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
		Charset:         "UTF-8",                 // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON:      true,                    // Output human readable JSON
		IndentXML:       true,                    // Output human readable XML
		HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html"
	}))

	// http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	staticOption := martini.StaticOptions{
		Prefix: "assets",
	}

	m.Use(martini.Static("assets", staticOption))

	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit", editHandler)
	m.Get("/delete", deleteHandler)
	m.Post("/SavePost", savePostHandler)

	// http.ListenAndServe(":8080", nil)
	m.RunOnAddr(":8080")
}
