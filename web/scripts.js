function addExam() {
    const axios = window.axios;
    const formData = {
        "function": "createExam",
        "exam": {
            "name": document.getElementById("name").value,
            "start": parseInt(document.getElementById("start").value),
            "end": parseInt(document.getElementById("end").value),
            "students": parseInt(document.getElementById("students").value)
        }
    }
    const url = "https://w65yfftrcg.execute-api.us-east-1.amazonaws.com/api/exams"
    axios.post(url, formData);
}