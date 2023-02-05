import AuthService from "../services/auth.service";
import React, { useState, useRef } from "react";


const Tutor = () => {
  const currentUser = AuthService.getCurrentUser();


  return (
    <div className="auth-inner">
      <header className="tutorPage">
      <h3>Welcome!</h3> 
      <h3> Tutor {currentUser.name}</h3>
      </header>
      
    </div>
  );
};

export default Tutor; 