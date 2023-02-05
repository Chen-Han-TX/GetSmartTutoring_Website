import AuthService from "../services/auth.service";
import React, { useState, useRef } from "react";


const Student = () => {
  const currentUser = AuthService.getCurrentUser();


  return (
    <div className="auth-inner">
      <header className="studentPage">
      <h3>Welcome!</h3> 
      <h3> Student {currentUser.name}</h3>
      </header>
      
    </div>
  );
};

export default Student; 