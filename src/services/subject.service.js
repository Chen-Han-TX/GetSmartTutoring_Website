import axios from "axios";

const SUB_URL = "http://localhost:5051/api/getsubjects/"



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



