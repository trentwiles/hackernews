import {
  BrowserRouter as Router,
  Routes,
  Route,
  useLocation,
} from "react-router-dom";
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
import Forbidden from "./components/Forbidden";
import { Helmet } from "react-helmet";
import ReactGA from "react-ga4";
import { useEffect } from "react";
import Privacy from "./components/Privacy";
import UrlCheck from "./components/UrlCheck";

const SERVICE_NAME = import.meta.env.VITE_SERVICE_NAME;

function AnalyticsTracker() {
  const location = useLocation();

  useEffect(() => {
    ReactGA.send({
      hitType: "pageview",
      page: location.pathname + location.search,
    });
  }, [location]);

  return null;
}

function App() {
  ReactGA.initialize(import.meta.env.VITE_GOOGLE_ANALYTICS_ID);

  const REQUIRED_ENV_VARS: string[] = ["VITE_API_ENDPOINT"];

  REQUIRED_ENV_VARS.map((envVar) => {
    if (import.meta.env[envVar] === undefined) {
      throw new Error("unable to find required .env variable: " + envVar);
    }
  });

  return (
    <GoogleReCaptchaProvider
      reCaptchaKey={import.meta.env.VITE_GOOGLE_RECAPTCHA_PUBLIC_KEY}
    >
      <Router>
        <AnalyticsTracker />
        <div className="App">
          <Routes>
            {/* PUBLIC ROUTES */}
            <Route
              path="/"
              element={
                <>
                  <Helmet>
                    <title>Home | {SERVICE_NAME}</title>
                  </Helmet>
                  <Home />
                </>
              }
            />
            <Route
              path="/latest"
              element={
                <>
                  <Helmet>
                    <title>Latest | {SERVICE_NAME}</title>
                  </Helmet>
                  <Home sortType="latest" />
                </>
              }
            />
            <Route
              path="/top"
              element={
                <>
                  <Helmet>
                    <title>Top | {SERVICE_NAME}</title>
                  </Helmet>
                  <Home sortType="best" />
                </>
              }
            />
            <Route
              path="/login"
              element={
                <>
                  <Helmet>
                    <title>Login | {SERVICE_NAME}</title>
                  </Helmet>
                  <Login serviceName={SERVICE_NAME} />
                </>
              }
            />
            <Route
              path="/logout"
              element={
                <>
                  <Helmet>
                    <title>Logout | {SERVICE_NAME}</title>
                  </Helmet>
                  <Logout />
                </>
              }
            />
            <Route
              path="/magic"
              element={
                <>
                  <Helmet>
                    <title>Magic Link | {SERVICE_NAME}</title>
                  </Helmet>
                  <MagicLink serviceName={SERVICE_NAME} />
                </>
              }
            />
            <Route
              path="/login-thanks"
              element={
                <>
                  <Helmet>
                    <title>Login Thanks | {SERVICE_NAME}</title>
                  </Helmet>
                  <LoginThanks />
                </>
              }
            />
            <Route
              path="/u/:username"
              element={
                <>
                  <Helmet>
                    <title>User Profile | {SERVICE_NAME}</title>
                  </Helmet>
                  <User />
                </>
              }
            />
            <Route
              path="/u/"
              element={
                <>
                  <Helmet>
                    <title>User | {SERVICE_NAME}</title>
                  </Helmet>
                  <UserRedirect />
                </>
              }
            />
            <Route
              path="/account/submissions"
              element={
                <>
                  <Helmet>
                    <title>My Submissions | {SERVICE_NAME}</title>
                  </Helmet>
                  <UserRedirect />
                </>
              }
            />
            <Route
              path="/submission/:sid"
              element={
                <>
                  <Helmet>
                    <title>Submission | {SERVICE_NAME}</title>
                  </Helmet>
                  <Submission />
                </>
              }
            />
            <Route
              path="/urlCheck"
              element={
                <>
                  <Helmet>
                    <title>Please Wait... | {SERVICE_NAME}</title>
                  </Helmet>
                  <UrlCheck />
                </>
              }
            />
            <Route
              path="/search"
              element={
                <>
                  <Helmet>
                    <title>Search | {SERVICE_NAME}</title>
                  </Helmet>
                  <Search />
                </>
              }
            />
            <Route
              path="/403"
              element={
                <>
                  <Helmet>
                    <title>Forbidden | {SERVICE_NAME}</title>
                  </Helmet>
                  <Forbidden />
                </>
              }
            />
            {/* PROTECTED (AUTH REQUIRED) ROUTES */}
            <Route
              path="/submit"
              element={
                <ProtectedRoute>
                  <>
                    <Helmet>
                      <title>Submit | {SERVICE_NAME}</title>
                    </Helmet>
                    <Submit />
                  </>
                </ProtectedRoute>
              }
            />
            <Route
              path="/account/settings"
              element={
                <ProtectedRoute>
                  <>
                    <Helmet>
                      <title>Settings | {SERVICE_NAME}</title>
                    </Helmet>
                    <Settings />
                  </>
                </ProtectedRoute>
              }
            />
            <Route
              path="/account/privacy"
              element={
                <ProtectedRoute>
                  <>
                    <Helmet>
                      <title>Privacy Settings | {SERVICE_NAME}</title>
                    </Helmet>
                    <Privacy />
                  </>
                </ProtectedRoute>
              }
            />
            <Route
              path="/account/admin"
              element={
                <ProtectedRoute>
                  <>
                    <Helmet>
                      <title>Admin Panel | {SERVICE_NAME}</title>
                    </Helmet>
                    <AdminPanel />
                  </>
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
