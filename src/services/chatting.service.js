import { upload } from "@testing-library/user-event/dist/upload";
import axios from "axios";

const CHATTING_URL = "http://localhost:5053/api/"

axios.defaults.withCredentials = true

let axiosConfig = {
    headers: {
        'Content-Type': 'text/plain',
       // "Access-Control-Allow-Origin": "http://localhost:3000"
    },
    withCredentials : true,
}

const createChatList = () => {
    return axios.post(CHATTING_URL + "createchatlist",  {},
    axiosConfig).then(
        (response) => {
            return response;
        }
    )
}

 
const getChatList = (user_id, user_type) => {
    return axios.get(CHATTING_URL + "getlist/" + user_id + "/" + user_type, 
        axiosConfig)
            .then((response) => {
                localStorage.setItem("chatList", JSON.stringify(response.data));
        return response;
      });
};

const sendMsg = (user_id, opp_user_id, user_type, content) => {
    return axios.post(CHATTING_URL + "sendmessages/" + user_id + "/" + opp_user_id + "/" + user_type, { "content": content},
        axiosConfig)
            .then((response) => {
        return response;
      });
};



const getCurrentChatList = () => {
    return JSON.parse(localStorage.getItem("chatList"));
};

const updateChatList = (chatId, message) => {
    var chatList = JSON.parse(localStorage.getItem("chatList"))

    for (var i = 0; i < chatList.length; i++) {
        if (chatList[i].chat_id === chatId) {
            if (chatList[i].messages === null) {
                chatList[i].messages = [message]
            } else {
                chatList[i].messages.push(message)
            }
            localStorage.setItem("chatList", JSON.stringify(chatList))
            return chatList
        }
    }
}


const ChattingServices = {
    createChatList,
    sendMsg, 
    getChatList,
    getCurrentChatList,
    updateChatList
}

  
export default ChattingServices;

