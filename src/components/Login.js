import React, { useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import Form from "react-validation/build/form";
import Input from "react-validation/build/input";
import CheckButton from "react-validation/build/button";

import AuthService from "../services/auth.service";
//import RideServices from "../services/ride.service";

const required = (value) => {
  if (!value) {
    return (
      <div className="invalid-feedback d-block">
        This field is required!
      </div>
    );
  }
};

const Login = () => {
  const form = useRef();
  const checkBtn = useRef();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  const navigate = useNavigate();

  const onChangeEmail = (e) => {
    const email = e.target.value;
    setEmail(email);
  };

  const onChangePassword = (e) => {
    const password = e.target.value;
    setPassword(password);
  };

  const handleLogin = (e) => {
    e.preventDefault();

    setMessage("");
    setLoading(true);

    form.current.validateAll();

    if (checkBtn.current.context._errors.length === 0) {
      AuthService.login(email, password).then(
        () => {
            navigate("/");
            window.location.reload();
            /*

          // After login successfully, the user should have cookie with jwt token
          // Can then call the currentRide function to get the current ride and save in the localstorage
          // and call the allRides function to retrieve all the rides for the user
         
          RideServices.currentRide().then(
            ()=> {
              
              RideServices.allrides("").then(
                () => {
                  navigate("/");
                  window.location.reload();
                }
              )
            },
            (error) => {
              var resMessage = "Error: " + error.response.status.toString() + ": " + error.response.data;
              setLoading(false);
              setMessage(resMessage);

            }
            
          )

          */
        },
        (error) => {
          var resMessage = "Error: " + error.response.status.toString() + ": " + error.response.data;
          setLoading(false);
          setMessage(resMessage);
        }
      );
    } else {
      setLoading(false);
    }
  };

  return (
        <Form onSubmit={handleLogin} ref={form}>
          <h3>Login</h3>
          <div className="mb-3">
            <label htmlFor="email">Email Address</label>
            <Input
              type="text"
              className="form-control"
              name="email"
              value={email}
              onChange={onChangeEmail}
              validations={[required]}
            />
          </div>

          <div className="mb-3">
            <label htmlFor="password">Password</label>
            <Input
              type="password"
              className="form-control"
              name="password"
              value={password}
              onChange={onChangePassword}
              validations={[required]}
            />
          </div>

          <div className="d-grid">
            <button className="btn btn-primary btn-block" disabled={loading}>
              {loading && (
                <span className="spinner-border spinner-border-sm"></span>
              )}
              <span>Login</span>
            </button>
          </div>

          {message && (
            <div className="mb-3">
              <div className="alert alert-danger" role="alert">
                {message}
              </div>
            </div>
          )}
          <CheckButton style={{ display: "none" }} ref={checkBtn} />
        </Form>
  );
};

export default Login;