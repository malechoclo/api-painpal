package models

import "time"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}



// USUARIO
type Direccion struct {
	Calle  string `json:"calle" firestore:"calle"`
	Numero int    `json:"numero" firestore:"numero"`
	Ciudad string `json:"ciudad" firestore:"ciudad"`
	Estado string `json:"estado" firestore:"estado"`
}

type Enfermedad struct {
	Nombre           string `json:"nombre" firestore:"nombre"`
	FechaDiagnostico string `json:"fecha_diagnostico" firestore:"fecha_diagnostico"`
	Tratamiento      string `json:"tratamiento" firestore:"tratamiento"`
}

type Alergia struct {
	Nombre string `json:"nombre" firestore:"nombre"`
	Tipo   string `json:"tipo" firestore:"tipo"`
}

type HistorialMedico struct {
	Fecha      string `json:"fecha" firestore:"fecha"`
	Descripcion string `json:"descripcion" firestore:"descripcion"`
	Doctor     string `json:"doctor" firestore:"doctor"`
}

type ContactoEmergencia struct {
	Nombre   string `json:"nombre" firestore:"nombre"`
	Relacion string `json:"relacion" firestore:"relacion"`
	Telefono string `json:"telefono" firestore:"telefono"`
}

type Usuario struct {
	UID               string              `json:"UID" firestore:"UID"`
	Nombre            string              `json:"nombre" firestore:"nombre"`
	Apellido          string              `json:"apellido" firestore:"apellido"`
	Edad              int                 `json:"edad" firestore:"edad"`
	Email             string              `json:"email" firestore:"email"`
	Telefono          string              `json:"telefono" firestore:"telefono"`
	Tipo              string              `json:"tipo" firestore:"tipo"`
	Activo            bool                `json:"activo" firestore:"activo"`
	Eliminado         bool                `json:"eliminado" firestore:"eliminado"`
	CreatedAt         time.Time           `json:"createdAt" firestore:"createdAt"`
	Direccion         Direccion           `json:"direccion" firestore:"direccion"`
	Enfermedades      []Enfermedad        `json:"enfermedades" firestore:"enfermedades"`
	Alergias          []Alergia           `json:"alergias" firestore:"alergias"`
	HistorialMedico   []HistorialMedico   `json:"historial_medico" firestore:"historial_medico"`
	ContactoEmergencia ContactoEmergencia `json:"contacto_emergencia" firestore:"contacto_emergencia"`
}


//FORM
type Respuesta struct {
	Pregunta string `json:"pregunta" firestore:"pregunta"`
	Respuesta string `json:"respuesta" firestore:"respuesta"`
}

type Formulario struct {
	Nombre     string      `json:"nombre" firestore:"nombre"`
	UID        string      `json:"UID" firestore:"UID"`
	Fecha      string      `json:"fecha" firestore:"fecha"` // puedes usar time.Time si prefieres
	Form       string      `json:"form" firestore:"form"`
	Respuestas []Respuesta `json:"respuestas" firestore:"respuestas"`
}