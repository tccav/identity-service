package hsfixtures

// POST /v1/identities/students fixtures
const (
	ValidStudentRequestBody = `{
    "id": "123451271",
    "name": "John",
    "cpf": "11111111030",
    "birth_date": "1994-03-19",
    "email": "jdoe@ol.com",
    "secret": "123456",
    "course_id": "6579705e-7e40-4b12-8ca1-7774ec3d6c3f"
}`
	InvalidJSONRequestBody = "invalid_json"
	InvalidIDRequestBody   = `{
    "id": "1as",
    "name": "John",
    "cpf": "11111111030",
    "birth_date": "1994-03-19",
    "email": "jdoe@ol.com",
    "secret": "123456",
    "course_id": "6579705e-7e40-4b12-8ca1-7774ec3d6c3f"
}`
	InvalidCPFRequestBody = `{
    "id": "123451271",
    "name": "John",
    "cpf": "111.111.110-30",
    "birth_date": "1994-03-19",
    "email": "jdoe@ol.com",
    "secret": "123456",
    "course_id": "6579705e-7e40-4b12-8ca1-7774ec3d6c3f"
}`
	InvalidBirthDateRequestBody = `{
    "id": "123451271",
    "name": "John",
    "cpf": "11111111030",
    "birth_date": "19940319",
    "email": "jdoe@ol.com",
    "secret": "123456",
    "course_id": "6579705e-7e40-4b12-8ca1-7774ec3d6c3f"
}`
	InvalidEmailRequestBody = `{
    "id": "123451271",
    "name": "John",
    "cpf": "11111111030",
    "birth_date": "1994-03-19",
    "email": "ol.com",
    "secret": "123456",
    "course_id": "6579705e-7e40-4b12-8ca1-7774ec3d6c3f"
}`
	InvalidCourseIDRequestBody = `{
    "id": "123451271",
    "name": "John",
    "cpf": "11111111030",
    "birth_date": "1994-03-19",
    "email": "jdoe@ol.com",
    "secret": "123456",
    "course_id": "657970"
}`
)

// POST /v1/identities/students/login fixtures
const (
	ValidStudentLoginRequestBody = `{
    "student_id": "123451271",
    "secret": "123456"
}`
)

// POST /v1/identities/students/verify-auth fixtures
const (
	ValidAuthHeader = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)
