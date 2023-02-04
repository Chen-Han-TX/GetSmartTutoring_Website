import React, { useState, useRef, useEffect } from 'react';
import { Container, Row, Col, Form, Button, ListGroup } from 'react-bootstrap';
import AuthService from "../services/auth.service";
import ChattingServices from '../services/chatting.service';

function ChatRoom({ chatDetail }) {

  const currentUser = AuthService.getCurrentUser();

  const [message, setMessage] = useState('');
  
  const [messages, setMessages] = useState(chatDetail.messages);
  const listGroupRef = useRef(null);

  useEffect(() => {
    listGroupRef.current.scrollTop = listGroupRef.current.scrollHeight;
  }, [messages]);

  const handleSubmit = (event) => {
    event.preventDefault();
    if (message !== "") {
          ChattingServices.updateChatList(chatDetail.chat_id, {
            "sender_id": currentUser.user_id,
            "content": message,
            "timestamp": Date.now()
        })
        var opp_user = ""
        if (currentUser.user_type === "Tutor") {
            opp_user = chatDetail.student_id
        } else {
            opp_user = chatDetail.tutor_id
        }
        ChattingServices.sendMsg(opp_user, message).then(
            (response) => {
              console.log(response)
              if (response.status === 200) {
                if (messages === null) {
                  setMessages([{
                      "sender_id": currentUser.user_id,
                      "content": message,
                      "timestamp": Date.now()
                    }])
                } else {
                  setMessages([...messages, {
                    "sender_id": currentUser.user_id,
                    "content": message,
                    "timestamp": Date.now()
                }]);
                }
                setMessage('');
              } else {
                console.log("response status: " + response.status);
                alert("Server problem")
              }
            },
            (error) => {
              if (error.response.status == 404){
                alert("something went wrong")
              }
            }
        );
    }
  };
  return (
    <div className='auth-inner'>
    {currentUser.user_type === "Student" && (
      <strong>Chat with: {chatDetail.tutor_name}</strong>
    )}
    {currentUser.user_type === "Tutor" && (
          <strong>Chat with: {chatDetail.student_name}</strong>
    )}
    <Container>
      <Row className="">
        <Col  className="mx-auto">
          <div>
            <ListGroup ref={listGroupRef} className="mb-3" style={{height: '50vh', maxHeight: '50vh', overflowY: 'auto' }}>
             { messages && (
                messages.map((msg, index) => (
                  <ListGroup.Item 
                    key={index} 
                    style={{ 
                      wordWrap: 'break-word',
                      backgroundColor: msg.sender_id === currentUser.user_id ? '#F2F2F2' : '#ADD8E6',
                      textAlign:  msg.sender_id === currentUser.user_id ? 'right' : 'left'
                    }}
                  >
                    {msg.content}
                  </ListGroup.Item>
                ))
              )}

            </ListGroup>
            <Form onSubmit={handleSubmit}>
              <Form.Group className="d-flex align-items-center">
                <Form.Control
                  type="text"
                  placeholder="Enter your message"
                  value={message}
                  onChange={(event) => setMessage(event.target.value)}
                  style={{ backgroundColor: '#F2F2F2', flexGrow: 1 }}
                />
                <Button type="submit" style={{ backgroundColor: '#7DCEA0', border: 'none' }}>
                  Send
                </Button>
              </Form.Group>
            </Form>
          </div>
        </Col>
      </Row>
    </Container>

    </div>
  );
}

export default ChatRoom;
