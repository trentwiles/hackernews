import Cookies from "js-cookie";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function UserRedirect() {
  const navigate = useNavigate();

  useEffect(() => {
    if (Cookies.get("token") == undefined) {
      navigate("/login");
      return;
    }

    fetch("http://localhost:3000/api/v1/me", {
      headers: {
        Authorization: "Bearer " + Cookies.get("token"),
      },
    })
      .then((res) => {
        if (!res.ok) {
          throw new Error("Network response was not ok");
        }
        return res.json();
      })
      .then((json) => {
        navigate("/u/" + json.username);
        return;
      })
      .catch(() => {
        navigate("/login");
        return;
      });
  }, []);

  return <p>Loading...</p>;
}
