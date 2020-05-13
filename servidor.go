package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
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

func (this *Calificaciones) Evaluar(datos []string, respuesta *string) error {
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
		} else {
			*respuesta = "Error, Evaluación ya existente"
		}
	} else {
		*respuesta = "Error, Evaluación debe ser numerica"
	}

	return nil
}

func (this *Calificaciones) Promedio(datos []string, respuesta *float64) error {
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
		} else {
			*respuesta = promedio / total
		}
	} else {
		*respuesta = 0.0
	}
	return nil
}

func server() {
	rpc.Register(new(Calificaciones))
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func form(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("form.html"),
		alumnosHTML(),
		materiasHTML(),
	)
}

func root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(res, cargarHtml("index.html"))
}

func cargarHtml(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func main() {
	http.HandleFunc("/", root)
	fmt.Println("Arrancando el servidor...")
	http.ListenAndServe(":9000", nil)
}
