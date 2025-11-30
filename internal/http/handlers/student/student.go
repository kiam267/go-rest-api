package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kiam267/student-api/internal/storage"
	"github.com/kiam267/student-api/internal/types"
	response "github.com/kiam267/student-api/internal/utils/reponse"
)




func New(storage storage.Storage ) http.HandlerFunc{
 return func(w http.ResponseWriter, r *http.Request)  {
   var student types.Student
	 err :=json.NewDecoder(r.Body).Decode(&student)

	 if errors.Is(err, io.EOF){
    response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")) )
		return
	 };

   if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err)) 
		return 
	 }
    //Reques Validation 

	 if err :=	validator.New().Struct(student); err != nil {
		validateError := err.(validator.ValidationErrors)
		response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateError))
		return
	 }

    
     
		lastId, err :=  storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		 )
		 
    slog.Info("User created successfully", slog.String("userId", fmt.Sprint(lastId)))

		 if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		 }


		response.WriteJson(w, http.StatusCreated, map[string] int64 {"id": lastId})
	}
}

