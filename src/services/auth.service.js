// This ride.service.js containing all the API request to the microservices - Auth.go required in the project.
import { upload } from "@testing-library/user-event/dist/upload";
import axios from "axios";

const AUTH_URL = "http://localhost:5050/api/auth/"

axios.defaults.withCredentials = true
let axiosConfig = {
    headers: {
        'Content-Type': 'text/plain',
       // "Access-Control-Allow-Origin": "http://localhost:3000"
    },
    withCredentials : true,
}
const register_student = (name, email, password, school, subjects) => {
    return axios.post(AUTH_URL + "signup/student", {
        "name" : name,
        "email": email,
        "password" : password,
        "school" : school,
        "area_of_interest": subjects
    }, axiosConfig);
};


const register_tutor = (name, email, password, hourlyRate, availability, subjects, uploadedFiles) => {
    return axios.post(AUTH_URL + "signup/tutor", {
        "name" : name,
        "email": email,
        "password" : password,
        "hourly_rate": hourlyRate,
        "availability": availability,
        "area_of_interest": subjects,
        "cert_of_evidence": uploadedFiles
    }, axiosConfig);
};



const login = (email_address, password) => {
    return axios.post(AUTH_URL + "login", {
        "email": email_address,
        "password": password },
        axiosConfig)
            .then((response) => {
                console.log(response.data)
            if (response.data.email) {

                delete response.data.password
                localStorage.setItem("user", JSON.stringify(response.data));
            }
        return response.data;
      });
      
};
  
const logout = () => {
    localStorage.clear();
};


const getCurrentUser = () => {
    return JSON.parse(localStorage.getItem("user"));
};
  

const AuthService = {
    register_student,
    register_tutor,
    login,
    logout,
    getCurrentUser,
}

  
export default AuthService;



