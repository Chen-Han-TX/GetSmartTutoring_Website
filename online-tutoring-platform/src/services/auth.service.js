// This ride.service.js containing all the API request to the microservices - Auth.go required in the project.
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

const register_tutor = (email_address, password, first_name, last_name, mobile_number) => {
    return axios.post(AUTH_URL + "signup/tutor", {
      "email_address" : email_address,
      "password": password,
      "first_name" : first_name,
      "last_name" : last_name,
      "mobile_number": mobile_number
    }, axiosConfig)
};

const register_student = (email_address, password, first_name, last_name, mobile_number, ic_number, car_lic_number) => {
    return axios.post(AUTH_URL + "signup/student", {
        "email_address" : email_address,
        "password": password,
        "first_name" : first_name,
        "last_name" : last_name,
        "mobile_number": mobile_number,
        "ic_number": ic_number,
        "car_lic_number": car_lic_number
    }, axiosConfig);
};


const login = (email_address, password) => {
    return axios.post(AUTH_URL + "login", {
        "email": email_address,
        "password": password },
        axiosConfig)
            .then((response) => {
                console.log(response.data)
            if (response.data.email_address) {

                delete response.data.password
                localStorage.setItem("user", JSON.stringify(response.data));
            }
        return response.data;
      });
      
};
  
const logout = () => {
    localStorage.clear();

    return axios.get(AUTH_URL + "logout", axiosConfig)
    .then((response) => {
        return response.data;
    });
};
  

const getCurrentUser = () => {
    return JSON.parse(localStorage.getItem("user"));
};
  

const AuthService = {
    register_passenger,
    register_rider,
    login,
    logout,
    getCurrentUser,
}

  
export default AuthService;



