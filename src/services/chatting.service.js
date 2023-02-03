import { upload } from "@testing-library/user-event/dist/upload";
import axios from "axios";

const CHATTING_URL = "http://localhost:5070/api/"

axios.defaults.withCredentials = true

let axiosConfig = {
    headers: {
        'Content-Type': 'text/plain',
       // "Access-Control-Allow-Origin": "http://localhost:3000"
    },
    withCredentials : true,
}



const sendMsg = (opp_user_id, content) => {
    return axios.post(CHATTING_URL + "sendmessages/" + opp_user_id, { "content": content},
        axiosConfig)
            .then((response) => {
        return response;
      });
};


const getChatList = () => {
    return axios.get(CHATTING_URL + "getlist", 
        axiosConfig)
            .then((response) => {
        return response;
      });
};


const ChattingServices = {
    sendMsg, 
    getChatList
}

  
export default ChattingServices;

