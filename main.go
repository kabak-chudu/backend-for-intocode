package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	connectdatabase "github.com/kabak-chudu/backend-for-intocode/connect_database"
	"github.com/kabak-chudu/backend-for-intocode/groups"
	"github.com/kabak-chudu/backend-for-intocode/notes"
	"github.com/kabak-chudu/backend-for-intocode/students"
)

func main() {
	db, err := connectdatabase.Connect()
	if err != nil {
		panic(err)
	} else {
		if err := db.AutoMigrate(&groups.Group{}, &students.Student{}, &notes.Note{}); err != nil {
			log.Println(fmt.Errorf("не удалось совершить автомигрирацию %w", err))
		}
	}

	router := gin.Default()
	// students
	router.GET("/students", GetStudents)
	router.GET("/students/:id", GetStudentByID)
	router.GET("/groups/:id/students", GetStudentsByGroup)
	router.POST("/students", AddStudent)
	router.PATCH("/students/:id", PATCHStudent)
	router.DELETE("/students/:id", DELETEStudent)
	// groups
	router.GET("/groups/:id/stats/offer", GETOffer) // не сделал пока percent_offer
	router.GET("/groups", GetGroups)
	router.GET("/groups/:id", GetGroupByID)
	router.POST("/groups", AddGroup)
	router.PATCH("/groups/:id", PATCHGroup)
	router.DELETE("/groups/:id", DELETEGroup)
	//notes
	router.GET("/notes/:id", NoteByID)
	router.GET("/students/:id/notes", NotesForStudent) // — все заметки по студенту;
	router.POST("/notes", AddNote)                     // — создание заметки (в теле запроса должен быть student_id);
	router.PATCH("/notes/:id", PATCHNote)              //— редактирование текста;
	router.DELETE("/notes/:id", DELETENote)            //— удаление заметки.
	router.Run()
}

func GetOfferPercent(group_id uint) (float64, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return 0, err
	}
	var total int64
	var offers int64
	db.Model(&students.Student{}).Where("group_id", group_id).Count(&total)
	if total == 0 {
		return 0, nil
	}

	db.Model(&students.Student{}).Where("group_id = ? AND study_status = ?", group_id, "offer").Count(&offers)

	percent := (float64(offers) / float64(total)) * 100
	return percent, nil

}

func GETOffer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	percent, err := GetOfferPercent(uint(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"group_id":      id,
		"offer_percent": percent,
	})
}

func NotesForStudent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notes, err := notes.GetAllNotesForStudent(uint(id))
	if err != nil {
		ctx.JSON(http.StatusOK, notes)
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"notes for student": notes})
}

func AddNote(ctx *gin.Context) {
	var req struct {
		Student_ID *uint  `json:"student_id"`
		Author     string `json:"author"`
		Text       string `json:"text" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Student_ID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "поле student_id обязательно надо указать"})
		return
	}
	if *req.Student_ID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "поле student_id должнео быть больше 0"})
		return
	}

	note, err := notes.CreateNote(*req.Student_ID, req.Author, req.Text)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"created": note})

}

func PATCHNote(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		Text string `json:"text" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := notes.UpdateNote(uint(id), req.Text)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"edited": note})
}

func DELETENote(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := notes.DeleteNote(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func NoteByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := notes.GetNoteID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, note)
}

func GetStudentsByGroup(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	students, err := students.GetStudentsByGroupID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusOK, students)
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"students": students})

}

func DELETEStudent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := students.DeleteStudent(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func AddStudent(ctx *gin.Context) {
	var req struct {
		Full_name     string `json:"full_name" binding:"required"`
		Email         string `json:"email"`
		Telegram      string `json:"telegram"`
		Group_id      *uint  `json:"group_id"`
		Tuition_total *uint  `json:"tuition_total"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Group_id == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "введите группу для студента"})
		return
	}
	if req.Tuition_total == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "введите общую оплату для студента"})
		return
	}
	if *req.Tuition_total < 200000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "минимальная общая оплата 200000!"})
		return
	}

	student, err := students.CreateStudent(req.Full_name, req.Email, req.Telegram, *req.Group_id, *req.Tuition_total)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"created": student})
}

func PATCHStudent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		Tuition_paid *uint   `json:"tuition_paid"`
		Email        *string `json:"email"`
		Telegram     *string `json:"telegram"`
		Study_status *string `json:"study_status"`
		Full_name    *string `json:"full_name"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, err := students.UpdateStudent(uint(id), req.Tuition_paid, req.Email, req.Telegram, req.Study_status, req.Full_name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"edited": student})
}

func DELETEGroup(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := groups.DeleteGroup(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func AddGroup(ctx *gin.Context) {
	var req struct {
		Title       string `json:"title"`
		Total_weeks *uint  `json:"total_weeks"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Title == "" || len(req.Title) < 5 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "title пустой или короткий введите минимум 5 символов"})
		return
	}
	if req.Total_weeks != nil {
		if *req.Total_weeks <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "total_weeks не может быть отрицателен или равен нулю"})
			return
		}
		if *req.Total_weeks == 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "total_weeks должен быть хотя бы больше единицы"})
			return
		}
		if *req.Total_weeks > 15 {
			*req.Total_weeks = 15
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "обязательно введите total_weeks!"})
		return
	}

	group, err := groups.CreateGroup(req.Title, *req.Total_weeks)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"created": group})
}

func PATCHGroup(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Current_week uint   `json:"current_week"`
		Title        string `json:"title"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	group, err := groups.UpdateGroup(uint(id), req.Current_week, req.Title)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"edited": group})

}

func GetGroupByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	group, err := groups.GetGroupID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusOK, group)
		return
	}

	ctx.JSON(http.StatusOK, group)
}

func GetGroups(ctx *gin.Context) {
	week_s := ctx.Query("week")
	is_finished_s := ctx.Query("finished")
	if week_s != "" {
		week, err := strconv.Atoi(week_s)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		groups, err := groups.GetGroupsByWeek(uint(week))
		if err != nil {
			ctx.JSON(http.StatusOK, groups)
			return
		}
		ctx.IndentedJSON(http.StatusOK, gin.H{"groups": groups})
		return
	}
	if is_finished_s != "" {
		isFinished, err := strconv.ParseBool(is_finished_s)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "используйте true или false для этого фильтра"})
			return
		}
		groups, err := groups.GetGroupsFinished(isFinished)
		if err != nil {
			ctx.JSON(http.StatusOK, groups)
			return
		}

		ctx.IndentedJSON(http.StatusOK, gin.H{"groups": groups})
		return
	}
	groups, err := groups.GetAllGroups()
	if err != nil {
		ctx.JSON(http.StatusOK, groups)
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"groups": groups})
}

func GetStudentByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if id < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "айди < 0!!!"})
		return
	}
	student, err := students.GetStudentByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, student)
}

func GetStudents(ctx *gin.Context) {
	group_ids := ctx.Query("group_id")
	payment_status := ctx.Query("payment_status")
	study_status := ctx.Query("study_status")

	studentsDB, err := students.GetStudentsFiltered(group_ids, payment_status, study_status)
	if err != nil {
		ctx.JSON(http.StatusOK, studentsDB)
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"students": studentsDB})

	// if group_ids != "" {
	// 	group_id, err := strconv.Atoi(group_ids)
	// 	if err != nil {
	// 		ctx.JSON(400, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	student, err := students.GetStudentsByGroupID(uint(group_id))
	// 	if err != nil {
	// 		ctx.JSON(200, student)
	// 		return
	// 	}

	// 	ctx.JSON(200, student)
	// 	return
	// }
	// if payment_status != "" {
	// 	student, err := students.GetStudentsByPaymentStatus(payment_status)
	// 	if err != nil {
	// 		ctx.JSON(404, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	ctx.JSON(200, student)
	// 	return
	// }
	// if study_status != "" {
	// 	student, err := students.GetStudentsByStudyStatus(study_status)
	// 	if err != nil {
	// 		ctx.JSON(404, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	ctx.JSON(200, student)
	// 	return
	// }

	// studentsDB, err := students.GetAllStudents()
	// if err != nil {
	// 	ctx.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	// ctx.IndentedJSON(200, gin.H{"students": studentsDB})
}
