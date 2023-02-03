import React, { useState } from 'react';
import { Dropdown, Container, Row, Col } from 'react-bootstrap';
import Chatroom from './Chatroom';
import AuthService from "../services/auth.service";

function ChatRoomPage() {

  const currentUser = AuthService.getCurrentUser();

  const chatrooms = [
    {
    "chat_id":"d4UJ9CuF7qbHHaGBgAq5",
    "student_id":"OpoyFtGk74TZ0ZQgK2tkYTlBCJ33",
    "tutor_id":"vXPiFWy1EoYlUUgoeEQlNSnirFB3",
    "student_name":"Chen Han",
    "tutor_name":"Ben Low",
    "messages":[
        {"sender_id":"OpoyFtGk74TZ0ZQgK2tkYTlBCJ33","content":"hello", "timestamp":1675398394},
        {"sender_id":"OpoyFtGk74TZ0ZQgK2tkYTlBCJ33","content":"hello no2","timestamp":1675398412},
        {"sender_id":"OpoyFtGk74TZ0ZQgK2tkYTlBCJ33","content":"hello no3","timestamp":1675399198},
        {"sender_id":"vXPiFWy1EoYlUUgoeEQlNSnirFB3","content":"bruh","timestamp":1675408622}]
    }
]
  const [activeChatroom, setActiveChatroom] = useState(chatrooms[0]);

/*
  return (
    <Container fluid>
      <Row>
        <Col md={3}>
          <Nav className="flex-column">
            {chatrooms.map((chatroom) => (
              <Nav.Link
                key={chatroom}
                onClick={() => setActiveChatroom(chatroom.chat_id)}
              >
                {chatroom.chat_id}
              </Nav.Link>
            ))}
          </Nav>
        </Col>
        <Col md={9}>
          {chatrooms.map((chatroom) => {
            if (chatroom.chat_id === activeChatroom) {
              return <Chatroom key={activeChatroom} chatDetail={chatroom} />;
            }
            return null;
          })}
        </Col>
      </Row>
    </Container>
  );
}
*/

return (
 <div className='auth-inner2' style={{'text-align': 'center'}}>
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
                    onClick={() => setActiveChatroom(chatroom)}
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
    </div>
  );
}

export default ChatRoomPage;
