package main

import (
	"net/http"
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
	Position string
	Email    string
}

var employees = []Employee{
	{1, "Budi Santoso", "Software Engineer", "budi@example.com"},
	{2, "Siti Aminah", "Product Manager", "siti@example.com"},
	{3, "Agus Prasetyo", "UI/UX Designer", "agus@example.com"},
	{4, "Dewi Lestari", "Data Analyst", "dewi@example.com"},
	{5, "Reza Rahadian", "DevOps Engineer", "reza@example.com"},
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
	r.LoadHTMLGlob("templates/*")

	// Rute Publik
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/dashboard") })
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })

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
			c.HTML(http.StatusOK, "login.html", gin.H{"Error": "Username atau Password salah"})
		}
	})

	// Proses Logout
	r.POST("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
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
			c.HTML(http.StatusOK, "dashboard.html", gin.H{"User": user, "Employees": employees})
		})

		// Rute untuk HTMX Search
		private.GET("/search", func(c *gin.Context) {
			query := strings.ToLower(c.Query("q"))
			var results []Employee

			for _, emp := range employees {
				if strings.Contains(strings.ToLower(emp.Name), query) ||
					strings.Contains(strings.ToLower(emp.Position), query) {
					results = append(results, emp)
				}
			}
			// Hanya me-render bagian baris tabel (partial)
			c.HTML(http.StatusOK, "rows.html", gin.H{"Employees": results})
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
