//import logo from './logo.svg';
import './App.css';
import { useState, useEffect } from 'react';
import { Routes, Route, Link } from "react-router-dom";
import Login from "./components/Login";
import Tutor from "./components/Tutor/Tutor";
import Student from "./components/Student/Student";
import '../node_modules/bootstrap/dist/css/bootstrap.min.css'


import EventBus from "./common/EventBus";
import AuthService from './services/auth.service';


const App = () => {
  const [currentUser, setCurrentUser] = useState(undefined);
  const [showTutor, setShowTutor] = useState(false);
  const [showStudent, setShowStudent] = useState(false);

  useEffect(() => {
    const user = AuthService.getCurrentUser();

    
    if (user) {
      setCurrentUser(user);

      if (user.user_type === "student") {
        setShowStudent(true)
        setShowTutor(false)
      } else if (user.user_type === "tutor") {
        setShowTutor(true)
        setShowStudent(false)  
      }
    }
    EventBus.on("logout", () => {
      logOut();
    });
    return () => {
      EventBus.remove("logout");
    };
  }, []);

  const logOut = () => {
    AuthService.logout();
    setCurrentUser(undefined);
    setShowStudent(false)
    setShowTutor(false)
    localStorage.removeItem('user');
  };

  return (
    <div>
      {
        !currentUser ? (
          <nav className="navbar navbar-expand-lg navbar-light">
            <div className='container'>
              <Link to={"/"} className="navbar-brand">
                Get Smart Tutoring (GST)
              </Link>
              <div className='collapse navbar-collapse'>
                <ul className="navbar-nav ml-auto" >
                  <li className='nav-item'>
                    <Link to={"/login"} className="nav-link">
                      Login
                    </Link>
                  </li>
                  <li className='nav-item'>
                    <Link to={"/register/tutor"} className="nav-link">
                      Register Tutor
                    </Link>
                  </li>
                  <li className='nav-item'>
                    <Link to={"/register/student"} className="nav-link">
                      Register Student
                    </Link>
                  </li>        
                </ul>
              </div>
            </div>
          </nav>

        ) : (
            showStudent === true && showTutor === false ? (
              <nav className="navbar navbar-expand-lg navbar-light">
              <div className='container'>
                <Link to={"/"} className="navbar-brand">
                  GoRidey
                </Link>
                <div className='collapse navbar-collapse'>
                  <ul className="navbar-nav ml-auto" >
                    <li className='nav-item'>
                      <Link to={"/student"} className="nav-link">
                        Home
                      </Link>
                    </li>
                    <li className='nav-item'>
                      <Link to={"/profile"} className="nav-link">
                        Profile
                      </Link>
                    </li>
                    <li className="nav-item">
                      <Link to={"/"} className="nav-link" onClick={logOut}>
                        Logout
                      </Link>
                    </li>        
                  </ul>
                </div>
              </div>
            </nav> 
            ) : (
              showStudent === false && showTutor === true && (              
              <nav className="navbar navbar-expand-lg navbar-light">
                <div className='container'>
                  <Link to={"/"} className="navbar-brand">
                    Get Smart Tutoring (GST)
                  </Link>
                  <div className='collapse navbar-collapse'>
                    <ul className="navbar-nav ml-auto" >
                      <li className='nav-item'>
                        <Link to={"/tutor"} className="nav-link">
                          Home
                        </Link>
                      </li>
                      <li className='nav-item'>
                        <Link to={"/profile"} className="nav-link">
                          Profile
                        </Link>
                      </li>
                      <li className="nav-item">
                        <Link to={"/"} className="nav-link" onClick={logOut}>
                          Logout
                        </Link>
                      </li>        
                    </ul>
                  </div>
                </div>
              </nav> 
              )
            )
        ) 
      }

      
      { !currentUser ? (
          <div className="auth-wrapper">
            <div className='auth-inner'>
              <Routes>
                <Route exact path={"/"} element={<Login />} />
                <Route exact path="/login" element={<Login />} />
                <Route exact path="/register/student" element={<RegisterStudent />} />
                <Route exact path="/register/tutor" element={<RegisterTutor />} />
              </Routes>
            </div>
        </div>
        ) : ( 
          showStudent === true && showTutor === false ? (
            <div className="auth-wrapper">
              <div>
              <Routes>
                <Route exact path={"/"} element={<Student />} />
                <Route path="/student" element={<Student />} />
                <Route exact path="/profile" element={<Profile />} />
              </Routes>
              </div>
           </div>
          ) : (
            showPassenger === false && showRider === true && (
              <div className="auth-wrapper">
                <div>
                <Routes>
                  <Route exact path={"/"} element={<Rider />} />
                  <Route path="/rider" element={<Rider  />} />
                  <Route exact path="/profile" element={<Profile />} />
                  <Route exact path="/rides" element={<RideHistory />} />
                </Routes>
                </div>
              </div>
            )      
          )
        )
      }

    </div>
  );
};

export default App;
