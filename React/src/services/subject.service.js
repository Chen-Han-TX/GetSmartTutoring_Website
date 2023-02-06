import axios from "axios";

//const SUB_URL = "https://subject-4dcnj7fm6a-uc.a.run.app/api/getsubjects/"
//const SUB_URL = "http://localhost:5051/api/getsubjects/"
const SUB_URL = "http://34.134.136.241:5051/api/getsubjects/"

axios.defaults.withCredentials = true
let axiosConfig = {
    headers: {
        'Content-Type': 'text/plain',
       // "Access-Control-Allow-Origin": "http://localhost:3000"
    },
    withCredentials : true,
}

const allSubjects = () => {
    return axios.get(SUB_URL + "all" , axiosConfig)
    .then((response) => {
        localStorage.setItem("subjects", JSON.stringify(response.data));
        return response.data;
    });
}


const getAllSubjects = () => {
    return JSON.parse(localStorage.getItem("subjects"));
};
  

const SubjectServices = {
    allSubjects,
    getAllSubjects
}

  
export default SubjectServices;



