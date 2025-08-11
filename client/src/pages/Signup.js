import React, { useState, useEffect, startTransition } from "react";
import axios from "axios";
import { useNavigate, Link } from "react-router-dom";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import * as constants from "../utils/constant/Constant";

export default function Signup() {
  const navigate = useNavigate();
  const toastOptions = {
    position: "top-right",
    autoClose: 3000,
    pauseOnHover: true,
    draggable: true,
    theme: "dark",
  };

  const [values, setValues] = useState({
    name: "",
    user_name: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      navigate("/chat");
    }
  }, [navigate]);

  const handleChange = (event) => {
    setValues({ ...values, [event.target.name]: event.target.value });
  };

  const handleValidation = () => {
    const { name, password, confirmPassword, user_name, email } = values;
    if (password !== confirmPassword) {
      toast.error("Password and confirm password should be the same.", toastOptions);
      return false;
    } else if (name === "") {
      toast.error("Name is required.", toastOptions);
      return false;
    } else if (user_name.length < 3) {
      toast.error("Username should be greater than 3 characters.", toastOptions);
      return false;
    } else if (password.length < 6) {
      toast.error("Password should be equal or greater than 6 characters.", toastOptions);
      return false;
    } else if (email === "") {
      toast.error("Email is required.", toastOptions);
      return false;
    }
    return true;
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (handleValidation()) {
      const { name, email, user_name, password } = values;
      try {
        await new Promise((resolve) => {
          startTransition(() => {
            axios
              .post(constants.ROUTES.SIGNUP, { user_name, name, email, password })
              .then(({ data }) => {
                if (data?.code === constants.API_STATUS.FAILURE_CODE) {
                  toast.error(data?.message, toastOptions);
                } else {
                  localStorage.setItem("token", data?.meta?.token);
                  localStorage.setItem("userId", data?.data?._id);
                  localStorage.setItem("userName", data?.data?.user_name);
                  localStorage.setItem("name", data?.data?.name);
                  navigate(0);
                }
              })
              .catch((error) => {
                console.log("Registration Failed: ", error.message);
                if (error.response && error.response.status === 500) {
                  toast.error("Backend is currently unavailable. Please try again later.", toastOptions);
                } else {
                  navigate("/error", { state: { errorMessage: "An error occurred during registration. Please try again." } });
                }
                resolve();
              });
          });
        });
      } catch (error) {
        console.error("Error during registration:", error);
      }
    }
  };

  return (
    <>
      <div className="signup-body">
        <div className="form-container">
          <form onSubmit={(event) => handleSubmit(event)} className="form">
            <h2 className="form-title">Create an Account</h2>
            <input
              type="text"
              placeholder="Name"
              autoComplete="off"
              name="name"
              onChange={(e) => handleChange(e)}
              className="form-input"
            />
            <input
              type="text"
              placeholder="Username"
              autoComplete="off"
              name="user_name"
              onChange={(e) => handleChange(e)}
              className="form-input"
            />
            <input
              type="email"
              placeholder="Email"
              autoComplete="off"
              name="email"
              onChange={(e) => handleChange(e)}
              className="form-input"
            />
            <input
              type="password"
              placeholder="Password"
              name="password"
              onChange={(e) => handleChange(e)}
              className="form-input"
            />
            <input
              type="password"
              placeholder="Confirm Password"
              name="confirmPassword"
              onChange={(e) => handleChange(e)}
              className="form-input"
            />
            <button type="submit" className="form-button">
              Create User
            </button>
            <p className="form-text">
              Already have an account?{" "}
              <Link to="/login" className="form-link">
                Login
              </Link>
            </p>
          </form>
        </div>
      </div>
      <ToastContainer />
    </>
  );
}
