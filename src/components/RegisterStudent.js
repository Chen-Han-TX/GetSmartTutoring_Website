import React, { useState, useRef } from "react";
import Form from "react-validation/build/form";
import Input from "react-validation/build/input";
import CheckButton from "react-validation/build/button";
import DropdownMultiselect from "react-multiselect-dropdown-bootstrap";
import { isEmail } from "validator";

import AuthService from "../services/auth.service";
import SubjectServices from "../services/subject.service";
import { useNavigate } from "react-router-dom";

const required = (value) => {
  if (!value) {
    return (
      <div className="invalid-feedback d-block">
        This field is required!
      </div>
    );
  }
};

const validEmail = (value) => {
  if (!isEmail(value)) {
    return (
      <div className="invalid-feedback d-block">
        This email cannot be registered!
      </div>
    );
  }
};


const vpassword = (value) => {
  if (value.length < 8 || value.length > 40) {
    return (
      <div className="invalid-feedback d-block">
        The password must be between 8 and 40 characters.
      </div>
    );
  }
};


const RegisterStudent = () => {
  const form = useRef();
  const checkBtn = useRef();


  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [school, setSchool] = useState("");

  var selectedSubjects = {
    "PSLE": [], "O-Level": [], "A-Level": []
  }

  const subjects = SubjectServices.getAllSubjects()
  const PSLESubjects = subjects["PSLE"].sort()
  const OlevelSubjects = subjects["O-Level"].sort()
  const AlevelSubjects = subjects["A-Level"].sort()

  const PSLEArray = [];
  for (let i = 0; i < PSLESubjects.length; i++) {
    var value = PSLESubjects[i]
    PSLEArray.push({key: value, label: value })
  } 

  const OLevelArray = [];
  for (let i = 0; i < OlevelSubjects.length; i++) {
    var value = OlevelSubjects[i] 
    OLevelArray.push({key: value, label: value })
  } 

  const ALevelArray = [];
  for (let i = 0; i < AlevelSubjects.length; i++) {
    var value = AlevelSubjects[i] 
    ALevelArray.push({key: value, label: value })
  } 

  const [successful, setSuccessful] = useState(false);
  const [message, setMessage] = useState("");
  const navigate = useNavigate();


  const onChangeEmail = (e) => {
    const email = e.target.value;
    setEmail(email);
  };

  const onChangeName = (e) => {
    const name = e.target.value;
    setName(name);
  };


  const onChangePassword = (e) => {
    const password = e.target.value;
    setPassword(password);
  };

  const onChangeSchool = (e) => {
    const school = e.target.value;
    setSchool(school);
  };

  const OnChangePSLE = (subject) => {
    selectedSubjects["PSLE"] = subject
  }

  const onChangeOlevel = (subject) => {
    selectedSubjects["O-Level"] = subject
  }

  const onChangeAlevel = (subject) => {
    selectedSubjects["A-Level"] = subject
  }


  const handleRegister = (e) => {
    e.preventDefault();

    setMessage("");
    setSuccessful(false);

    form.current.validateAll();

    // If passed validation, call auth service to send the API request
    if (checkBtn.current.context._errors.length === 0) {
      
      AuthService.register_student(name, email, password, school, selectedSubjects).then(
        (response) => {
          if (response.status == 200) {
            setMessage("Registered Successfully!");
            setSuccessful(true);
            setTimeout(function () {
              navigate("/login");
              window.location.reload();
            }, 2000);
          } else {
            setMessage("Error, try again!");
            setSuccessful(false);
            setTimeout(function () {
              setMessage("");
            }, 5000);
          }
        },
        (error) => {
          var resMessage = ""
          if (error.response.status === 406) {
            resMessage = "This email has been registered!"
          } else {
            resMessage =
            (error.response &&
              error.response.data &&
              error.response.data.message) ||
            error.message ||
            error.toString();
          }
          setMessage(resMessage);
          setSuccessful(false);
        }
      );
    }
  };


  return (
        <Form onSubmit={handleRegister} ref={form}>
          <h3>Register Student</h3>
          {!successful && (

            <div>

              <div className="mb-3">
                <label htmlFor="name">Name</label>
                <Input
                  type="text"
                  className="form-control"
                  name="name"
                  value={name}
                  onChange={onChangeName}
                  validations={[required]}
                />
              </div>

              <div className="mb-3">
                <label htmlFor="email">Email</label>
                <Input
                  type="text"
                  className="form-control"
                  name="email"
                  value={email}
                  onChange={onChangeEmail}
                  validations={[required, validEmail]}
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
                  validations={[required, vpassword]}
                />
              </div>

              <div className="mb-3">
                <label htmlFor="school">School</label>
                <Input
                  type="text"
                  className="form-control"
                  name="school"
                  value={school}
                  onChange={onChangeSchool}
                  validations={[required]}
                />
              </div>

              <label htmlFor="options">-----Area of Interests----</label>
              <div className="mb-3">
                  <label htmlFor="options">PSLE</label>
                  <DropdownMultiselect options={PSLEArray} name="pslesubjects" 
                  handleOnChange={(selected) => {
                    OnChangePSLE(selected);
                  }}/>
              </div>

              <div className="mb-3">
                  <label htmlFor="options">O-Level</label>
                  <DropdownMultiselect options={OLevelArray} name="olevelsubjects"
                  handleOnChange={(selected) => {
                    onChangeOlevel(selected);
                  }}/>
              </div>

              <div className="mb-3">
                  <label htmlFor="options">A-Level</label>
                  <DropdownMultiselect options={ALevelArray} name="alevelsubjects"
                  handleOnChange={(selected) => {
                    onChangeAlevel(selected);
                  }}/>
              </div>

              <div className="d-grid">
                <button className="btn btn-success btn-block">Register</button>
              </div>
            </div>
          )}

          {message && (
            <div className="mb-3">
              <div
                className={
                  successful ? "alert alert-success" : "alert alert-danger"
                }
                role="alert"
              >
                {message}
              </div>
            </div>
          )}
          <CheckButton style={{ display: "none" }} ref={checkBtn} />
        </Form>
  );
};

export default RegisterStudent