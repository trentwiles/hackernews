import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./components/Login";
import Home from "./components/Home";
import LoginThanks from "./components/LoginThanks";
import MagicLink from "./components/MagicLink";
import { GoogleReCaptchaProvider } from "react-google-recaptcha-v3";

function App() {
  return (
    <GoogleReCaptchaProvider reCaptchaKey="6LcopnUrAAAAACZBUINoyS__gkqGOTm-Nj4qhIm1">
      <Router>
        <div className="App">
          <Routes>
            {/* PUBLIC ROUTES */}
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login serviceName="HackerNews" />} />
            <Route path="/magic" element={<MagicLink serviceName="HackerNews" />} />
            <Route path="/login-thanks" element={<LoginThanks />} />

            {/* PROTECTED (AUTH REQUIRED) ROUTES */}
            {/* <Route path="/submit" element={
              <ProtectedRoute>
                <Submit />
              </ProtectedRoute>
            } /> */}
          </Routes>
        </div>
      </Router>
    </GoogleReCaptchaProvider>
  );
}

export default App;