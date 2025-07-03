import { BrowserRouter as Router, Routes, Route, Link, Navigate } from "react-router-dom";
import Login from "./components/Login";
import MagicLink from "./components/MagicLink";
import { GoogleReCaptchaProvider } from "react-google-recaptcha-v3";

function App() {
  return (
    <GoogleReCaptchaProvider reCaptchaKey="6LcopnUrAAAAACZBUINoyS__gkqGOTm-Nj4qhIm1">
      <Router>
        <div className="App">
          {/* Routes */}
          <Routes>
            <Route path="/" element={<Navigate to="/login" replace />} />
            <Route path="/login" element={<Login serviceName="HackerNews" />} />
            <Route path="/magic" element={<MagicLink serviceName="HackerNews" />} />
          </Routes>
        </div>
      </Router>
    </GoogleReCaptchaProvider>
  );
}

export default App;