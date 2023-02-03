import React, { useState, useRef, useEffect } from 'react';
import { Container, Row, Col, Form, Button, ListGroup } from 'react-bootstrap';
import AuthService from "../services/auth.service";

function ChatRoom({ chatDetail }) {

  const currentUser = AuthService.getCurrentUser();

  const [message, setMessage] = useState('');
  
  const [messages, setMessages] = useState(chatDetail.messages);

  const handleSubmit = (event) => {
    event.preventDefault();
    if (message !== "") {
        setMessages([...messages, {
            "sender_id": currentUser.user_id,
            "content": message,
            "timestamp": Date.now()
        }]);
        setMessage('');
    }
  };

  return (
    <div className='auth-inner'>
    <Container>
      <Row className="">
        <Col  className="mx-auto">
          <div>
            <ListGroup className="mb-3" style={{height: '50vh', maxHeight: '50vh', overflowY: 'auto' }}>
              {messages.map((msg, index) => (
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
              ))}
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