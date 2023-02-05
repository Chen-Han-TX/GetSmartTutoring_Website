import React, { useState, useEffect } from 'react';
import { Dropdown, Container, Row, Col } from 'react-bootstrap';
import Chatroom from './Chatroom';
import AuthService from "../services/auth.service";
import ChattingServices from '../services/chatting.service';

function ChatRoomPage() {

  const currentUser = AuthService.getCurrentUser();
  const [chatrooms, setChatrooms] = useState(ChattingServices.getCurrentChatList());

  const [activeChatroom, setActiveChatroom] = useState(
    chatrooms && chatrooms.length > 0 ? chatrooms[0] : null
  );
  
  const changeRoom = (chatroom) => {
    ChattingServices.getChatList(currentUser.user_id, currentUser.user_type).then(
      (response) => {
        setChatrooms(response.data)
        if (chatroom.messages === null){
          chatroom.messages = []
        }
        setActiveChatroom(chatroom)
      },
      (error) => {
        if (error.response.status == 404){
          console.log(error)
          alert(error)
        }
      }
    );
  }
  /*
  useEffect(() => {
    ChattingServices.getChatList().then(
        (response) => {
          setChatRooms(response.data)
          setActiveChatroom(chatrooms[0])
        },
        (error) => {
          if (error.response.status == 404){
            console.log(error)
            alert(error)
          }
        }
      );

  }, [chatrooms]);
  */
return (
 <div className='auth-inner2' style={{textAlign: "center"}}>

    {activeChatroom ? (

    <Container fluid>
    <Row>
        <Col md={12}>
            <br></br>
        <Dropdown>
            <Dropdown.Toggle variant="light" size='lg' id="dropdown-basic">
                {currentUser.user_type === "Tutor" && (
                    "Student: " + activeChatroom.student_name
                )}
                {currentUser.user_type === "Student" && (
                    "Tutor: " + activeChatroom.tutor_name
                )}
            </Dropdown.Toggle>

            <Dropdown.Menu>
            {chatrooms.map((chatroom) => (
                <Dropdown.Item
                key={chatroom.chat_id}
                onClick={() => changeRoom(chatroom)}
                >
                {currentUser.user_type === "Tutor" && (
                    chatroom.student_name
                )}
                {currentUser.user_type === "Student" && (
                    chatroom.tutor_name
                )}
                </Dropdown.Item>
            ))}
            </Dropdown.Menu>
        </Dropdown>
        </Col>
        <Col md={12}>
        {chatrooms.map((chatroom) => {
            if (chatroom.chat_id === activeChatroom.chat_id) {
            return <Chatroom key={chatroom.chat_id} chatDetail={chatroom} />;
            }
            return null;
        })}
        </Col>
    </Row>
    </Container>


    ) : (
      <div className='auth-inner'>
        <h3>You do not have a chatroom</h3>
        {currentUser.user_type == "student" && (
          <div>please book a session with your desired tutor first!
            </div>
        )} 
        </div>
    )}
    </div>
  );
}

export default ChatRoomPage;
