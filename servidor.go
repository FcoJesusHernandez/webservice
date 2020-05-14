package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Materia struct {
	Nombre string
}

type Alumno struct {
	Nombre string
}

type Calificacion struct {
	Alumno       Alumno
	Materia      Materia
	Calificacion float64
}

type Calificaciones struct {
	Calificaciones list.List
}

var lista_calificaciones list.List
var lista_alumnos list.List
var lista_materias list.List

func (this *Calificaciones) Evaluar(datos []string, respuesta *string, danger *bool) error {
	alumno_ := Alumno{
		Nombre: datos[0],
	}

	materia_ := Materia{
		Nombre: datos[1],
	}

	num, err := strconv.ParseFloat(datos[2], 64)
	if err == nil {
		evaluacion := Calificacion{
			Alumno:       alumno_,
			Materia:      materia_,
			Calificacion: num,
		}

		var bandera = false
		for e := lista_alumnos.Front(); e != nil; e = e.Next() {
			if e.Value.(Alumno).Nombre == alumno_.Nombre {
				bandera = true
			}
		}

		if !bandera {
			lista_alumnos.PushBack(alumno_)
		}

		bandera = false
		for e := lista_materias.Front(); e != nil; e = e.Next() {
			if e.Value.(Materia).Nombre == materia_.Nombre {
				bandera = true
			}
		}

		if !bandera {
			lista_materias.PushBack(materia_)
		}

		bandera = false
		for e := lista_calificaciones.Front(); e != nil; e = e.Next() {
			if e.Value.(Calificacion).Alumno.Nombre == evaluacion.Alumno.Nombre && e.Value.(Calificacion).Materia.Nombre == evaluacion.Materia.Nombre {
				bandera = true
			}
		}

		if !bandera {
			lista_calificaciones.PushBack(evaluacion)
			*respuesta = "Evaluación anexada con éxito"
			*danger = false
		} else {
			*respuesta = "Error, Evaluación ya existente"
			*danger = true
		}
	} else {
		*respuesta = "Error, Evaluación debe ser numerica"
		*danger = true
	}

	return nil
}

func (this *Calificaciones) Promedio(datos []string, respuesta *float64, danger *bool) error {
	var total float64
	var promedio float64

	tipo := datos[0]
	auxiliar := datos[1]

	if tipo == "1" { // promedio de alumno
		alumno_ := Alumno{
			Nombre: auxiliar,
		}

		for e := lista_calificaciones.Front(); e != nil; e = e.Next() {
			if e.Value.(Calificacion).Alumno == alumno_ {
				total += 1
				promedio += e.Value.(Calificacion).Calificacion
			}
		}
		if total == 0 {
			*respuesta = 0
		} else {
			*respuesta = promedio / total
		}
	} else if tipo == "2" { // promedio general / todos
		for e := lista_calificaciones.Front(); e != nil; e = e.Next() {
			total += 1
			promedio += e.Value.(Calificacion).Calificacion
		}
		if total == 0 {
			*respuesta = 0
		} else {
			*respuesta = promedio / total
		}
	} else if tipo == "3" { // promedio de materia
		materia_ := Materia{
			Nombre: auxiliar,
		}

		for e := lista_calificaciones.Front(); e != nil; e = e.Next() {
			if e.Value.(Calificacion).Materia == materia_ {
				total += 1
				promedio += e.Value.(Calificacion).Calificacion
			}
		}
		if total == 0 {
			*respuesta = 0
			*danger = true
		} else {
			*respuesta = promedio / total
			*danger = false
		}
	} else {
		*respuesta = 0.0
		*danger = true
	}
	return nil
}

func form(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	mensaje := "Calificación exitosa"
	fmt.Fprintf(
		res,
		cargarHtml("form.html", "Agregar", mensaje, "", false, false),
	)
}

func root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(res, cargarHtml("index.html", "inicio", "Hola Bienvenido", "", false, false))
}

var clf = new(Calificaciones)

func calificacion(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		datos := []string{req.FormValue("alumno"), req.FormValue("materia"), req.FormValue("calificacion")}
		var result string
		var danger bool

		clf.Evaluar(datos, &result, &danger)

		res.Header().Set(
			"Content-Type",
			"text/html",
		)

		fmt.Fprintf(
			res,
			cargarHtml("index.html", "inicio", "Hola, Bienvenido", result, danger, false),
		)
	}
}

func promedio(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		datos := []string{}
		var salida string
		var result float64
		var danger bool

		if req.FormValue("alumno") != "" {
			salida = "El promedio de " + req.FormValue("alumno") + " es : "
			datos = []string{"1", req.FormValue("alumno")}
		} else if req.FormValue("materia") != "" {
			salida = "El promedio de " + req.FormValue("materia") + " es : "
			datos = []string{"3", req.FormValue("materia")}
		} else {
			salida = "Petición desconocida"
			datos = []string{"4", "error"}
			result = 0.0
			danger = true
		}

		clf.Promedio(datos, &result, &danger)

		res.Header().Set(
			"Content-Type",
			"text/html",
		)

		fmt.Fprintf(
			res,
			cargarHtml("index.html", "promedio", "Hola, Bienvenido", salida+fmt.Sprintf("%f", result), danger, false),
		)
	}
}

func promedio_gen(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		datos := []string{}
		var salida string
		var result float64
		var danger bool

		salida = "El promedio general es : "
		datos = []string{"2", ""}

		clf.Promedio(datos, &result, &danger)

		res.Header().Set(
			"Content-Type",
			"text/html",
		)

		fmt.Fprintf(
			res,
			cargarHtml("index.html", "promedio general", "Hola, Bienvenido", salida+fmt.Sprintf("%f", result), danger, false),
		)
	}
}

func cargarHtml(a string, titulo string, mensaje string, auxiliar string, danger_aux bool, danger_msj bool) string {
	html, _ := ioutil.ReadFile(a)
	salida := strings.Replace(string(html), "$__TITULO__$", titulo, -1)
	salida = strings.Replace(salida, "$__MENSAJE__$", mensaje, -1)
	salida = strings.Replace(salida, "$__AUXILIAR__$", auxiliar, -1)

	if danger_aux {
		salida = strings.Replace(salida, "$__CLASS_AUX__$", "alert-danger", -1)
	} else {
		salida = strings.Replace(salida, "$__CLASS_AUX__$", "alert-success", -1)
	}

	if danger_msj {
		salida = strings.Replace(salida, "$__CLASS_MSJ__$", "alert-danger", -1)
	} else {
		salida = strings.Replace(salida, "$__CLASS_MSJ__$", "alert-secondary", -1)
	}

	//alumnosHTML(),
	//materiasHTML(),
	return salida
}

func main() {
	http.HandleFunc("/calificacion", calificacion)
	http.HandleFunc("/promedio", promedio)
	http.HandleFunc("/general", promedio_gen)
	http.HandleFunc("/inicio", root)
	fmt.Println("Arrancando el servidor...")
	http.ListenAndServe(":9000", nil)
}
