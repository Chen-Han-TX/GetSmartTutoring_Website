import { upload } from "@testing-library/user-event/dist/upload";
import axios from "axios";

const TUTORING_URL = "http://localhost:5054/api/tutoring/"

axios.defaults.withCredentials = true

let axiosConfig = {
    headers: {
        'Content-Type': 'text/plain',
       // "Access-Control-Allow-Origin": "http://localhost:3000"
    },
    withCredentials : true,
}

const matchTutors = (searchedSubjects) => {
    return axios.post(TUTORING_URL + "matchtutors", searchedSubjects,
        axiosConfig)
            .then((response) => {
                console.log(response.data)
        return response.data;
      });
};

const applyForTutor = (application) => {
    return axios.post(TUTORING_URL + "apply", application,
        axiosConfig)
            .then((response) => {
                console.log(response.data)
        return response.data;
      });
};

const getApplications = () => {
    return axios.get(TUTORING_URL + "getapplications", 
        axiosConfig)
            .then((response) => {
                console.log(response.data)
        return response.data;
      });
};

const handleApplications = (application) => {
    return axios.put(TUTORING_URL + "handleapplications", application,
        axiosConfig)
            .then((response) => {
                console.log(response.data)
        return response.data;
      });
};

const TutoringService = {
    matchTutors,
    applyForTutor,
    getApplications,
    handleApplications
}

  
export default TutoringService;

