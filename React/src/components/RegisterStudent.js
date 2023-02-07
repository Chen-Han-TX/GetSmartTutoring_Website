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

  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [school, setSchool] = useState("");
  const [selectedSubjects, setSelectedSubjects] = useState({
    "PSLE" : [],
    "O-Level": [],
    "A-Level": []
  });

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
    let updated = selectedSubjects
    updated["PSLE"] = subject
    setSelectedSubjects(selectedSubjects => ({
      ...updated
    }));
  }

const onChangeOlevel = (subject) => {
    let updated = selectedSubjects
    updated["O-Level"] = subject
    setSelectedSubjects(selectedSubjects => ({
        ...updated
    }));
}

const onChangeAlevel = (subject) => {
    let updated = selectedSubjects
    updated["A-Level"] = subject
    setSelectedSubjects(selectedSubjects => ({
        ...updated
    }));
}

  const handleRegister = (e) => {
    e.preventDefault();


    if  (name == "" || email == "" || password == "" || school == "") {
      alert("please fill up all the required fields!")
      return
    } else if (selectedSubjects["A-Level"].length == 0 && selectedSubjects["PSLE"].length == 0 && selectedSubjects["O-Level"].length == 0){
      alert("please choose at least one area of interest!")
      return
    } else {
      AuthService.register_student(name, email, password, school, selectedSubjects).then(
        (response) => {
          if (response.status == 200) {
            alert("Registered successfully!")
            navigate("/login");
            window.location.reload();
          } else {
            alert("Error, try again")
          }
        },
        (error) => {
          var resMessage = ""
          if (error.response.status === 406) {
            resMessage = "This email has been registered!"
            alert(resMessage)
          } else {
            resMessage =
            (error.response &&
              error.response.data &&
              error.response.data.message) ||
            error.message ||
            error.toString();
            alert(resMessage)
          }
        }
      );
    }
  };


  return (

        <Form className="auto-inner">
            <div>
            <h3>Register Student</h3>
              <hr className="hr"></hr>
              <h4>Basic Details</h4>
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

              <hr className="hr"></hr>
              <h4>Area of Interests</h4>
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
                <button onClick={handleRegister} className="btn btn-success btn-block">Register</button>
              </div>
            </div>
        </Form>
  );
};

export default RegisterStudent