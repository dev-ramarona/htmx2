package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Struktur data dummy
type Employee struct {
	ID       int
	Name     string
	Gender   string
	Age      int
	Position string
	Email    string
	Phone    string
	Address  string
	City     string
	Status   string
}

var employees = []Employee{
	{1, "Budi Santoso", "Male", 28, "Software Engineer", "budi@example.com", "081234567801", "Jl. Melati No. 12", "Jakarta", "Active"},
	{2, "Siti Aminah", "Female", 31, "Product Manager", "siti@example.com", "081234567802", "Jl. Mawar No. 8", "Bandung", "Active"},
	{3, "Agus Prasetyo", "Male", 26, "UI/UX Designer", "agus@example.com", "081234567803", "Jl. Anggrek No. 15", "Surabaya", "Active"},
	{4, "Dewi Lestari", "Female", 29, "Data Analyst", "dewi@example.com", "081234567804", "Jl. Kenanga No. 21", "Yogyakarta", "Active"},
	{5, "Reza Rahadian", "Male", 32, "DevOps Engineer", "reza@example.com", "081234567805", "Jl. Flamboyan No. 5", "Bekasi", "Active"},
	{6, "Rina Kusuma", "Female", 27, "Backend Developer", "rina@example.com", "081234567806", "Jl. Sakura No. 10", "Depok", "Active"},
	{7, "Andi Wijaya", "Male", 30, "Frontend Developer", "andi@example.com", "081234567807", "Jl. Cemara No. 19", "Bogor", "Active"},
	{8, "Fajar Nugroho", "Male", 25, "QA Engineer", "fajar@example.com", "081234567808", "Jl. Merpati No. 7", "Semarang", "Active"},
	{9, "Nadia Putri", "Female", 28, "Business Analyst", "nadia@example.com", "081234567809", "Jl. Teratai No. 14", "Malang", "Active"},
	{10, "Yoga Pratama", "Male", 34, "System Analyst", "yoga@example.com", "081234567810", "Jl. Dahlia No. 18", "Tangerang", "Inactive"},
	{11, "Lina Marlina", "Female", 29, "HR Specialist", "lina@example.com", "081234567811", "Jl. Cendana No. 22", "Bandung", "Active"},
	{12, "Dimas Saputra", "Male", 27, "Mobile Developer", "dimas@example.com", "081234567812", "Jl. Nusa Indah No. 3", "Jakarta", "Active"},
	{13, "Maya Sari", "Female", 35, "Project Manager", "maya@example.com", "081234567813", "Jl. Kamboja No. 6", "Surabaya", "Active"},
	{14, "Hendra Gunawan", "Male", 33, "Cloud Engineer", "hendra@example.com", "081234567814", "Jl. Pahlawan No. 27", "Medan", "Active"},
	{15, "Putri Ayu", "Female", 24, "Technical Writer", "putri@example.com", "081234567815", "Jl. Kartini No. 9", "Solo", "Active"},
	{16, "Rizky Hidayat", "Male", 31, "Security Engineer", "rizky@example.com", "081234567816", "Jl. Diponegoro No. 11", "Makassar", "Active"},
	{17, "Nanda Permata", "Female", 30, "Database Administrator", "nanda@example.com", "081234567817", "Jl. Ahmad Yani No. 45", "Palembang", "Active"},
	{18, "Arif Setiawan", "Male", 36, "Scrum Master", "arif@example.com", "081234567818", "Jl. Sudirman No. 55", "Bandar Lampung", "Active"},
	{19, "Intan Maharani", "Female", 28, "Marketing Specialist", "intan@example.com", "081234567819", "Jl. Imam Bonjol No. 13", "Denpasar", "On Leave"},
	{20, "Farhan Akbar", "Male", 26, "AI Engineer", "farhan@example.com", "081234567820", "Jl. Gatot Subroto No. 17", "Balikpapan", "Active"},
}

// Middleware untuk mengecek apakah user sudah login
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Fungsi untuk mengubah struct menjadi slice of interface{}
func StructToSlice(v any) ([]string, []string) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	fields := make([]string, val.NumField())
	headers := make([]string, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		fields[i] = fmt.Sprintf("%v", val.Field(i).Interface())
		headers[i] = typ.Field(i).Name
	}
	return fields, headers
}

func main() {
	r := gin.Default()

	// Setup session store
	store := cookie.NewStore([]byte("secret_key_super_aman"))
	// TAMBAHKAN KONFIGURASI INI
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,                // Umur sesi dalam detik (1 hari)
		HttpOnly: true,                 // Mencegah akses cookie dari JavaScript (menghindari XSS)
		Secure:   false,                // PENTING: Harus 'false' agar bisa jalan di HTTP/IP biasa
		SameSite: http.SameSiteLaxMode, // Memastikan cookie dikirim di IP yang sama
	})
	r.Use(sessions.Sessions("mysession", store))

	// Load HTML templates
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Rute Publik
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/dashboard") })
	r.GET("/login", func(c *gin.Context) {
		c.Header("HX-Redirect", "/login")
		c.HTML(http.StatusOK, "login.html", nil)
		c.Header("HX-Refresh", "true")
	})

	// Proses Login
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Hardcoded login untuk tutorial (admin/admin)
		if username == "admin" && password == "admin" {
			session := sessions.Default(c)
			session.Set("user", username)
			session.Save()
			// Gunakan HX-Redirect untuk memberitahu HTMX agar redirect halaman penuh
			c.Header("HX-Redirect", "/dashboard")
			c.Status(http.StatusOK)
		} else {
			c.HTML(http.StatusOK, "login_form", gin.H{"Error": "Username atau Password salah"})
		}
	})

	// Proses Logout
	r.POST("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		session.Options(sessions.Options{MaxAge: -1})
		c.Redirect(http.StatusFound, "/login")
		c.Header("HX-Refresh", "true")
	})

	// Rute Privat (Membutuhkan Login)
	private := r.Group("/")
	private.Use(AuthRequired())
	{
		private.GET("/dashboard", func(c *gin.Context) {
			session := sessions.Default(c)
			user := session.Get("user")
			results := [][]string{}
			headers := []string{}
			for _, emp := range employees {
				strslc, header := StructToSlice(emp)
				results = append(results, strslc)
				if len(headers) == 0 {
					headers = header
				}
			}
			c.HTML(http.StatusOK, "dashboard.html", gin.H{"User": user, "Employees": results, "Headers": headers})
		})

		// Rute untuk HTMX Search
		private.GET("/search", func(c *gin.Context) {
			query := strings.ToLower(c.Query("q"))
			var results [][]string
			var headers []string
			for _, emp := range employees {
				if strings.Contains(strings.ToLower(emp.Name), query) ||
					strings.Contains(strings.ToLower(emp.Position), query) {
					strslc, header := StructToSlice(emp)
					if len(headers) == 0 {
						headers = header
					}
					results = append(results, strslc)
				}
			}
			// Hanya me-render bagian baris tabel (partial)
			c.HTML(http.StatusOK, "rows.html", gin.H{"Employees": results, "Headers": headers})
		})

		// 1. Rute untuk MENAMPILKAN modal edit
		private.GET("/edit/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var emp Employee

			// Cari data employee berdasarkan ID
			for _, e := range employees {
				if e.ID == id {
					emp = e
					break
				}
			}
			c.HTML(http.StatusOK, "edit_modal.html", emp)
		})

		// 2. Rute untuk MENYIMPAN perubahan data
		private.POST("/edit/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))

			// Update data di slice Go
			for i, e := range employees {
				if e.ID == id {
					employees[i].Name = c.PostForm("name")
					employees[i].Position = c.PostForm("position")
					employees[i].Email = c.PostForm("email")
					break
				}
			}

			// Render ulang baris tabel + hapus modal (OOB Swap)
			c.HTML(http.StatusOK, "update_success.html", gin.H{"Employees": employees})
		})
	}

	r.Run("0.0.0.0:8080")
}
