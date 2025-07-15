import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./components/Login";
import Home from "./components/Home";
import LoginThanks from "./components/LoginThanks";
import MagicLink from "./components/MagicLink";
import { GoogleReCaptchaProvider } from "react-google-recaptcha-v3";
import Submit from "./components/Submit";
import ProtectedRoute from "./components/auth/ProtectedRoute";
import User from "./components/User";
import Submission from "./components/Submission";
import UserRedirect from "./components/UserRedirect";
import Logout from "./components/Logout";
import Settings from "./components/Settings";
import Search from "./components/Search";
import AdminPanel from "./components/AdminPanel";

function App() {
  const REQUIRED_ENV_VARS: string[] = ["VITE_API_ENDPOINT"];

  REQUIRED_ENV_VARS.map((envVar) => {
    if (import.meta.env[envVar] === undefined) {
      throw new Error("unable to find required .env variable: " + envVar);
    }
  });

  return (
    <GoogleReCaptchaProvider reCaptchaKey="6LcopnUrAAAAACZBUINoyS__gkqGOTm-Nj4qhIm1">
      <Router>
        <div className="App">
          <Routes>
            {/* PUBLIC ROUTES */}
            <Route path="/" element={<Home />} />
            <Route path="/latest" element={<Home sortType="latest" />} />
            <Route path="/top" element={<Home sortType="best" />} />

            <Route path="/login" element={<Login serviceName="HackerNews" />} />
            <Route path="/logout" element={<Logout />} />
            <Route
              path="/magic"
              element={<MagicLink serviceName="HackerNews" />}
            />
            <Route path="/login-thanks" element={<LoginThanks />} />
            <Route path="/u/:username" element={<User />} />
            <Route path="/u/" element={<UserRedirect />} />
            <Route path="/account/submissions" element={<UserRedirect />} />
            <Route path="/submission/:sid" element={<Submission />} />
            <Route path="/search" element={<Search />} />

            {/* PROTECTED (AUTH REQUIRED) ROUTES */}
            <Route
              path="/submit"
              element={
                <ProtectedRoute>
                  <Submit />
                </ProtectedRoute>
              }
            />
            <Route
              path="/account/settings"
              element={
                <ProtectedRoute>
                  <Settings />
                </ProtectedRoute>
              }
            />
            <Route
              path="/account/admin"
              element={
                <ProtectedRoute>
                  <AdminPanel />
                </ProtectedRoute>
              }
            />
          </Routes>
        </div>
      </Router>
    </GoogleReCaptchaProvider>
  );
}

export default App;
